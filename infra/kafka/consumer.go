package kafka

import (
	"fmt"
	"log"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	MsgChan chan *ckafka.Message
}

func NewKafkaConsumer(myChan chan *ckafka.Message) *KafkaConsumer {
	return &KafkaConsumer{
		MsgChan: myChan,
	}
}

func (k *KafkaConsumer) Consume() {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KafkaBootstrapServers"),
		"group.id":          os.Getenv("KafkaConsumerGroupId"),
	}
	c, err := ckafka.NewConsumer(configMap)
	if err != nil {
		log.Fatalf("Error consuming kafka messages: " + err.Error())
	}

	topics := []string{os.Getenv("KafkaReadTopic")}
	c.SubscribeTopics(topics, nil)
	fmt.Println("Kafka consumer has been started.")

	//start consuming endlessly
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			//links the message to the channel
			//it allows retrieving the messaging elsewhere in the app
			k.MsgChan <- msg
		}
	}
}
