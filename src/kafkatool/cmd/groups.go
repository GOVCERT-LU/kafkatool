// Copyright (C) 2021, CERT Gouvernemental (GOVCERT.LU)
// Author: Daniel Struck <daniel.struck@govcert.etat.lu>

package cmd

import (
	"github.com/spf13/cobra"
)

// topicCmd represents the topic command
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Consumer groups related commands",
}

func init() {
	rootCmd.AddCommand(groupsCmd)
}
