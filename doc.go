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

    // do an snmp get; RFC 4088 is used for uri's
    uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0)`
    results, err := gsnmpgo.Query(uri, gsnmpgo.GNET_SNMP_V2C)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // use results; result.Value has an interface that supports Stringer and Integer()
    for _, result := range results {
        fmt.Printf("%T:%s STRING:%s INTEGER:%d\n",
            result.Value, result.Oid, result.Value, result.Value.Integer())
    }

    // or if you just want to print your results, use Dump()
    gsnmpgo.Dump(results)

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

The results are returned as a slice of QueryResult:

    type QueryResults []QueryResult
    type QueryResult struct {
        Oid   string
        Value Varbinder
    }

If you want access to the snmp type for each result returned, you could
use a type switch:

    for _, result := range results {
        switch result.Value.(type) {
        case gsnmpgo.VBT_OctetString:
            fmt.Printf("result is a an octet string: %s\n", result)
        default:
            fmt.Println("result is some other type")
        }
    }

Often you just want the result as a string or a number, Varbinder is
an interface that provides two convenience functions:

    type Varbinder interface {
        Integer() int64
        fmt.Stringer
    }

    for _, result := range results {
        fmt.Printf("OID %s as a number: %d\n", result.Oid, result.Value.Integer())
        fmt.Printf("OID %s as a string: %s\n", result.Oid, result.Value)
    }

Some of the Stringers are smart, for example gsnmpgo.VBT_Timeticks will be
formatted as days, hours, etc when returned as a string:

    OID 1.3.6.1.2.1.1.3.0 as a number: 4381200
    OID 1.3.6.1.2.1.1.3.0 as a string: 0 days, 12:10:12.00

Sonia Hamilton, sonia@snowfrog.net, http://www.snowfrog.net.
*/
package gsnmpgo

// Copyright 2013 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.
