// Package enumconv provides helper functions for gocog, used in
// gsnmpgo/stringers.go.
package enumconv

// Copyright 2013 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

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

	// function comment header
	fmt.Printf("// Stringer for %s\n//\n", ctypename)

	// C code include
	fmt.Println("// C:")
	for _, line := range strings.Split(ccode, "\n") {
		fmt.Printf("//    %s\n", line)
	}

	// start stringer function
	fmt.Printf("func (%s %s) String() string {\n", strings.ToLower(gotypename), ctypename)

	// switch statement
	fmt.Printf("	switch %s(%s) {\n", gotypename, strings.ToLower(gotypename))
	for _, enum := range enums {
		fmt.Printf("	case %s:\n", enum)
		fmt.Printf("		return \"%s\"\n", enum)
	}
	fmt.Println("	}")
	fmt.Printf("	return \"UNKNOWN %s\"\n", ctypename)

	// end stringer function
	fmt.Println("}")
}
