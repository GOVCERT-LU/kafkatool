// Copyright (C) 2018, CERT Gouvernemental (GOVCERT.LU)
// Author: Daniel Struck <daniel.struck@govcert.etat.lu>

package cmd

import (
	"fmt"
	"kafkatool/helper"
	"sort"

	"github.com/Shopify/sarama"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	groupsListLagThreshold int64
)

// consumerGroupListCmd represents the list command
var groupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all consumer groups",
	Run: func(cmd *cobra.Command, args []string) {

		kafkaClient := helper.ConnectKafkaClient()
		defer kafkaClient.Close()

		clusterAdmin, err := sarama.NewClusterAdminFromClient(kafkaClient)
		helper.Check(err)
		defer clusterAdmin.Close()

		groups, err := clusterAdmin.ListConsumerGroups()
		helper.Check(err)

		sortedGroupKeys := make([]string, 0, len(groups))
		for k := range groups {
			sortedGroupKeys = append(sortedGroupKeys, k)
		}
		sort.Strings(sortedGroupKeys)

		for _, consumerGroup := range sortedGroupKeys {

			details, err := clusterAdmin.DescribeConsumerGroups([]string{consumerGroup})
			helper.Check(err)

			colorMain := color.New(color.FgHiYellow, color.BgBlue, color.Bold).SprintFunc()
			colorOK := color.New(color.FgGreen, color.Bold).SprintFunc()
			colorNOK := color.New(color.FgHiRed, color.Bold).SprintFunc()
			colorBold := color.New(color.Bold).SprintFunc()

			var state string
			if details[0].State == "Stable" {
				state = colorOK("Stable")
			} else {
				state = colorNOK(details[0].State)
			}

			fmt.Printf("consumer_group=%s state=%s\n", colorMain(consumerGroup), state)

			for _, groupMemberDescription := range details[0].Members {
				fmt.Printf("client_id=%s, host=%s", colorBold(groupMemberDescription.ClientId), colorBold(groupMemberDescription.ClientHost[1:]))

				assignment, err := groupMemberDescription.GetMemberAssignment()
				helper.Check(err)

				for topic, partitions := range assignment.Topics {

					fmt.Printf(", topic=%s (", colorBold(topic))

					for i, p := range partitions {

						offsetFetchResponse, err := clusterAdmin.ListConsumerGroupOffsets(consumerGroup, assignment.Topics)
						helper.Check(err)
						partitionOffset, err := kafkaClient.GetOffset(topic, p, sarama.OffsetNewest)
						helper.Check(err)

						if i > 0 {
							fmt.Print(", ")
						}

						lag := partitionOffset - offsetFetchResponse.Blocks[topic][p].Offset
						if lag > groupsListLagThreshold {
							fmt.Printf("partition=%s lag=%s", colorBold(p), colorNOK(lag))
						} else {
							fmt.Printf("partition=%s lag=%s", colorBold(p), colorOK(lag))
						}

					}
					fmt.Print(")")

				}
				fmt.Println()
			}

			fmt.Println()

		}

	},
}

func init() {
	groupsCmd.AddCommand(groupsListCmd)

	groupsListCmd.Flags().Int64Var(&groupsListLagThreshold, "lag-threshold", 10000, "threshold to flag the lag red")
}
