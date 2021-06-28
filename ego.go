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
	"github.com/thatisuday/commando"
)

func main() {
	// setup commando
	commando.
		SetExecutableName("ego").
		SetVersion("0.1.0").
		SetDescription("Command Line Interface, which is designed to handle things, like Dev(Op) and other common tasks, much faster")

	Setup_chuck_Command()
	Setup_gitpull_Command()
	Setup_gitpush_Command()
	Setup_gitsync_Command()
	Setup_localip_Command()
	Setup_nodeinstall_Command()
	Setup_publicip_Command()

	// parse CLI args
	commando.Parse(nil)
}
