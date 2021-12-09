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

	utils "github.com/egomobile/ego-cli/utils"
	"github.com/thatisuday/commando"
)

func update_execute(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	softwareRepo, err := utils.UpdateLocalSoftwareRepository()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Updated software repository with", len(softwareRepo), "entries")
}

func Setup_update_Command() {
	commando.
		Register("update").
		SetShortDescription("updates the local list of software repositories").
		SetAction(update_execute)
}
