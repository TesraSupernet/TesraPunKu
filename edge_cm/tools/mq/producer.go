package mq

import (
	"context"
	"fmt"
	"github.com/DOSNetwork/core/edge_cm/conf"
	"github.com/DOSNetwork/core/edge_cm/tools"
	"github.com/apache/pulsar-client-go/pulsar"
	log "github.com/sirupsen/logrus"
	"time"
)

var clientInstance pulsar.Client

func client() pulsar.Client {
	if clientInstance == nil {
		tools.Logger.Info(fmt.Sprintf("url:%s,timeOut:%d,token:%s", conf.Conf.Pulsar.PulsarURL, conf.Conf.Pulsar.ConnectionTimeout, conf.Conf.Pulsar.Token))
		cli, err := pulsar.NewClient(pulsar.ClientOptions{
			URL:               conf.Conf.Pulsar.PulsarURL,
			ConnectionTimeout: time.Second * time.Duration(conf.Conf.Pulsar.ConnectionTimeout),
			OperationTimeout:  time.Second * time.Duration(conf.Conf.Pulsar.ConnectionTimeout),
			Authentication:    pulsar.NewAuthenticationToken(conf.Conf.Pulsar.Token),
		})
		if err != nil {
			log.Panicf("create pulsar client error: %s", err.Error())
		}
		clientInstance = cli
	}
	return clientInstance
}

// init pulsar logger
func initLogger(l string) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05.000",
	})
	var level log.Level
	switch l {
	case "info":
		level = log.InfoLevel
	case "warn":
		level = log.WarnLevel
	case "error":
		level = log.ErrorLevel
	default:
		level = log.DebugLevel
	}
	log.SetLevel(level)
}

func init() {
	initLogger("error")
}

func Send(sync bool, topic string, payload []byte) (status bool, retErr error) {
	producer, err := client().CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return false, err
	}
	defer producer.Close()

	ctx := context.Background()
	msg := &pulsar.ProducerMessage{
		Payload: payload,
	}
	if sync {
		// Attempt to send the message
		if _, err := producer.Send(ctx, msg); err != nil {
			status = false
			retErr = err
		} else {
			status = true
		}
	} else {
		// Attempt to send the message asynchronously and handle the response
		producer.SendAsync(ctx, msg, func(MessageID pulsar.MessageID, pMsg *pulsar.ProducerMessage, err error) {
			if err != nil {
				status = false
				retErr = err
			}
		})
		status = true
	}
	return
}
