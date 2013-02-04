package gsnmpgo

// Copyright 2013 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

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

func TestOidAsString(t *testing.T) {
	for i, test := range oidAsStringTests {
		ret := OidAsString(test.in)
		if test.ok && ret != test.out {
			t.Errorf("#%d: Bad result: %v (expected %v)", i, ret, test.out)
		}
	}
}
