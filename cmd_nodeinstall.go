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

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/briandowns/spinner"
	"github.com/thatisuday/commando"
)

func nodeinstall_run(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	audit, err := flags["audit"].GetBool()
	runAuditFix := err == nil && audit

	update, err := flags["update"].GetBool()
	runUpdate := err == nil && update

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var package_json = filepath.Join(cwd, "/package.json")
	if _, err := os.Stat(package_json); os.IsNotExist(err) {
		log.Fatal("No package.json found")
	}

	var node_modules = filepath.Join(cwd, "/node_modules")

	// remove node_modules?
	if _, err := os.Stat(node_modules); !os.IsNotExist(err) {
		WithSpinner("Removing node_modules/ folder ...", func(s *spinner.Spinner) {
			err = os.RemoveAll(node_modules)
			if err != nil {
				log.Fatal(err)
			}
		})

		fmt.Println("")
	}

	// execute npm install
	WithSpinner("Executing npm install ...", func(s *spinner.Spinner) {
		cmd := exec.Command("npm", "install")
		cmd.Dir = cwd

		err := cmd.Run()
		if err != nil {
			output, err2 := cmd.Output()
			if err2 != nil {
				log.Fatal(err)
			} else {
				log.Fatal(string(output))
			}
		}

		fmt.Println("")
	})

	if runUpdate {
		// run npm update

		WithSpinner("Executing npm update ...", func(s *spinner.Spinner) {
			cmd := exec.Command("npm", "update")
			cmd.Dir = cwd

			err := cmd.Run()
			if err != nil {
				output, err2 := cmd.Output()
				if err2 != nil {
					log.Fatal(err)
				} else {
					log.Fatal(string(output))
				}
			}

			fmt.Println("")
		})
	}

	if runAuditFix {
		// run npm audit fix

		WithSpinner("Executing npm audit fix ...", func(s *spinner.Spinner) {
			cmd := exec.Command("npm", "audit", "fix")
			cmd.Dir = cwd

			err := cmd.Run()
			if err != nil {
				output, err2 := cmd.Output()
				if err2 != nil {
					log.Fatal(err)
				} else {
					log.Fatal(string(output))
				}
			}

			fmt.Println("")
		})
	}
}

func Setup_nodeinstall_Command() {
	commando.
		Register("node-install").
		SetShortDescription("runs \"npm install\"").
		SetDescription("Deletes \"node_modules\" and executes \"npm install\" with optional \"npm update\" and \"npm audit fix\"").
		AddFlag("audit,a", "run \"audit fix\" after \"install\" or \"update\"", commando.Bool, nil).
		AddFlag("update,u", "run \"update\" after \"install\"", commando.Bool, nil).
		SetAction(nodeinstall_run)
}
