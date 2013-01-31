package gsnmpgo

// Copyright 2013 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

// common.go contains "common" or miscellanous functions

import (
	"fmt"
	"strings"
)

// AsString returns the string representation of an Oid
func OidAsString(o []int) string {
	if len(o) == 0 {
		return ""
	}
	result := fmt.Sprintf("%v", o)
	result = result[1 : len(result)-1] // strip [ ] of Array representation
	return "." + strings.Join(strings.Split(result, " "), ".")
}
