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

	"github.com/briandowns/spinner"
	"github.com/thatisuday/commando"
)

func gitsync_run(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	reverse, err := flags["reverse"].GetBool()
	pushThenPull := err == nil && reverse

	repo, err := GetGitBranchAndRemotes()
	if err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var pull = func(remote string) {
		WithSpinner("Pulling from '"+remote+"' ...", func(s *spinner.Spinner) {
			cmd := exec.Command("git", "pull", remote, repo.Branch)
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

	var push = func(remote string) {
		WithSpinner("Pushing to '"+remote+"' ...", func(s *spinner.Spinner) {
			cmd := exec.Command("git", "push", remote, repo.Branch)
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

	for _, r := range repo.Remotes {
		if pushThenPull {
			push(r)
			pull(r)
		} else {
			pull(r)
			push(r)
		}
	}
}

func Setup_gitsync_Command() {
	commando.
		Register("git-sync").
		SetShortDescription("git pull and push").
		SetDescription("Does a \"pull\" and \"push\" in one command, in the current branch for all remotes").
		AddFlag("reverse,r", "first do a push, then a pull", commando.Bool, nil).
		SetAction(gitsync_run)
}
