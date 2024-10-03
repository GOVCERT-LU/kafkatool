package helper

import (
	"fmt"
	"strings"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

var (
	metadataRequestVersion = int16(1)
)

func kafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_1_0_0
	config.Net.TLS.Enable = viper.GetBool("tls")
	config.Producer.Return.Successes = true

	if viper.GetBool("sasl_plain") {
		config.Net.SASL.Handshake = true
		config.Net.SASL.Enable = true
		config.Net.SASL.User = viper.GetString("sasl_username")
		config.Net.SASL.Password = viper.GetString("sasl_password")
	}

	return config

}

func getFirstBrokerName() string {
	return fmt.Sprintf("%s:%d", strings.Split(viper.GetString("broker"), ",")[0], viper.GetInt("port"))
}

func getBrokerNames() []string {

	brokerNames := make([]string, 0)

	for _, broker := range strings.Split(viper.GetString("broker"), ",") {
		brokerNames = append(brokerNames, fmt.Sprintf("%s:%d", broker, viper.GetInt("port")))
	}

	return brokerNames

}

// ConnectKafkaClient makes a client connection to a Kafka cluster.
func ConnectKafkaClient() sarama.Client {

	client, err := sarama.NewClient(getBrokerNames(), kafkaConfig())
	Check(err)

	return client

}

// RetrieveMetadata retrieve the metadata from the Kafka broker.
func RetrieveMetadata(broker *sarama.Broker) *sarama.MetadataResponse {

	request := new(sarama.MetadataRequest)
	request.Version = metadataRequestVersion

	response, err := broker.GetMetadata(request)
	Check(err)

	return response

}

// GetConsumer returns a Kafka consumer.
func GetConsumer() sarama.Consumer {
	consumer, err := sarama.NewConsumer(getBrokerNames(), kafkaConfig())
	Check(err)
	return consumer
}

// GetProducer returns a Kafka consumer.
func GetProducer() sarama.AsyncProducer {
	producer, err := sarama.NewAsyncProducer(getBrokerNames(), kafkaConfig())
	Check(err)
	return producer
}
