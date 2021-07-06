package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type KafkaConsumer struct {
	Brokers    []string
	Topic      string
	Group      string
	Offset     string
	Partitions []int32
	Cs         *cluster.Consumer
}

//
func NewKafkaConsumer(brokers []string, group, topic, offset string) *KafkaConsumer {
	return &KafkaConsumer{
		Brokers: brokers,
		Group:   group,
		Topic:   topic,
		Offset:  offset,
	}
}

//
func (kc *KafkaConsumer) LoopConsume(ctx context.Context, dealFunc func(string, []byte) error) error {
	if kc.Cs == nil {
		if err := kc.newConsumer(); err != nil {
			log.Printf("topic:[%s] kafka consumer connect error:[%s].", kc.Topic, err.Error())
			return err
		}
	}
	defer kc.Cs.Close()

	// Consume all channels, wait for context to exit
	for {
		select {
		case <-ctx.Done():
			log.Printf("topic:[%s] partition:[%#v] kafka consumer exiting... by %#v", kc.Topic, kc.Partitions, ctx.Err().Error())
			return ctx.Err()

		case ccp, ok := <-kc.Cs.Partitions():
			if ok {
				kc.Partitions = append(kc.Partitions, ccp.Partition())
				log.Printf("topic:[%s] partition:[%#v] kafka consumer working...", ccp.Topic(), ccp.Partition())
			}

		case ntf, more := <-kc.Cs.Notifications():
			if more {
				kc.Partitions = ntf.Current[kc.Topic]
				log.Printf("topic:[%s] kafka consumer rebalancing... NotificationType:[%s], partition list is:[%#v]", kc.Topic, ntf.Type.String(), kc.Partitions)
			}

		case msg, more := <-kc.Cs.Messages():
			if more {
				// sdk_{businessid}_{topic}_{partition}_{offset} ---> md5
				idsrc := fmt.Sprintf("sdk_%s_%s_%d_%d", "business_id", msg.Topic, msg.Partition, msg.Offset)
				dealFunc(idsrc, msg.Value) // idempotent
				kc.Cs.MarkOffset(msg, "")

				//zlog.Info(fmt.Sprintf("topic:[%s] kafka consumer message...", kc.Topic), "detail", fmt.Sprintf("topic:%s, partition:%d, offset:%d, value:%s", msg.Topic, msg.Partition, msg.Offset, string(msg.Value)))
			}

		case err, more := <-kc.Cs.Errors():
			if more {
				log.Printf("topic:[%s] kafka consumer error:[%s]...", kc.Topic, err.Error())
			}
		} //end select
	}
	return nil
}

func (kc *KafkaConsumer) newConsumer() error {
	config := cluster.NewConfig()
	//config.Consumer.MaxWaitTime = 5 * time.Second
	config.Consumer.Offsets.CommitInterval = 1 * time.Second //自动提交offset频率
	config.Group.Return.Notifications = true
	config.Group.PartitionStrategy = cluster.StrategyRange     //分区分配策略
	config.Consumer.Group.Session.Timeout = 6 * time.Second    //(default 10s)
	config.Consumer.Group.Heartbeat.Interval = 2 * time.Second //(default 3s)
	//config.Consumer.Group.Rebalance.Timeout = 120 * time.Second //(default 60s).
	config.Version = sarama.V0_10_2_0
	// offset setting
	if kc.Offset == "earliest" {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	var err error
	kc.Cs, err = cluster.NewConsumer(kc.Brokers, kc.Group, []string{kc.Topic}, config)
	return err
}
