package gsnmp

// Copyright 2012 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

// TODO some of these tests are hardcoded against public@192.168.1.10
// How to make the tests more generic; writing mocks for snmp would be
// timeconsuming and painful. Packet captures??

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
