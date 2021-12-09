// Copyright (C) 2018, CERT Gouvernemental (GOVCERT.LU)
// Author: Daniel Struck <daniel.struck@govcert.etat.lu>

package cmd

import (
	"fmt"
	"kafkatool/helper"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

var (
	groupNameDelete string
	groupYesPrompt  bool
)

// deleteCmd represents the delete command
var groupsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a topic",
	Run: func(cmd *cobra.Command, args []string) {

		kafkaClient := helper.ConnectKafkaClient()
		defer kafkaClient.Close()

		clusterAdmin, err := sarama.NewClusterAdminFromClient(kafkaClient)
		helper.Check(err)
		defer clusterAdmin.Close()

		fmt.Println()

		if yesPrompt || helper.Confirmation(fmt.Sprintf("Do you really want to delete the consumer group %s?", groupNameDelete)) {

			err := clusterAdmin.DeleteConsumerGroup(groupNameDelete)
			fmt.Println()
			helper.Check(err)

		} else {

			fmt.Println("\nnothing done")

		}

		fmt.Printf("consumer group %s deleted\n", groupNameDelete)

	},
}

func init() {

	groupsCmd.AddCommand(groupsDeleteCmd)

	groupsDeleteCmd.Flags().StringVar(&groupNameDelete, "name", "", "group name")
	groupsDeleteCmd.MarkFlagRequired("name")

	groupsDeleteCmd.Flags().BoolVar(&groupYesPrompt, "yes", false, "automatic yes to confirmation prompt")

}
