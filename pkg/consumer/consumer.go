package consumer

import (
	"github.com/golang-base-template/pkg/consumer/ping"
	"log"
	"runtime/debug"
	"time"

	"github.com/bitly/go-nsq"

	"github.com/golang-base-template/util/config"
)

//Consumer Identifier
const (
	pingIdentifier = "ping-consumer"
)

func Init(cfg *config.Config) {
	//Init Consumer configurations
	consumerCfg := nsq.NewConfig()

	consumerCfg.MaxAttempts = cfg.Consumer.DefaultMaxAttempts
	consumerCfg.MaxBackoffDuration = time.Duration(cfg.Consumer.MaxBackoffDuration) * time.Second
	consumerCfg.DefaultRequeueDelay = time.Duration(cfg.Consumer.DefaultRequeueDelay) * time.Second

	for key, value := range cfg.ConsumerList {
		if value.Config == nil {
			if value.MaxInFlight != 0 {
				consumerCfg.MaxInFlight = value.MaxInFlight
			} else {
				consumerCfg.MaxInFlight = cfg.Consumer.DefaultMaxInflight
			}
			value.Config = consumerCfg
		}

		if !value.Switch {
			continue
		}

		for i := 0; i < value.WorkerAmount; i++ {
			consumer, err := nsq.NewConsumer(value.Topic, value.Channel, value.Config)
			if err != nil {
				//TODO: log error
				log.Printf("[Consumer] error when creating consumer for %s, err : %s\n", key, err.Error())
			}

			switch key {
			case pingIdentifier:
				consumer.AddHandler(GuardConsumer(ping.PingHandler{
					Topic:   value.Topic,
					Channel: value.Channel,
				}, value.Topic, value.Channel))
			default:
				//TODO: log error
				log.Printf("[Consumer] fail to consume from topic: %s", value.Topic)
			}

			err = consumer.ConnectToNSQLookupds(cfg.Consumer.LookupdAddress)
			if err != nil {
				//TODO: log error
				log.Printf("[Consumer] err connecting to NSQLookupd %s #%d, err : %s", key, i+1, err.Error())
			} else {
				//TODO: log access
				log.Printf("[Consumer] %s #%d is listening", key, i+1)
			}
		}
	}
}

// GuardConsumer protect consumer when it gets panic and then eventually recover
func GuardConsumer(fn nsq.Handler, topic, channel string) nsq.Handler {
	return nsq.HandlerFunc(func(message *nsq.Message) error {
		defer nsqPanicHandler(topic, channel, message.Body)
		return fn.HandleMessage(message)
	})
}

func nsqPanicHandler(topic, channel string, message []byte) {
	if r := recover(); r != nil {
		//TODO: log error
		log.Printf("[PANIC] %v topic : %s, channel : %s", r, topic, channel)
		log.Printf("%s", string(debug.Stack()))
	}
}
