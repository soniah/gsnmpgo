/*
Package gsnmpgo is a go/cgo wrapper around gsnmp; it currently provides support
for snmp v1 and v2c, and snmp get, snmp getnext, and snmp walk.

gsnmpgo is pre 1.0, therefore API's may change, and tests are minimal.

INSTALLATION

gsnmpgo requires the following libraries, as well as library header files:

    glib2.0, gsnmp, gnet-2.0

Here is an example of installation on Ubuntu Precise 12.04.1:

    # setup Go
    sudo aptitude install golang git
    cat >> ~/.bashrc
    export GOPATH="${HOME}/go"
    ^D
    mkdir ~/go && source ~/.bashrc && cd ~/go

    # download only - troubleshooting builds is easier
    go get -d github.com/soniah/gsnmpgo

    # install prerequisites for gsnmpgo and build
    sudo aptitude install libglib2.0-dev libgsnmp0-dev libgnet-dev
    go install github.com/soniah/gsnmpgo

    # test working (you will need to edit example.go and
    # provide different uris)
    cd src/github.com/soniah/gsnmpgo/examples
    go run example.go

SUMMARY

Here is a summary of usage:

    // do an snmp get; RFC 4088 is used for uris
    uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0)`
    params := &gsnmpgo.QueryParams{
        Uri:     uri,
        Version: gsnmpgo.GNET_SNMP_V2C,
        Timeout: 200,
        Retries: 2,
    }
    results, err := gsnmpgo.Query(params)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // use results; result.Value has an interface that supports Stringer and Integer()
    for oid, value := range results {
        fmt.Printf("oid, type: %s, %T\n", oid, value)
        fmt.Printf("INTEGER: %d\n", value.Integer())
        fmt.Printf("STRING : %s\n", value)
        fmt.Println()
    }

    // or if you just want to check your results, use Dump()
    gsnmpgo.Dump(results)

    // turn on debugging
    gsnmpgo.Debug = true

SPECIFYING URIS

http://tools.ietf.org/html/rfc4088 is used for specifying snmp uris; as well
as doing an snmp get you can also do a snmp getnext or snmp walk:

    // GET - notice you can have multiple OIDs
    uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0)`

    // NEXT - notice the plus sign at the end
    uri := `snmp://public@192.168.1.10//1.3.6.1.2.1+`

    // WALK - notice the star at the end
    // uri := `snmp://public@192.168.1.10//1.3.6.1.2.1.*`

RESULTS

The results are returned as a map of oid to Varbinder:

    type QueryResults map[string]Varbinder

If you want access to the snmp type for each result returned, you could
use a type switch:

    for oid, value := range results {
        switch value.(type) {
        case gsnmpgo.VBT_OctetString:
            fmt.Printf("OID %s is an octet string: %s\n", oid, value)
        default:
            fmt.Printf("OID %s is some other type\n", oid)
        }
    }

Often you just want the result as a string or a number, Varbinder is
an interface that provides two convenience functions:

    type Varbinder interface {
        // Integer() needs to handle both signed numbers (int32), as well as
        // unsigned int 64 (uint64). Therefore it returns a *big.Int.
        Integer() *big.Int
        fmt.Stringer
    }

    for oid, value := range results {
        fmt.Printf("OID %s as a number: %d\n", oid, value.Integer())
        fmt.Printf("OID %s as a string: %s\n", oid, value)
    }

Some of the Stringers are smart, for example gsnmpgo.VBT_Timeticks will be
formatted as days, hours, etc when returned as a string:

    OID 1.3.6.1.2.1.1.3.0 as a number: 4381200
    OID 1.3.6.1.2.1.1.3.0 as a string: 0 days, 12:10:12.00

Sonia Hamilton, sonia@snowfrog.net, http://www.snowfrog.net.
*/
package gsnmpgo

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
