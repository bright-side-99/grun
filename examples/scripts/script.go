#!/usr/bin/env grun
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("=== Simple grun Example ===")
	fmt.Println()
	fmt.Println("This is a simple Go script with no external dependencies.")
	fmt.Println("It uses only the standard library and runs instantly with grun!")
	fmt.Println()
	
	if len(os.Args) > 1 {
		fmt.Printf("Arguments received: %s\n", strings.Join(os.Args[1:], ", "))
	} else {
		fmt.Println("Try running: grun script.go arg1 arg2")
	}
}
