package main

import (
	"fmt"
	g "github.com/soniah/gsnmpgo"
)

func main() {
	fmt.Println("hello testing world!")
	results, _ := g.ReadVeraxResults("device/os/os-linux-std.txt")
	fmt.Printf("%v\n", results)
}
