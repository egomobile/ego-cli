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
	"net"

	"github.com/thatisuday/commando"
)

func localip_run(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	ipOnly, err := flags["ip-only"].GetBool()
	outputIPAddressesOnly := err == nil && ipOnly

	ifaces, err := net.Interfaces()

	if err != nil {
		log.Fatal((err))
	}

	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err == nil {
			var ipAddreses []string
			var addAddrIfNeeded = func(a string) {
				for _, b := range ipAddreses {
					if b == a {
						return
					}
				}

				ipAddreses = append(ipAddreses, a)
			}

			for _, addr := range addrs {
				var ip net.IP

				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}

				if !ip.IsLoopback() {
					if ip.To4() != nil {
						addAddrIfNeeded(ip.To4().String())
					}

					if ip.To16() != nil {
						addAddrIfNeeded(ip.To16().String())
					}
				}
			}

			if len(ipAddreses) > 0 {
				if !outputIPAddressesOnly {
					fmt.Println(i.Name)
				}

				for _, ipAddr := range ipAddreses {
					if outputIPAddressesOnly {
						fmt.Println(ipAddr)
					} else {
						fmt.Println("\t", ipAddr)
					}
				}
			}
		}
	}
}

func Setup_localip_Command() {
	commando.
		Register("local-ip").
		SetShortDescription("lists local net addresses").
		SetDescription("Lists all IP addresses of all known network interfaces").
		AddFlag("ip-only,i", "IP addresses only", commando.Bool, nil).
		SetAction(localip_run)
}
