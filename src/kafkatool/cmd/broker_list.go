// Copyright (C) 2018, CERT Gouvernemental (GOVCERT.LU)
// Author: Daniel Struck <daniel.struck@govcert.etat.lu>

package cmd

import (
	"fmt"
	"kafkatool/helper"

	"github.com/spf13/cobra"
)

// brokerListCmd represents the list command
var brokerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Kafka brokers",

	Run: func(cmd *cobra.Command, args []string) {

		client, controller := helper.ConnectKafkaClient()
		defer client.Close()

		fmt.Println()

		fmt.Printf("%-3s %-40s %s\n", "id", "addr", "controller")

		for _, broker := range client.Brokers() {

			fmt.Printf("%-3d %-40s %t\n", broker.ID(), broker.Addr(), broker.ID() == controller.ID())

		}

		fmt.Println()

	},
}

func init() {

	brokerCmd.AddCommand(brokerListCmd)

}
