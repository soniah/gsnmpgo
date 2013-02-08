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
	gsnmpgo.Debug = true

	// GET
	// uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0)`

	// WALK:
	// uri := `snmp://public@192.168.1.10//1.3.6.1.*`

	// NEXT:
	// uri := `snmp://public@192.168.1.10//1.3.6.1.2.1+`

	// Verax GET - string, oid, timeticks
	uri := `snmp://public@127.0.0.1:161//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0,1.3.6.1.2.1.1.3.0)`

	// Verax NEXT:
	// uri := `snmp://public@127.0.0.1:161//(1.3.6.1.2.1.1.1.0)+`

	// Verax WALK
	// uri := `snmp://public@127.0.0.1:161//1.3.6.1.*`

	params := &gsnmpgo.QueryParams{
		Uri:     uri,
		Version: gsnmpgo.GNET_SNMP_V2C,
		Timeout: 1000,
		Retries: 5,
	}
	results, err := gsnmpgo.Query(params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gsnmpgo.Dump(results)
	fmt.Println()

	ch := results.IterAscend()
	for {
		r := <-ch
		if r == nil {
			break
		}
		result := r.(gsnmpgo.QueryResult)
		switch result.Value.(type) {
		case gsnmpgo.VBT_OctetString:
			fmt.Printf("OID %s is an octet string: %s\n", result.Oid, result.Value)
		default:
			fmt.Printf("OID %s is some other type\n", result.Oid)
		}
	}
	fmt.Println()

	ch2 := results.IterAscend()
	for {
		r := <-ch2
		if r == nil {
			break
		}
		result := r.(gsnmpgo.QueryResult)
		fmt.Printf("OID %s type: %T\n", result.Oid, result.Value)
		fmt.Printf("OID %s as a number: %d\n", result.Oid, result.Value.Integer())
		fmt.Printf("OID %s as a string: %s\n\n", result.Oid, result.Value)
	}
}
