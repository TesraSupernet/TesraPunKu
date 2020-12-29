package mq

import (
	"github.com/apache/pulsar-client-go/pulsar"
)

// consume message
type Message struct {
	C <-chan pulsar.Message
}
