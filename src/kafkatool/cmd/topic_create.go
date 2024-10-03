// Copyright (C) 2018, CERT Gouvernemental (GOVCERT.LU)
// Author: Daniel Struck <daniel.struck@govcert.etat.lu>

package cmd

import (
	"fmt"
	"kafkatool/helper"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/cobra"
)

var topicNameCreate string
var numPartitions int
var replicationFactor int

// createCmd represents the create command
var topicCreateCmd = &cobra.Command{

	Use: "create",

	Short: "Create a topic",

	Run: func(cmd *cobra.Command, args []string) {

		broker, err := helper.ConnectKafkaClient().Controller()
		helper.Check(err)

		defer broker.Close()

		topicRequest := &sarama.CreateTopicsRequest{
			Version: 2,               // version 2 requires Kafka 1.0.0
			Timeout: time.Second * 2, // wait 2 seconds for the Kafka server to create the topic
			TopicDetails: map[string]*sarama.TopicDetail{
				topicNameCreate: &sarama.TopicDetail{
					NumPartitions:     int32(numPartitions),
					ReplicationFactor: int16(replicationFactor),
				},
			},
		}

		resp, err := broker.CreateTopics(topicRequest)
		helper.Check(err)

		if resp.TopicErrors[topicNameCreate].ErrMsg != nil {
			log.Fatalln(*resp.TopicErrors[topicNameCreate].ErrMsg)
		}

		fmt.Printf("\ntopic %s successfully created\n\n", topicNameCreate)

	},
}

func init() {

	topicCmd.AddCommand(topicCreateCmd)

	topicCreateCmd.Flags().StringVar(&topicNameCreate, "name", "", "topic name")
	topicCreateCmd.Flags().IntVar(&numPartitions, "num_partitions", 3, "number of partitions a topic is split into")
	topicCreateCmd.Flags().IntVar(&replicationFactor, "replication_factor", 3, "how many times each partition is replicated")

	// required parameters
	topicCreateCmd.MarkFlagRequired("name")

}
