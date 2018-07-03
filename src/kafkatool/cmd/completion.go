// Copyright (C) 2018, CERT Gouvernemental (GOVCERT.LU)
// Author: Daniel Struck <daniel.struck@govcert.etat.lu>

package cmd

import (
	"kafkatool/helper"
	"os"

	"github.com/spf13/cobra"
)

// brokerCmd represents the broker command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate the bash completion file (kafkatool_completion.sh)",
	Run: func(cmd *cobra.Command, args []string) {
		// for the bash completion generation
		file, err := os.Create("kafkatool_completion.sh")
		helper.Check(err)
		rootCmd.GenBashCompletion(file)
		defer file.Close()
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)

}
