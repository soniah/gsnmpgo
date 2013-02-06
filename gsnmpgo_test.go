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
	"strconv"
	"testing"
)

var _ = fmt.Sprintf("dummy") // dummy

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

var veraxDevices = []struct {
	path string
	port int
}{
	{"testing/device/os/os-linux-std.txt", 161},
	{"testing/device/cisco/cisco_router.txt", 162},
}

func TestQuery(t *testing.T) {
	for i, test := range veraxDevices {
		var err error
		var vresults, gresults QueryResults

		if vresults, err = ReadVeraxResults(test.path); err != nil {
			t.Errorf("#%d, %s: ReadVeraxResults error: %s", i, test.path, err)
		}
		fmt.Println(vresults) // dummy

		uri := `snmp://public@127.0.0.1:` + strconv.Itoa(test.port) + `//1.3.6.1.*`
		if gresults, err = Query(uri, GNET_SNMP_V2C); err != nil {
			t.Errorf("#%d, %s: Query error: %s", i, uri, err)
		}
		fmt.Println(gresults) // dummy

	}
}
