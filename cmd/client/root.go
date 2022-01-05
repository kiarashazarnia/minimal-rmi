package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "minirmi",
	Short: "Mini RMI is a minimal yet powerful RMI client",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func ExecuteCommand() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
