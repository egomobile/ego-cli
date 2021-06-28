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
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

type GitBranchAndRemotes struct {
	Branch  string
	Remotes []string
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
