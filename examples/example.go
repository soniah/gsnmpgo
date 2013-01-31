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
	//uri := `snmp://public@192.168.1.10/`
	//uri := `snmp://public@192.168.1.10//1.3.6.1.2.1.1.*`
	//uri := `snmp://public@192.168.1.10//1.3.6.1.2.1.1.1.0`

	// string, oid, timeticks, integer32, gauge32, hex-string, ipaddress, counter32
	//uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0,1.3.6.1.2.1.1.3.0,1.3.6.1.2.1.1.7.0,1.3.6.1.2.1.2.2.1.5.6,1.3.6.1.2.1.2.2.1.6.1,1.3.6.1.2.1.4.20.1.1.192.168.1.10,1.3.6.1.2.1.2.2.1.10.1)`

	// simulator 127.0.0.1
	// string, counter64, hex string
	uri := `snmp://public@127.0.0.1:162//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.31.1.1.1.8.19,1.3.6.1.2.1.3.1.1.2.14.1.10.0.0.1)`

	// work

	parsed_uri, err := gsnmpgo.ParseURI(uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("ParseURI():", parsed_uri)

	vbl, uritype, err := gsnmpgo.ParsePath(uri, parsed_uri)
	defer gsnmpgo.UriDelete(parsed_uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("ParsePath(): {oids: %s, uritype: %s, err: %v}\n", gsnmpgo.GListOidsString(vbl), uritype, err)

	session, err := gsnmpgo.NewUri(uri, gsnmpgo.GNET_SNMP_V2C, parsed_uri)
	// defer gsnmpgo.UriDelete(parsed_uri) is already setup
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("session:", session)

	switch gsnmpgo.UriType(uritype) {
	case gsnmpgo.GNET_SNMP_URI_GET:
		fmt.Println("doing GNET_SNMP_URI_GET")
		results, err := gsnmpgo.Get(session, vbl)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		gsnmpgo.Dump(results)
	case gsnmpgo.GNET_SNMP_URI_NEXT:
		fmt.Println("doing GNET_SNMP_URI_NEXT")
	case gsnmpgo.GNET_SNMP_URI_WALK:
		fmt.Println("doing GNET_SNMP_URI_WALK")
	}
}
