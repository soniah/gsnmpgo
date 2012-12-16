package main

import (
	"fmt"

	g "github.com/soniah/gsnmp"
)

func main() {
	uri := `snmp://public@192.168.1.10//1.3.6.1.2.1.1.*`
	parsed_uri := g.ParseURI(uri)
	fmt.Println("parsed_uri:", parsed_uri)

	// if (! gnet_snmp_parse_path(uri->path, &vbl, &type, &error)) {
	vbl, _type, ok := g.ParsePath(parsed_uri)
	if ok {
		fmt.Printf("ok: %t, vbl: %s, _type: %s\n", ok, vbl, _type)
	} else {
		fmt.Println("fail! ParsePath")
	}
}
