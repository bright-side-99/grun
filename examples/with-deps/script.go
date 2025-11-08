package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	// Demonstrate using an external dependency
	color.Cyan("=== grun Example with Dependencies ===")
	fmt.Println()

	color.Green("✓ Successfully imported and used github.com/fatih/color")
	color.Yellow("✓ This script requires go.mod to work")
	color.Magenta("✓ grun automatically detects go.mod and builds accordingly")
	
	fmt.Println()
	
	if len(os.Args) > 1 {
		color.Blue("Arguments received: %v", os.Args[1:])
	} else {
		color.White("Try running: grun script.go arg1 arg2")
	}
}

