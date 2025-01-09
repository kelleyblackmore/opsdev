package main

import (
	"fmt"
	"os"

	"github.com/kelleyblackmore/opsdev/internal/installer"
)

func main() {
	installer := installer.NewToolInstaller()
	if err := installer.StartSetup(); err != nil {
		fmt.Printf("Error during setup: %v\n", err)
		os.Exit(1)
	}
}