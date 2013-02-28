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
	"github.com/soniah/gsnmpgo"
)

func main() {
	gsnmpgo.Debug = true

	/*
		.1.3.6.1.2.1.1.1.0 "Samsung ML-2850 Series OS 1.01.12.37 11-03-2008;Engine 1.01.23;NIC V4.01.03(ML-285x) 08-21-2008;S/N 4F50BAGS500082F "
		.1.3.6.1.2.1.1.2.0 .1.3.6.1.4.1.236.11.5.1
		.1.3.6.1.2.1.1.3.0 15:16:07:27.00
		.1.3.6.1.2.1.1.4.0 "Administrator"
		.1.3.6.1.2.1.1.5.0 "SEC00159937762B"
		.1.3.6.1.2.1.1.6.0 ""
		.1.3.6.1.2.1.1.7.0 104
	*/

	oids := []string{"1.3.6.1.2.1.1.1.0", "1.3.6.1.2.1.1.2.0", "1.3.6.1.2.1.1.3.0", "1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.5.0", "1.3.6.1.2.1.1.6.0", "1.3.6.1.2.1.1.7.0"}

	params := &gsnmpgo.QueryParams{
		Community: "public",
		IPAddress: "192.168.1.10",
		Version:   gsnmpgo.GNET_SNMP_V2C,
		Oids:      oids,
	}

	params.GetMany()
	results := params.Tree

	gsnmpgo.Dump(results)
	fmt.Println()

	if results == nil {
		fmt.Println("IterAscend: results are NIL - exiting")
		return
	}

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