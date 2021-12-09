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

package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	utils "github.com/egomobile/ego-cli/utils"
	"github.com/thatisuday/commando"
)

func install_execute(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	softwareRepo, err := utils.EnsureLocalSoftwareRepository()
	if err != nil {
		log.Fatalln(err)
	}

	packageName := strings.TrimSpace(strings.ToLower(args["package"].Value))

	existingPackage, ok := softwareRepo[packageName]
	if !ok {
		os.Exit(404)
	}

	var existingArchitectureItem utils.SoftwarePackageSourcePlatformArchitectureItem

	findArchitecture := func() {
		found := false

		goOS := strings.TrimSpace(strings.ToLower(runtime.GOOS))
		goArch := strings.TrimSpace(strings.ToLower(runtime.GOARCH))

		for k, src := range existingPackage.Sources {
			keys := strings.Split(k, ",")
			for _, k := range keys {
				platform := strings.TrimSpace(strings.ToLower(k))

				if platform == goOS {
					for k, archItem := range src {
						keys := strings.Split(k, ",")

						for _, k := range keys {
							architecture := strings.TrimSpace(strings.ToLower(k))
							if architecture == goArch {
								existingArchitectureItem = archItem
								found = true
							}
						}
					}
				}
			}
		}

		if !found {
			fmt.Println("Could not find package", packageName, "that can be installed on", goOS, "/", goArch, "!")

			os.Exit(404)
		}
	}

	findArchitecture()

	fmt.Println("existingArchitectureItem", existingArchitectureItem)
}

func Setup_install_Command() {
	commando.
		Register("install").
		SetShortDescription("installs software").
		AddArgument("package", "the name of the package to install", "").
		SetAction(install_execute)
}
