#!/usr/bin/env grun
package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	color.Cyan("╔═══════════════════════════════════════╗")
	color.Cyan("║  Executable Go Script with Dependencies  ║")
	color.Cyan("╚═══════════════════════════════════════╝")
	fmt.Println()
	
	color.Green("✓ Shebang: #!/usr/bin/env grun")
	color.Green("✓ External dependencies via go.mod")
	color.Green("✓ Executable: chmod +x color-demo.go")
	color.Green("✓ Run directly: ./color-demo.go")
	
	fmt.Println()
	
	if len(os.Args) > 1 {
		color.Yellow("📦 Arguments: %v", os.Args[1:])
	} else {
		color.White("💡 Try: ./color-demo.go arg1 arg2")
	}
	
	fmt.Println()
	color.Magenta("🚀 Go scripts, made easy!")
}

