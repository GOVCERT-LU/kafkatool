// Copyright (C) 2018, CERT Gouvernemental (GOVCERT.LU)
// Author: Daniel Struck <daniel.struck@govcert.etat.lu>

package cmd

import (
	"github.com/spf13/cobra"
)

// brokerCmd represents the broker command
var brokerCmd = &cobra.Command{
	Use:   "broker",
	Short: "Commands related to the kafka brokers",
}

func init() {
	rootCmd.AddCommand(brokerCmd)

}
