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

func ip_detectIPv4(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
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

func ip_detectIPv6(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
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

func ip_detectIP(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	ip4, err := flags["ip4"].GetBool()
	detectV4 := err != nil || ip4

	ip6, err := flags["ip6"].GetBool()
	detectV6 := err != nil || ip6

	if !detectV4 && !detectV6 {
		detectV4 = true
		detectV6 = true
	}

	if detectV4 {
		ip_detectIPv4(args, flags)
	}

	if detectV6 {
		ip_detectIPv6(args, flags)
	}
}

func Setup_publicip_Command() {
	commando.
		Register("public-ip").
		SetShortDescription("detect public IP address(es)").
		SetDescription("Tries to detect public IPv4 and IPv6 address(es) by using ipify.org service").
		AddFlag("ip4,4", "detect IPv4", commando.Bool, nil).
		AddFlag("ip6,6", "detect IPv6", commando.Bool, nil).
		SetAction(ip_detectIP)
}
