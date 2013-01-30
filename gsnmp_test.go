package gsnmp

// gsnmp is a Go wrapper around the C gsnmp library.
//
// Copyright (C) 2013 Sonia Hamilton sonia@snowfrog.get.
//
// gsnmp is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// gsnmp is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser Public License for more details.
//
// You should have received a copy of the GNU Lesser Public License
// along with gsnmp.  If not, see <http://www.gnu.org/licenses/>.

// TODO these tests are hardcoded against public@192.168.1.10. How to make the
// tests more generic; writing mocks for snmp would be timeconsuming and
// painful. Packet captures??

import (
	"testing"
)

var oidAsStringTests = []struct {
	in  []int
	out string
	ok  bool
}{
	{[]int{1, 3, 6, 1, 4, 1, 2680, 1, 2, 7, 3, 2, 0}, ".1.3.6.1.4.1.2680.1.2.7.3.2.0", true},
	{[]int{}, "", true},
}

func TestNewObjectIdentifier(t *testing.T) {
	for i, test := range oidAsStringTests {
		ret := OidAsString(test.in)
		if test.ok && ret != test.out {
			t.Errorf("#%d: Bad result: %v (expected %v)", i, ret, test.out)
		}
	}
}

var parseURIs = []struct {
	uri string
	err bool
}{
	{`snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0)`, false},
	{`xyzz://public@192.168.1.10//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0)`, true},
}

func TestParseURI(t *testing.T) {
	for i, test := range parseURIs {
		_, ret_err := ParseURI(test.uri)
		if (ret_err != nil) != test.err {
			t.Errorf("#%d: Bad result: %v (expected %v)", i, ret_err, test.err)
		}
	}
}
