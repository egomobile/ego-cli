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
	"embed"
	b64 "encoding/base64"
	"errors"
	"html"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/thatisuday/commando"
)

// www holds static web server content
//go:embed www/*
var www embed.FS

type GitBranchAndRemotes struct {
	Branch  string
	Remotes []string
}

// information about the file and folder structre of
// a GitHub repository
type GitHubRepoFilesAndFolders struct {
	// the tree as item list
	Tree GitHubRepoTree `json:"tree,omitempty"`
}

// the tree of the GitHub repository
type GitHubRepoTree = []GitHubRepoTreeItem

// and item of GitHubRepoTree
type GitHubRepoTreeItem = struct {
	// the relative path inside the GitHub repository
	Path string `json:"path,omitempty"`
}

func BuildHtmlPage(content string, title string, css string, js string) (string, error) {
	indexHtml, err := GetWWW().ReadFile("www/index.html")
	if err != nil {
		return "", err
	}

	bootstrapCSS, err := www.ReadFile("www/css/bootstrap.min.css")
	if err != nil {
		return "", err
	}

	bootstrapJS, err := www.ReadFile("www/js/bootstrap.bundle.min.js")
	if err != nil {
		return "", err
	}

	egoLogoSVG, err := www.ReadFile("www/img/ego.svg")
	if err != nil {
		return "", err
	}

	egoLogo := b64.StdEncoding.EncodeToString([]byte(egoLogoSVG))
	egoLogo = "data:image/svg+xml;base64," + egoLogo

	titleSuffix := strings.TrimSpace(title)
	if titleSuffix != "" {
		titleSuffix = " :: " + titleSuffix
	}

	var h = string(indexHtml)
	h = strings.ReplaceAll(h, "<EGO-PAGE-TITLE-SUFFIX />", html.EscapeString(titleSuffix))
	h = strings.ReplaceAll(h, "<EGO-LOGO />", egoLogo)
	h = strings.ReplaceAll(h, "<EGO-HTML />", content)
	h = strings.ReplaceAll(h, "<EGO-BOOTSTRAP-CSS />", string(bootstrapCSS))
	h = strings.ReplaceAll(h, "<EGO-BOOTSTRAP-JS />", string(bootstrapJS))
	h = strings.ReplaceAll(h, "<EGO-CSS />", css)
	h = strings.ReplaceAll(h, "<EGO-JS />", js)
	h = strings.ReplaceAll(h, "<EGO-APP-VERSION />", html.EscapeString(commando.DefaultCommandRegistry.Version))

	return h, nil
}

// EnsureEgoDir() - returns the full path of  ~/.ego directory
// and ensures that it exists
func EnsureEgoDir() (string, error) {
	usersHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err // cannot get home directory
	}

	// ~/.ego
	egoDir := filepath.Join(usersHomeDir, ".ego")

	stat, err := os.Stat(egoDir)
	if err != nil {
		if os.IsNotExist(err) {
			// does not exist => try create it
			err := os.MkdirAll(egoDir, os.ModePerm)

			if err == nil {
				return egoDir, nil // directory created
			} else {
				return "", err // could not create directory
			}
		} else {
			return "", err // could not get file information
		}
	}

	if stat.IsDir() {
		return egoDir, nil // directory already exists
	}

	// no directory
	return "", errors.New(".ego is no directory")
}

// GetGitBranchAndRemotes detects the current branch
// and all remotes of the current working directory
func GetGitBranchAndRemotes() (GitBranchAndRemotes, error) {
	var repo GitBranchAndRemotes

	cwd, err := os.Getwd()
	if err != nil {
		return repo, err
	}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = cwd

	output, err := cmd.Output()
	if err != nil {
		return repo, err
	}

	repo.Branch = strings.TrimSpace(strings.Join(SplitStringAndCleanup(string(output), "\n"), "\n"))

	cmd = exec.Command("git", "remote")
	cmd.Dir = cwd

	output, err = cmd.Output()
	if err != nil {
		return repo, err
	}

	repo.Remotes = SplitStringAndCleanup(string(output), "\n")

	return repo, nil
}

func GetWWW() embed.FS {
	return www
}

// SplitStringAndCleanup splits a string in s using a separator in sep
// and removes all empty lines
func SplitStringAndCleanup(s string, sep string) []string {
	var result []string

	lines := strings.Split(s, sep)
	for _, l := range lines {
		var strToAdd = strings.TrimSpace(l)

		if len(strToAdd) > 0 {
			result = append(result, strToAdd)
		}
	}

	return result
}

// WithSpinner invokes and action by using an initial preifx text
func WithSpinner(text string, action func(s *spinner.Spinner)) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.UpdateCharSet(spinner.CharSets[14])
	s.Suffix = " " + text
	s.Color("green", "bold")

	s.Start()
	action(s)
	s.Stop()
}
