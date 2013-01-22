package main

import (
	"fmt"
	"os"

	"github.com/soniah/gsnmp"
)

func main() {
	// home
	//uri := `snmp://public@192.168.1.10/`
	uri := `snmp://public@192.168.1.10//1.3.6.1.2.1.1.*`
	//uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0)`

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

	session, err := gsnmp.SnmpNewUri(uri, parsed_uri)
	fmt.Println("session:", session)
}
