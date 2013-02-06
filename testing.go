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
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var _ = fmt.Sprintf("dummy")        // dummy
var _ = strings.Split("dummy", "m") // dummy

func ReadVeraxResults(filename string) (results QueryResults, err error) {
	var lines []byte
	if lines, err = ioutil.ReadFile(filename); err != nil {
		return nil, fmt.Errorf("unable to open file %s", filename)
	}
	results = make(QueryResults)

	// some lines have newlines in them, therefore can't just split on newline
	lines_split := re_split(regexp.MustCompile(`\r\n\.`), string(lines), -1)
	for _, line := range lines_split {
		splits_a := strings.SplitN(line, " = ", 2)
		oid := splits_a[0]
		splits_b := strings.SplitN(splits_a[1], ": ", 2)
		oidtype := splits_b[0]
		oidval := splits_b[1]

		// removing leading . first oid
		if string(oid[0]) == "." {
			oid = oid[1:]
		}

		var value Varbinder
		switch oidtype {

		case "STRING", "String", "Hex-STRING":
			value = VBT_OctetString(oidval)

		case "OID":
			value = VBT_ObjectID(oidval)

		case "IpAddress":
			value = VBT_IPAddress(oidval)

		case "INTEGER":
			if n, err := strconv.Atoi(oidval); err != nil {
				value = VBT_Integer32(n)
			}

		case "Gauge32":
			if n, err := strconv.Atoi(oidval); err != nil {
				value = VBT_Unsigned32(n)
			}

		case "Counter32":
			if n, err := strconv.Atoi(oidval); err != nil {
				value = VBT_Counter32(n)
			}

		case "Timeticks":
			if n, err := strconv.Atoi(oidval); err != nil {
				value = VBT_Timeticks(n)
			}

		case "Counter64":
			if n, err := strconv.Atoi(oidval); err != nil {
				value = VBT_Counter64(n)
			}

		case "Network Address":
			// ?? Network Address, C0:A8:68:01

		case "BITS":
			// ?? BITS, 80 0

		default:
			fmt.Printf("Unhandled type: %s, %s\n", oidtype, oidval)
		}
		results[oid] = value
	}
	return results, nil
}

// adapted from http://codereview.appspot.com/6846048/
//
// re_split slices s into substrings separated by the expression and returns a slice of
// the substrings between those expression matches.
//
// The slice returned by this method consists of all the substrings of s
// not contained in the slice returned by FindAllString(). When called on an exp ression
// that contains no metacharacters, it is equivalent to strings.SplitN().
// Example:
// s := regexp.MustCompile("a*").re_split("abaabaccadaaae", 5)
// // s: ["", "b", "b", "c", "cadaaae"]
//
// The count determines the number of substrings to return:
// n > 0: at most n substrings; the last substring will be the unsplit remaind er.
// n == 0: the result is nil (zero substrings)
// n < 0: all substrings
func re_split(re *regexp.Regexp, s string, n int) []string {
	if n == 0 {
		return nil
	}
	if len(s) == 0 {
		return []string{""}
	}
	matches := re.FindAllStringIndex(s, n)
	strings := make([]string, 0, len(matches))
	beg := 0
	end := 0
	for _, match := range matches {
		if n > 0 && len(strings) >= n-1 {
			break
		}
		end = match[0]
		if match[1] != 0 {
			strings = append(strings, s[beg:end])
		}
		beg = match[1]
	}
	if end != len(s) {
		strings = append(strings, s[beg:])
	}
	return strings
}
