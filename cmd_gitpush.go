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

func gitpush_run(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	repo, err := GetGitBranchAndRemotes()
	if err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range repo.Remotes {
		WithSpinner("Pushing to '"+r+"' ...", func(s *spinner.Spinner) {
			cmd := exec.Command("git", "push", r, repo.Branch)
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

func Setup_gitpush_Command() {
	commando.
		Register("git-push").
		SetShortDescription("git push").
		SetDescription("Does a \"git push\", in the current branch for all remotes in one command").
		SetAction(gitpush_run)
}
