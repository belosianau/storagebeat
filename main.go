package main

import (
	"os"

	"storagebeat/cmd"

	_ "storagebeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
