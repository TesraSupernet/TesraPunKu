package mq

import (
	"context"
	"fmt"
	"github.com/DOSNetwork/core/edge_cm/tools"
	"github.com/apache/pulsar-client-go/pulsar"
)

// listen topic and return message
// get message   e.g: <- message
func Receive(topic, subName, mod string) *Message {
	tools.Logger.Info(fmt.Sprintf("%s,%s,%s", topic, subName, mod))
	var subType pulsar.SubscriptionType
	switch mod {
	case "shared":
		subType = pulsar.Shared
	case "failover":
		subType = pulsar.Failover
	case "keyShared":
		subType = pulsar.KeyShared
	default:
		subType = pulsar.Exclusive
	}
	// Use the client object to instantiate a consumer
	consumer, err := client().Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: subName,
		Type:             subType,
	})
	if err != nil {
		tools.Logger.Panic(err.Error())
		tools.Logger.Panic("create consumer err.")
	}
	ctx := context.Background()
	c := make(chan pulsar.Message, 1)
	//Listen indefinitely on the topic
	go func(ch chan pulsar.Message) {
		defer consumer.Close()
		for {
			msg, err := consumer.Receive(ctx)
			if err != nil {
				// Failed to process messages
				consumer.Nack(msg)
			} else {
				// process message
				ch <- msg
				// Message processed successfully
				consumer.Ack(msg)
			}
		}
	}(c)
	return &Message{C: c}
}
