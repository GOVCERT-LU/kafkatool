// Copyright (C) 2018, CERT Gouvernemental (GOVCERT.LU)
// Author: Daniel Struck <daniel.struck@govcert.etat.lu>

package cmd

import (
	"fmt"
	"kafkatool/helper"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

var topicNameDelete string

// deleteCmd represents the delete command
var topicDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a topic",
	Run: func(cmd *cobra.Command, args []string) {

		client, controller := helper.ConnectKafkaClient()
		defer client.Close()

		topicRequest := &sarama.DeleteTopicsRequest{
			Version: 1,               // version 2 requires Kafka 1.0.0
			Timeout: time.Second * 2, // wait 2 seconds for the Kafka server to delete the topic
			Topics:  []string{strings.TrimSpace(topicNameDelete)},
		}

		fmt.Println()

		if helper.Confirmation(fmt.Sprintf("Do you really want to delete the topic %s?", topicNameDelete)) {

			resp, err := controller.DeleteTopics(topicRequest)
			helper.Check(err)

			fmt.Println()

			if resp.TopicErrorCodes[topicNameDelete] == 0 {

				fmt.Printf("%s: deleted\n", topicNameDelete)

			} else {

				fmt.Printf("%s: %s\n", topicNameDelete, resp.TopicErrorCodes[strings.TrimSpace(topicNameDelete)])

			}

		} else {

			fmt.Println("\nnothing done")

		}

		fmt.Println()

	},
}

func init() {
	topicCmd.AddCommand(topicDeleteCmd)

	topicDeleteCmd.Flags().StringVar(&topicNameDelete, "name", "", "topic name")
	topicDeleteCmd.MarkFlagRequired("name")
}
