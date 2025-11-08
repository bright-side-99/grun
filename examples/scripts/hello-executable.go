#!/usr/bin/env grun
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("🚀 Executable Go Script!")
	fmt.Println()
	fmt.Println("This script can be run directly: ./hello-executable.go")
	fmt.Println("No need to type 'grun' - the shebang does it for you!")
	fmt.Println()
	
	if len(os.Args) > 1 {
		fmt.Printf("Args: %v\n", os.Args[1:])
	}
}

