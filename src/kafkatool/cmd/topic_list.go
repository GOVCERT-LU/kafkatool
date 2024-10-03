// Copyright (C) 2018, CERT Gouvernemental (GOVCERT.LU)
// Author: Daniel Struck <daniel.struck@govcert.etat.lu>

package cmd

import (
	"fmt"
	"kafkatool/helper"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/IBM/sarama"
	termbox "github.com/nsf/termbox-go"
)

type partitionData struct {
	leader        int32
	colorReplicas termbox.Attribute
	colorIsr      termbox.Attribute
	id            int32
	replicas      []int32
	isr           []int32
}

type topicData struct {
	name        string
	internal    bool
	colorLeader termbox.Attribute
	partitions  []*partitionData
}

var (
	topicListInternal bool
	topicListDetails  bool
	topicListMonitor  bool
)

func printText(x int, y int, text string, fg termbox.Attribute) (int, int) {

	for _, r := range text {
		termbox.SetCell(x, y, r, fg, termbox.ColorDefault)
		x++
	}

	return x, y
}

func drawHead(metadata *sarama.MetadataResponse, page int, pageCount int) (int, int) {

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	x, y := 0, 0

	printText(x, y, fmt.Sprintf("press [q] or [ESC] to quit, page %d from %d page(s) (%v)", page, pageCount, time.Now()), termbox.ColorDefault)
	y += 2 // advance two lines
	x = 0  // reset to beginning

	printText(x, y, fmt.Sprintf("controller id: %d", metadata.ControllerID), termbox.AttrBold)
	y += 2
	x = 0

	return x, y

}

func retrieveTopicDataset(metadata *sarama.MetadataResponse) []*topicData {

	topicDataset := make([]*topicData, 0)

	for _, topic := range metadata.Topics {

		// skip internal topics by default
		if topic.IsInternal && !topicListInternal {
			continue
		}

		topicData := new(topicData)

		topicData.name = topic.Name
		topicData.internal = topic.IsInternal
		topicData.partitions = make([]*partitionData, 0)

		// how many different leaders does this topic have?
		leaderCount := make(map[int32]bool)
		for _, partition := range topic.Partitions {
			leaderCount[partition.Leader] = true
		}

		switch {
		case len(leaderCount) > 2:
			topicData.colorLeader = termbox.ColorGreen
		case len(leaderCount) == 2:
			topicData.colorLeader = termbox.ColorYellow
		default:
			topicData.colorLeader = termbox.ColorRed
		}

		// retrieve the topic partitions
		for _, partition := range topic.Partitions {

			partitionData := new(partitionData)

			switch {
			case len(partition.Replicas) > 2:
				partitionData.colorReplicas = termbox.ColorGreen
			case len(partition.Replicas) == 2:
				partitionData.colorReplicas = termbox.ColorYellow
			default:
				partitionData.colorReplicas = termbox.ColorRed
			}

			switch {
			case len(partition.Isr) < 2:
				partitionData.colorIsr = termbox.ColorRed
			case len(partition.Replicas) != len(partition.Isr):
				partitionData.colorIsr = termbox.ColorYellow
			default:
				partitionData.colorIsr = termbox.ColorGreen
			}

			partitionData.leader = partition.Leader
			partitionData.id = partition.ID
			partitionData.replicas = partition.Replicas
			partitionData.isr = partition.Isr

			sort.Slice(partitionData.replicas, func(i, j int) bool {
				return partitionData.replicas[i] < partitionData.replicas[j]
			})
			sort.Slice(partitionData.isr, func(i, j int) bool {
				return partitionData.isr[i] < partitionData.isr[j]
			})

			topicData.partitions = append(topicData.partitions, partitionData)
		}
		sort.Slice(topicData.partitions, func(i, j int) bool {
			return topicData.partitions[i].id < topicData.partitions[j].id
		})

		topicDataset = append(topicDataset, topicData)

	}

	sort.Slice(topicDataset, func(i, j int) bool {
		return topicDataset[i].name < topicDataset[j].name
	})

	return topicDataset

}

func draw() {

	broker, err := helper.ConnectKafkaClient().Controller()
	helper.Check(err)
	defer broker.Close()

	metadata := helper.RetrieveMetadata(broker)
	topicDataset := retrieveTopicDataset(metadata)

	// determine the number of pages
	_, termHeight := termbox.Size()
	headerSize := 4
	pageCount := 1
	y := headerSize
	for _, topic := range topicDataset {

		if y+len(topic.partitions)+2 > termHeight {
			pageCount++
			y = headerSize
		}

		y += 2 + len(topic.partitions)

	}

	// show content
	page := 1
	x, y := drawHead(metadata, page, pageCount)
	for _, topic := range topicDataset {

		if y+len(topic.partitions)+2 > termHeight {
			page++
			termbox.Flush()
			time.Sleep(3 * time.Second)
			x, y = drawHead(metadata, page, pageCount)
		}

		x = 0
		x, y = printText(x, y, "topic name: ", termbox.ColorDefault)
		x, y = printText(x, y, topic.name, termbox.AttrBold)
		x, y = printText(x, y, fmt.Sprintf(", internal: %t", topic.internal), termbox.ColorDefault)
		y++
		x = 0
		for _, partition := range topic.partitions {

			x = 0
			x, y = printText(x, y, fmt.Sprintf("partition id: %d, leader: ", partition.id), termbox.ColorDefault)
			x, y = printText(x, y, fmt.Sprintf("%d", partition.leader), topic.colorLeader)
			x, y = printText(x, y, ", replicas: ", termbox.ColorDefault)
			x, y = printText(x, y, fmt.Sprintf("%v", partition.replicas), partition.colorReplicas)
			x, y = printText(x, y, ", isr: ", termbox.ColorDefault)
			x, y = printText(x, y, fmt.Sprintf("%v", partition.isr), partition.colorIsr)
			y++
		}
		y++
	}

	termbox.Flush()
	time.Sleep(3 * time.Second)

}

func monitor() {
	// init termbox
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Ch == 113 || ev.Key == termbox.KeyCtrlC) {
				termbox.Close()
				os.Exit(0)
			}
		}
	}()

	// draw the content in a loop
	for {
		draw()
	}
}

// topicListCmd represents the list command
var topicListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all partitions",
	Run: func(cmd *cobra.Command, args []string) {

		if topicListMonitor {

			monitor()

		} else {

			broker, err := helper.ConnectKafkaClient().Controller()
			helper.Check(err)
			defer broker.Close()

			metadata := helper.RetrieveMetadata(broker)
			topicDataset := retrieveTopicDataset(metadata)

			for _, topic := range topicDataset {

				fmt.Println(topic.name)

				if topicListDetails {
					for _, partition := range topic.partitions {
						fmt.Printf("id=%d leader=%d replicas=%v isr=%v\n", partition.id, partition.leader, partition.replicas, partition.isr)
					}
					fmt.Println()
				}
			}
		}

	},
}

func init() {
	topicCmd.AddCommand(topicListCmd)

	topicListCmd.Flags().BoolVar(&topicListInternal, "internal", false, "Show internal topics")
	topicListCmd.Flags().BoolVar(&topicListDetails, "details", false, "Show topic details")
	topicListCmd.Flags().BoolVar(&topicListMonitor, "monitor", false, "Continuously monitor the topics")

}
