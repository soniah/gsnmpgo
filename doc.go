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

    # install GoLLRB; ignore error about "no Go source files" - GoLLRB has a
    # non-standard layout
    go get -d github.com/petar/GoLLRB
    go install github.com/petar/GoLLRB/llrb

    # install gsnmpgo
    go get -d github.com/soniah/gsnmpgo
    sudo aptitude install libglib2.0-dev libgsnmp0-dev libgnet-dev
    go install github.com/soniah/gsnmpgo

    # test working (you will need to edit example.go and
    # provide different uris)
    cd src/github.com/soniah/gsnmpgo/examples
    go run example.go

SUMMARY

(most of this code is in examples/example.go)

    // do an snmp get; RFC 4088 is used for uris
    uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0)`
    params := gsnmpgo.NewDefaultParams(uri)
    results, err := gsnmpgo.Query(params)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // check your results
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

The results are returned as an LLRB tree to provide "ordered map"
functionality (ie finding items by "key", and iterating "in order"). Items in
the tree are of type QueryResult:

    type QueryResult struct {
        Oid   string
        Value Varbinder
    }

See http://github.com/petar/GoLLRB for more documentation on using the LLRB
tree.

SNMP types are represented by Go types that implement the Varbinder interface
(eg "Octet String" is VBT_OctetString, "IP Address" is VBT_IPAddress). Use a
type switch to make decisions based on the SNMP type:

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

The Varbinder interface has two convenience functions Integer() and String()
that allow you to get all your results "as a string" or "as a number":

    type Varbinder interface {
        Integer() *big.Int
        fmt.Stringer
    }

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

Some of the Stringers are smart, for example gsnmpgo.VBT_Timeticks will be
formatted as days, hours, etc when returned as a string:

    OID 1.3.6.1.2.1.1.3.0 as a number: 4381200
    OID 1.3.6.1.2.1.1.3.0 as a string: 0 days, 12:10:12.00

TESTS

The tests use the Verax Snmp Simulator [1]; setup Verax before running "go test":

* download, install and run Verax with the default configuration

* in the gsnmpgo/testing directory, setup these symlinks (or equivalents for your system):

    ln -s /usr/local/vxsnmpsimulator/device device
    ln -s /usr/local/vxsnmpsimulator/conf/devices.conf.xml devices.conf.xml

* remove randomising elements from Verax device files:

    cd testing/device/cisco
    sed -i -e 's!\/\/\$.*!!' -e 's!^M!!' cisco_router.txt
    cd ../os
    sed -i -e 's!\/\/\$.*!!' -e 's!^M!!' os-linux-std.txt

[1] http://www.veraxsystems.com/en/products/snmpsimulator

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
