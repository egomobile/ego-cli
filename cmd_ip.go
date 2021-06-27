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
	"io/ioutil"
	"net/http"

	"github.com/thatisuday/commando"
)

func detectIPv4(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	var ipV4OrErrorMsg = ""

	respV4, errV4 := http.Get("https://api.ipify.org/")

	if errV4 != nil {
		ipV4OrErrorMsg = fmt.Sprintf("ERROR: %v", errV4)
	} else {
		defer respV4.Body.Close()

		body, err := ioutil.ReadAll(respV4.Body)
		if err != nil {
			ipV4OrErrorMsg = fmt.Sprintf("ERROR: %v", err)
		} else {
			ipV4OrErrorMsg = string(body)
		}
	}
	fmt.Printf("%v", ipV4OrErrorMsg)
	fmt.Println()
}

func detectIPv6(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	var ipV6OrErrorMsg = ""

	respV6, errV6 := http.Get("https://api6.ipify.org/")

	if errV6 != nil {
		ipV6OrErrorMsg = fmt.Sprintf("ERROR: %v", errV6)
	} else {
		defer respV6.Body.Close()

		body, err := ioutil.ReadAll(respV6.Body)
		if err != nil {
			ipV6OrErrorMsg = fmt.Sprintf("ERROR: %v", err)
		} else {
			ipV6OrErrorMsg = string(body)
		}
	}
	fmt.Printf("%v", ipV6OrErrorMsg)
	fmt.Println()
}

func detectIP(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	detectIPv4(args, flags)
	detectIPv6(args, flags)
}

func Setup_ip_Command() {
	commando.
		Register("ip").
		SetShortDescription("detect public IP address(es)").
		SetDescription("Tries to detect public IPv4 and IPv6 address(es) by using ipify.org service.").
		SetAction(detectIP)
}
