// Copyright (C) 2018, CERT Gouvernemental (GOVCERT.LU)
// Author: Daniel Struck <daniel.struck@govcert.etat.lu>

package cmd

import (
	"bufio"
	"kafkatool/helper"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

var topicNameWrite string
var fileNameWrite string

// createCmd represents the create command
var topicWriteCmd = &cobra.Command{
	Use:   "write",
	Short: "Write messages to a topic, by default an increasing number is written to the topic",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		producer := helper.GetProducer()

		defer func() {
			if err := producer.Close(); err != nil {
				log.Fatalln(err)
			}
		}()

		// Trap SIGINT to trigger a graceful shutdown.
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		var (
			wg                          sync.WaitGroup
			enqueued, successes, errors int
		)

		dataChannel := make(chan string)

		if fileNameWrite != "" {

			file, err := os.Open(fileNameWrite)
			helper.Check(err)

			scanner := bufio.NewScanner(file)

			go func(dataChan chan<- string, scanner *bufio.Scanner) {

				defer close(dataChan)

				for scanner.Scan() {
					dataChan <- scanner.Text()
				}

				if err := scanner.Err(); err != nil {
					log.Fatalf("error reading file %s", fileNameWrite)
				}

			}(dataChannel, scanner)

			log.Printf("file %s is written to topic %s", fileNameWrite, topicNameWrite)

		} else {

			// by default write a continous number to the topic
			go func(dataChan chan<- string) {

				count := 1
				for {
					dataChan <- strconv.Itoa(count)
					count++
				}

			}(dataChannel)

		}

		// count successes
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range producer.Successes() {
				successes++
			}
		}()

		// count errors
		wg.Add(1)
		go func() {
			defer wg.Done()
			for err := range producer.Errors() {
				log.Println(err)
				errors++
			}
		}()

		// write messages to the topic
	ProducerLoop:
		for msg := range dataChannel {

			message := &sarama.ProducerMessage{Topic: topicNameWrite, Value: sarama.StringEncoder(msg)}
			select {
			case producer.Input() <- message:
				enqueued++

			case <-signals:
				producer.AsyncClose() // Trigger a shutdown of the producer.
				break ProducerLoop
			}

		}

		producer.AsyncClose() // Trigger a shutdown of the producer.

		wg.Wait()

		log.Printf("finished: enqueued %d; produced: %d; errors: %d\n", enqueued, successes, errors)

	},
}

func init() {
	topicCmd.AddCommand(topicWriteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// topicWriteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// topicWriteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	topicWriteCmd.Flags().StringVar(&topicNameWrite, "name", "", "topic name")
	topicWriteCmd.Flags().StringVar(&fileNameWrite, "file", "", "file name (written to the topic line by line)")
	topicWriteCmd.MarkFlagRequired("name")
}
