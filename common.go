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
