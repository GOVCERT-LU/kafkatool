// Copyright (C) 2025, CERT Gouvernemental (GOVCERT.LU)
// Author: Daniel Struck <daniel.struck@govcert.etat.lu>

package cmd

import (
	"fmt"
	"kafkatool/helper"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"

	"github.com/IBM/sarama"
	"github.com/spf13/cobra"
)

var (
	topicNameRead string
	consumed      uint64
	wg            sync.WaitGroup
	offsetOldest  bool // Flag for oldest offset
	offsetNewest  bool // Flag for newest offset
	verbose       bool
)

func readFromPartition(partitionConsumer sarama.PartitionConsumer, partitionID int32, signals chan os.Signal) {

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			if verbose {
				fmt.Printf("partition %d offset %d: %s\n", partitionID, msg.Offset, msg.Value)
			} else {
				fmt.Printf("%s\n", msg.Value)
			}
			atomic.AddUint64(&consumed, 1)
		case <-signals:
			break ConsumerLoop
		}
	}

	wg.Done()

}

// createCmd represents the create command
var topicReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read messages from a topic",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		consumer := helper.GetConsumer()
		defer func() {
			if err := consumer.Close(); err != nil {
				log.Fatalln(err)
			}
		}()

		var offset int64 = sarama.OffsetNewest // Default offset
		if offsetOldest {
			offset = sarama.OffsetOldest
		}

		// Sarama's Consumer type does not currently support automatic consumer-group rebalancing and offset tracking.
		// workaround at the moment
		for i := int32(0); i < 64; i++ {

			partitionConsumer, err := consumer.ConsumePartition(topicNameRead, i, offset)
			if err != nil {
				break
			}

			// Trap SIGINT to trigger a shutdown.
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt)

			wg.Add(1)
			log.Printf("created consumer for partition %d", i)
			go readFromPartition(partitionConsumer, i, signals)

		}

		log.Println("Stop the consumers with CTRL+c")

		wg.Wait()

		log.Printf("consumed: %d\n", atomic.LoadUint64(&consumed))

	},
}

func init() {
	topicCmd.AddCommand(topicReadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// topicReadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// topicReadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	topicReadCmd.Flags().StringVar(&topicNameRead, "name", "", "topic name")
	topicReadCmd.MarkFlagRequired("name")

	// Add flags for offset configuration
	topicReadCmd.Flags().BoolVar(&offsetOldest, "oldest", false, "Consume from the oldest offset")
	topicReadCmd.Flags().BoolVar(&offsetNewest, "newest", true, "Consume from the newest offset (default)") // Newest is the default
	// Add a check to prevent setting both --oldest and --newest together.  Newest is the default, so this isn't strictly necessary, but good practice.
	topicReadCmd.MarkFlagsMutuallyExclusive("oldest", "newest")

	topicReadCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}
