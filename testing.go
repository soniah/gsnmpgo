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
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var _ = fmt.Sprintf("dummy")        // dummy
var _ = strings.Split("dummy", "m") // dummy
var _ = strconv.Itoa(0)             // dummy

func ReadVeraxResults(filename string) (results *llrb.Tree, err error) {
	var lines []byte
	if lines, err = ioutil.ReadFile(filename); err != nil {
		return nil, fmt.Errorf("unable to open file %s", filename)
	}
	results = llrb.New(lessOID)

	// some lines have newlines in them, therefore can't just split on newline
	lines_split := re_split(regexp.MustCompile(`\n\.`), string(lines), -1)
LINE:
	for _, line := range lines_split {
		splits_a := strings.SplitN(line, " = ", 2)
		oid := splits_a[0]
		splits_b := strings.SplitN(splits_a[1], ": ", 2)
		oidtype := splits_b[0]
		oidval := strings.TrimSpace(splits_b[1])

		// removing leading . first oid
		if string(oid[0]) == "." {
			oid = oid[1:]
		}

		var value Varbinder
		switch oidtype {

		case "STRING", "String", "Hex-STRING":
			oidval = strings.Trim(oidval, `"`)
			value = VBT_OctetString(oidval)

		case "OID":
			value = VBT_ObjectID(oidval)

		case "IpAddress", "Network Address":
			value = VBT_IPAddress(oidval)

		case "INTEGER":
			if n, err := strconv.Atoi(oidval); err == nil {
				value = VBT_Integer32(n)
			} else {
				panic(fmt.Sprintf("Err converting integer. oid: %s err: %v\n", oid, err))
			}

		case "Gauge32":
			if n, err := strconv.Atoi(oidval); err == nil {
				value = VBT_Unsigned32(n)
			}

		case "Counter32":
			if n, err := strconv.ParseUint(oidval, 10, 32); err == nil {
				value = VBT_Counter32(n)
			} else {
				panic(fmt.Sprintf("Counter32: oid: %s oidval: %s err: %v\n", oid, oidval, err))
			}

		case "Timeticks":
			matches := regexp.MustCompile(`\d+`).FindAllString(oidval, 1) // pull out "(value)"
			oidval := matches[0]
			if n, err := strconv.Atoi(oidval); err == nil {
				value = VBT_Timeticks(n)
			}

		case "Counter64":
			if n, err := strconv.ParseUint(oidval, 10, 64); err == nil {
				value = VBT_Counter64(n)
			}

		case "BITS":
			continue LINE
			// TODO is BITS Verax specific, or doesn't gsnmp handle?
			// .1.3.6.1.2.1.88.1.4.2.1.3.6.95.115.110.109.112.100.95.109.116.101.
			// 84.114.105.103.103.101.114.70.105.114.101.100
			// = BITS: 38 30 20 30 2 3 4 10 11 18 26 27

		default:
			panic(fmt.Sprintf("Unhandled type: %s, %s\n", oidtype, oidval))
		}
		result := QueryResult{Oid: oid, Value: value}
		results.ReplaceOrInsert(result)
	}
	return results, nil
}

func CompareVerax(t *testing.T, gresults, vresults *llrb.Tree) {
	ch := gresults.IterAscend()
	for {
		gr := <-ch
		if gr == nil {
			break
		}
		goresult := gr.(QueryResult)
		vstruct := QueryResult{Oid: goresult.Oid}
		vr := vresults.Get(vstruct)
		if vr == nil {
			continue
		}
		vresult := vr.(QueryResult)

		vstring := fmt.Sprintf("%s", vresult.Value)
		gostring := fmt.Sprintf("%s", goresult.Value)
		if gostring != vstring {
			// fmt.Printf("OK oid: %s type: %T value: %s\n", goresult.Oid, goresult.Value, gostring)
			if len(gostring) > 4 && gostring[0:5] == "07 DA" {
				// skip - weird Verax stuff
			} else if len(vstring) > 4 && vstring[0:5] == "4E:85" {
				// skip - weird Verax stuff
			} else if len(vstring) > 17 && vstring[0:18] == "Cisco IOS Software" {
				// skip - \n's have been stripped - ignore
			} else {
				t.Errorf("compare fail: oid: %s type: %T\ngostring: |%s|\nvstring : |%s|",
					goresult.Oid, goresult.Value, gostring, vstring)
			}
		}
	}
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
