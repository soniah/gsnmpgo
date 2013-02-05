// Package enumconv provides helper functions for gocog, used in
// gsnmpgo/stringers.go.
package enumconv

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
	"strings"
)

// Write uses gocog to produce Go boilerplate and Stringers for C enums.
//
// gotypename: the go type name of this C enum
//
// ctypename: the C type name of this enum (with "_Ctype_")
//
// enums: slice of strings containing names of enums
//
// ccode: any C code to be included in Stringer comment header
//
// start_at: value to start the enum at
//
// I decided not to parse the C enum typedefs and instead opted to pass in
// fields like gotypename and ctypename, as parsing would be overkill for this
// project.
func Write(gotypename string, ctypename string, enums []string, ccode string, start_at int) {

	// type
	fmt.Printf("\n// type and values for %s\n", ctypename)
	fmt.Printf("type %s int\n", gotypename)
	fmt.Println()

	// const
	fmt.Println("const (")
	for i, enum := range enums {
		if i == 0 {
			if start_at == 0 {
				fmt.Printf("	%s %s = iota\n", enum, gotypename)
			} else {
				fmt.Printf("	%s %s = iota %+d\n", enum, gotypename, start_at)
			}
		} else {
			fmt.Printf("	%s\n", enum)
		}
	}
	fmt.Println(")")
	fmt.Println()

	// Go type Stringer
	///////////////////

	// function comment header
	fmt.Printf("// Stringer for %s\n", gotypename)

	// start stringer function
	receiver_name := strings.ToLower(gotypename)
	fmt.Printf("func (%s %s) String() string {\n", receiver_name, gotypename)

	// switch statement
	fmt.Printf("	switch %s {\n", receiver_name)
	for _, enum := range enums {
		fmt.Printf("	case %s:\n", enum)
		fmt.Printf("		return \"%s\"\n", enum)
	}
	fmt.Println("	}")
	fmt.Printf("	return \"UNKNOWN %s\"\n", gotypename)

	// end stringer function
	fmt.Println("}")
	fmt.Println()

	// C type Stringer
	//////////////////

	// function comment header
	fmt.Printf("// Stringer for %s\n//\n", ctypename)

	// C code include
	fmt.Println("// C:")
	for _, line := range strings.Split(ccode, "\n") {
		fmt.Printf("//    %s\n", line)
	}

	// start stringer function
	fmt.Printf("func (%s %s) String() string {\n", receiver_name, ctypename)

	// stringer body - use gotype stringer
	ins_string := `	fmt.Sprintf("%s", `
	fmt.Printf("	return %s %s(%s))\n", ins_string, gotypename, receiver_name)

	// end stringer function
	fmt.Println("}")
}
