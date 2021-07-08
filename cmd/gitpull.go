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
	"os/exec"

	"github.com/briandowns/spinner"
	"github.com/thatisuday/commando"

	utils "github.com/egomobile/ego-cli/utils"
)

func gitpull_run(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	repo, err := utils.GetGitBranchAndRemotes()
	if err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range repo.Remotes {
		utils.WithSpinner("Pulling from '"+r+"' ...", func(s *spinner.Spinner) {
			cmd := exec.Command("git", "pull", r, repo.Branch)
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

func Setup_gitpull_Command() {
	commando.
		Register("git-pull").
		SetShortDescription("git pull").
		SetDescription("Does a \"git pull\", in the current branch for all remotes in one command").
		SetAction(gitpull_run)
}
