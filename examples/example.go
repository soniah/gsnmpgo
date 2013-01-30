package main

// Copyright 2013 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

import (
	"fmt"
	"os"

	"github.com/soniah/gsnmp"
)

func main() {
	// home
	//uri := `snmp://public@192.168.1.10/`
	//uri := `snmp://public@192.168.1.10//1.3.6.1.2.1.1.*`
	//uri := `snmp://public@192.168.1.10//1.3.6.1.2.1.1.1.0`
	uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0)` // string, oid

	// work

	parsed_uri, err := gsnmp.ParseURI(uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("ParseURI():", parsed_uri)

	vbl, uritype, err := gsnmp.ParsePath(uri, parsed_uri)
	defer gsnmp.UriDelete(parsed_uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("ParsePath(): {oids: %s, uritype: %s, err: %v}\n", gsnmp.OidToString(vbl), uritype, err)

	session, err := gsnmp.NewUri(uri, parsed_uri)
	// defer gsnmp.UriDelete(parsed_uri) already setup
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("session:", session)

	switch gsnmp.UriType(uritype) {
	case gsnmp.GNET_SNMP_URI_GET:
		fmt.Println("doing GNET_SNMP_URI_GET")
		results, err := gsnmp.Get(session, vbl)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//fmt.Println(results) // hack
		gsnmp.Dump(results)
	case gsnmp.GNET_SNMP_URI_NEXT:
		fmt.Println("doing GNET_SNMP_URI_NEXT")
	case gsnmp.GNET_SNMP_URI_WALK:
		fmt.Println("doing GNET_SNMP_URI_WALK")
	}
}
