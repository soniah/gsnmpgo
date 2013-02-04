package main

// gsnmpgo is a go/cgo wrapper around gsnmp.
//
// Copyright (C) 2012-2013 Sonia Hamilton sonia@snowfrog.net.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import (
	"fmt"
	"os"

	"github.com/soniah/gsnmpgo"
)

func main() {
	// GET
	uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0)`

	// WALK:
	// uri := `snmp://public@192.168.1.10//1.3.6.1.2.1.*`

	// NEXT:
	// uri := `snmp://public@192.168.1.10//1.3.6.1.2.1+`

	results, err := gsnmpgo.Query(uri, gsnmpgo.GNET_SNMP_V2C)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gsnmpgo.Dump(results)

	for _, result := range results {
		switch result.Value.(type) {
		case gsnmpgo.VBT_OctetString:
			fmt.Printf("result is a an octet string: %s\n", result)
		default:
			fmt.Println("result is some other type")
		}
	}

	for _, result := range results {
		fmt.Printf("OID %s as a number: %d\n", result.Oid, result.Value.Integer())
		fmt.Printf("OID %s as a string: %s\n", result.Oid, result.Value)
	}
}
