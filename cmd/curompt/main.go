package main

import (
	"fmt"
	"os"

	"github.com/curogom/curompt/internal/cli"
)

var (
	// Version indicates the build version embedded at compile time.
	Version = "dev"
	// BuildTime is the UTC timestamp when the binary was built.
	BuildTime = "unknown"
	// GitCommit is the short SHA corresponding to the build.
	GitCommit = "unknown"
)

func main() {
	rootCmd := cli.NewRootCmd(Version, BuildTime, GitCommit)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
