package main

import (
	"context"
	"log"

	"github.com/infantbee/skills/golang/kafka"
)

var (
	USER_BROKERS []string = []string{"127.0.0.1:9094"}
	TOPIC                 = "test1"
	TOPIC_GROUP           = "group1"
)

func main() {
	ctx := context.Background()

	// producer
	// msg := make(chan []byte)
	// go func() {
	// 	p := kafka.NewKafkaAsyncProducer(USER_BROKERS, TOPIC)
	// 	p.LoopAsyncProducer(ctx, msg)
	// }()

	// go func() {
	// 	tk := time.NewTicker(1 * time.Second)
	// 	for {
	// 		select {
	// 		case <-tk.C:
	// 			msg <- []byte("xxxxx")
	// 			//
	// 			log.Printf("11111")
	// 		}
	// 	}
	// }()

	//consumer
	c := kafka.NewKafkaConsumer(USER_BROKERS, TOPIC_GROUP, TOPIC, "")
	c.LoopConsume(ctx, func(tag string, data []byte) error {
		log.Printf("kafka consumer get tag:[%s], data:[%b].", tag, data)
		//todo
		return nil
	})

	select {}

	// 看
}
