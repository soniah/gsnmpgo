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

import (
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	"strconv"
	"testing"
)

var _ = fmt.Sprintf("dummy") // dummy
var _ = strconv.Itoa(0)      // dummy

var oidAsStringTests = []struct {
	in  []int
	out string
	ok  bool
}{
	{[]int{1, 3, 6, 1, 4, 1, 2680, 1, 2, 7, 3, 2, 0}, ".1.3.6.1.4.1.2680.1.2.7.3.2.0", true},
	{[]int{}, "", true},
}

func TestOidAsString(t *testing.T) {
	for i, test := range oidAsStringTests {
		ret := OidAsString(test.in)
		if test.ok && ret != test.out {
			t.Errorf("#%d: Bad result: %v (expected %v)", i, ret, test.out)
		}
	}
}

var lessOIDTests = []struct {
	oid_a string
	oid_b string
	less  bool
}{
	{"1.2.3", "1.2.4", true},
	{"1.2.3", "1.2.3.4", true},
	{"1.2.3", "1.2", false},
	{"", "1.2", true},
	{"1.2", "", false},
	{"", "", false},
	{"1.9", "1.10", true},
	{"1.12", "1.10", false},
	{"1.12.2", "1.12.1", false},
	{"1.12.1", "1.12.2", true},
}

func TestLessOID(t *testing.T) {
	for i, test := range lessOIDTests {
		astruct := QueryResult{Oid: test.oid_a}
		bstruct := QueryResult{Oid: test.oid_b}
		if res := lessOID(astruct, bstruct); res != test.less {
			t.Errorf("#%d: expected (%t) got (%t) oid_a (%s) oid_b (%s)",
				i, test.less, res, test.oid_a, test.oid_b)
		}
	}
}

var veraxDevices = []struct {
	path string
	port int
}{
	{"testing/device/os/os-linux-std.txt", 161},
	{"testing/device/cisco/cisco_router.txt", 162},
}

func TestQueryGets(t *testing.T) {
	for i, test := range veraxDevices {
		var err error

		var vresults *llrb.Tree
		if vresults, err = ReadVeraxResults(test.path); err != nil {
			t.Errorf("#%d, %s: ReadVeraxResults error: %s", i, test.path, err)
		}

		var counter int
		var uri, oids string
		gresults := llrb.New(lessOID)
		ch := vresults.IterAscend()
		for {
			r := <-ch
			if r == nil {
				break
			}
			counter += 1
			if counter <= 3 { // TODO random number of oids, not hardcoded n
				oid := r.(QueryResult)
				oids += "," + oid.Oid
			} else {
				uri = `snmp://public@127.0.0.1:` + strconv.Itoa(test.port) + "//(" + oids[1:] + ")"
				counter = 0 // reset
				oids = ""   // reset

				params := &QueryParams{
					Uri:     uri,
					Version: GNET_SNMP_V2C,
					Timeout: 200,
					Retries: 5,
					Tree:    gresults,
				}
				if _, err := Query(params); err != nil {
					t.Errorf("Query error: %s. Uri: %s", err, uri)
				}
			}
		}
		CompareVerax(t, gresults, vresults)
	}
}

var partitionAllPTests = []struct {
	cp int
	ps int
	sl int
	ok bool
}{
	{-1, 3, 8, false}, // test out of range
	{8, 3, 8, false},  // test out of range
	{0, 3, 8, false},  // test 0-7/3 per doco
	{1, 3, 8, false},
	{2, 3, 8, true},
	{3, 3, 8, false},
	{4, 3, 8, false},
	{5, 3, 8, true},
	{6, 3, 8, false},
	{7, 3, 8, true},
}

func TestPartitionAllP(t *testing.T) {
	for i, test := range partitionAllPTests {
		ok := PartitionAllP(test.cp, test.ps, test.sl)
		if ok != test.ok {
			t.Errorf("#%d: Bad result: %v (expected %v)", i, ok, test.ok)
		}
	}
}
