// envii is a TUI-first secrets / .env manager.
package main

import (
	"fmt"
	"os"

	"github.com/Trickster-ID/envii/internal/cli"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	if err := cli.Execute(version); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
