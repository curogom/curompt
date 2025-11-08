package main

import (
	"fmt"
	"os"

	"github.com/curogom/curompt/internal/cli"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	rootCmd := cli.NewRootCmd(Version, BuildTime, GitCommit)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
