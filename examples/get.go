package main

// Copyright 2013 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

import (
	"fmt"
	"os"

	"github.com/soniah/gsnmpgo"
)

func main() {
	// home

	// GET: string, oid, timeticks, integer32, gauge32, hex-string, ipaddress, counter32
	uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0,1.3.6.1.2.1.1.3.0,1.3.6.1.2.1.1.7.0,1.3.6.1.2.1.2.2.1.5.6,1.3.6.1.2.1.2.2.1.6.1,1.3.6.1.2.1.4.20.1.1.192.168.1.10,1.3.6.1.2.1.2.2.1.10.1)`

	// WALK:
	// uri := `snmp://public@192.168.1.10//1.3.6.1.2.1.*`

	// NEXT:
	// uri := `snmp://public@192.168.1.10//1.3.6.1.2.1+`

	// simulator 127.0.0.1
	// string, counter64, hex string
	// uri := `snmp://public@127.0.0.1:162//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.31.1.1.1.8.19,1.3.6.1.2.1.3.1.1.2.14.1.10.0.0.1)`

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
