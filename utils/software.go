// This file is part of the ego-cli distribution.
// Copyright (c) Next.e.GO Mobile SE, Aachen, Germany (https://e-go-mobile.com/)
//
// ego-cli is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, version 3.
//
// ego-cli is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
)

// a repository with all known (software) packages
// which can be installed
//
// each key is the name for the CLI, e.g.:
//
// "vscode" => ego install vscode
type SoftwarePackageRepository = map[string]SoftwarePackage

// a value for a key of SoftwarePackageRepository
//
// contains all information about a software package
type SoftwarePackage struct {
	// optional information about the author
	Author SoftwarePackageAuthor `json:"author,omitempty"`
	// a description for the user
	Description string `json:"description"`
	// optional homepage URL with more information
	Homepage string `json:"homepage"`
	// optional information about the maintainer of that package
	Maintainer SoftwarePackageMaintainer `json:"maintainer"`
	// optional list of source code repositories
	Repositories []SoftwarePackageSourceCodeRepository `json:"repositories,omitempty"`
	// the list of sources, grouped by platforms
	Sources SoftwarePackageSourcePlatforms `json:"sources"`
}

// information about the author
type SoftwarePackageAuthor struct {
	// optional list of contact addresses
	Contacts SoftwarePackageContacts `json:"contacts,omitempty"`
	// the URL to the author's homepage
	Homepage string `json:"homepage"`
	// the name
	Name string `json:"name"`
}

// information about the maintainer
type SoftwarePackageMaintainer struct {
	// optional list of contact addresses
	Contacts SoftwarePackageContacts `json:"contacts,omitempty"`
	// homepage URL
	Homepage string `json:"homepage"`
	// the name
	Name string `json:"name"`
}

// list of contact addresses
type SoftwarePackageContacts = []string

// list of software architectures
//
// each key is the name of one or more supported architecture
// separated by commas, if needed
type SoftwarePackageSourcePlatformArchitectures = map[string]SoftwarePackageSourcePlatformArchitectureItem

// a value for a key of SoftwarePackageSourcePlatformArchitectures
//
// contains all information to install the package on the
// underlying architecture
type SoftwarePackageSourcePlatformArchitectureItem struct {
	// optional list of requirement (sources)
	Requirements []string `json:"requirements,omitempty"`
	// the source (URL)
	Source string `json:"source"`
}

// source code repository information
type SoftwarePackageSourceCodeRepository struct {
	// the type
	Type string `json:"type"`
	// the URL
	Url string `json:"url"`
}

// list of supported platforms
//
// each key is the name of one or more supported platform
// separated by commas, if needed
type SoftwarePackageSourcePlatforms = map[string]SoftwarePackageSourcePlatformArchitectures

// TODO: update to "master"
const repoBranchUrl = "feature/install"

func DownloadSoftwareRepositoryFileList() ([]string, error) {
	// download file tree of repository on GitHib
	// TODO: update repoBranchUrl to "master"
	resp, err := http.Get("https://api.github.com/repos/egomobile/ego-cli/git/trees/" + repoBranchUrl + "?recursive=1")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err // could not parse body
	}

	var filesAndFolders GitHubRepoFilesAndFolders
	err = json.Unmarshal(bodyData, &filesAndFolders)
	if err != nil {
		return nil, err // could not parse
	}

	urls := []string{}

	if filesAndFolders.Tree != nil {
		for _, treeItem := range filesAndFolders.Tree {
			if strings.HasPrefix(treeItem.Path, "software/") { // only if in software/ sub folder
				if strings.HasSuffix(treeItem.Path, ".json") { // only JSON files
					// TODO: update repoBranchUrl to "master"
					fullUrl := "https://raw.githubusercontent.com/egomobile/ego-cli/" + repoBranchUrl + "/" + treeItem.Path

					urls = append(urls, fullUrl)
				}
			}
		}
	}

	sort.Strings(urls)

	return urls, nil
}

func DownloadSoftwareRepository() (SoftwarePackageRepository, error) {
	softwareRepo := make(SoftwarePackageRepository)

	packageFiles, err := DownloadSoftwareRepositoryFileList()
	if err != nil {
		return softwareRepo, err
	}

	for _, pkgFileUrl := range packageFiles {
		pkgName := path.Base(pkgFileUrl)
		pkgName = strings.TrimSpace(strings.ToLower(pkgName[:len(pkgName)-5]))

		if pkgName != "" {
			fmt.Println("Loading software package meta of", pkgName, "...")

			resp, err := http.Get(pkgFileUrl)
			if err == nil {
				defer resp.Body.Close()

				bodyData, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					var pkg SoftwarePackage
					err = json.Unmarshal(bodyData, &pkg)

					if err == nil {
						softwareRepo[pkgName] = pkg // OK => add package entry
					} else {
						fmt.Println("Parsing response body of software package", pkgName, " failed:", err)
					}
				} else {
					fmt.Println("Loading response body of software package", pkgName, " failed:", err)
				}
			} else {
				fmt.Println("Loading software package", pkgName, " failed:", err)
			}
		}
	}

	return softwareRepo, nil
}

func EnsureLocalSoftwareRepository() (SoftwarePackageRepository, error) {
	var softwareRepo SoftwarePackageRepository

	softwareFile, err := GetSoftwareFilePath()
	if err != nil {
		return softwareRepo, err // could not ensure ~/.ego directory
	}

	_, err = os.Stat(softwareFile)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := UpdateLocalSoftwareRepository() // needs to be downloaded
			if err != nil {
				return softwareRepo, err // could not download
			}
		} else {
			return softwareRepo, err // could not get stats
		}
	}

	jsonData, err := ioutil.ReadFile(softwareFile)
	if err != nil {
		return softwareRepo, err // could not read file
	}

	// try parse JSON data
	err = json.Unmarshal(jsonData, &softwareRepo)
	return softwareRepo, err
}

func GetSoftwareFilePath() (string, error) {
	egoDir, err := EnsureEgoDir()
	if err != nil {
		return "", err // could not ensure ~/.ego directory
	}

	softwareFile := path.Join(egoDir, "software.json")

	stat, err := os.Stat(softwareFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err // could not get stats
		}
	} else {
		if stat.IsDir() {
			// cannot be directory
			return "", errors.New(fmt.Sprintf("%v is no file", softwareFile))
		}
	}

	return softwareFile, nil
}

func UpdateLocalSoftwareRepository() (SoftwarePackageRepository, error) {
	softwareRepo, err := DownloadSoftwareRepository()
	if err != nil {
		return softwareRepo, err
	}

	softwareFile, err := GetSoftwareFilePath()
	if err != nil {
		return softwareRepo, err
	}

	jsonData, err := json.Marshal(softwareRepo)
	if err != nil {
		return softwareRepo, err
	}

	err = ioutil.WriteFile(softwareFile, jsonData, 0644)
	return softwareRepo, err
}