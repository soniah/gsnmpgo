// Copyright 2012 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
package enumconv

import (
	"fmt"
	"strings"
)

// Write uses gocog to produce Go boilerplate and Stringers for C enums
//
// gotypename: the go type name of this C enum
// ctypename: the C type name of this enum (without "_Ctype_")
// enums: slice of strings containing names of enums
// ccode: any C code to be included in Stringer comment header
//
func Write(gotypename string, ctypename string, enums []string, ccode string) {

	// type
	fmt.Printf("\n// type and values for %s\n//\n", ctypename)
	fmt.Printf("type %s int\n", gotypename)
	fmt.Println()

	// const
	fmt.Println("const (")
	for i, enum := range enums {
		if i == 0 {
			fmt.Printf("	%s %s = iota\n", enum, gotypename)
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
	fmt.Println("//")

	// start stringer function
	fmt.Printf("func (%s %s) String() string {\n", strings.ToLower(gotypename), "_Ctype_" + ctypename)

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

	// for line before [[[end]]]
	//fmt.Println()
}
