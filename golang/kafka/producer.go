package kafka

import (
	"context"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaAsyncProducer struct {
	Brokers []string
	Topic   string
	Cp      sarama.AsyncProducer
}

//
func NewKafkaAsyncProducer(brokers []string, topic string) *KafkaAsyncProducer {
	return &KafkaAsyncProducer{
		Brokers: brokers,
		Topic:   topic,
	}
}

//
func (kp *KafkaAsyncProducer) LoopAsyncProducer(ctx context.Context, text <-chan []byte) error {
	if kp.Cp == nil {
		if err := kp.newAsyncProducer(); err != nil {
			log.Printf("newAsyncProducer error:[%s]", err.Error())
			return err
		}
	}
	defer kp.Cp.AsyncClose()

	// produce msg
	for {
		select {
		case <-ctx.Done():
			//zlog.Info(fmt.Sprintf("kafka async producer stoped"), "reason", ctx.Err().Error())
			return ctx.Err()

		case suc := <-kp.Cp.Successes():
			log.Printf("kafka async producer success, offset: %d, timestamp: %s, partitions: %d", suc.Offset, suc.Timestamp.String(), suc.Partition)

		case fail := <-kp.Cp.Errors():
			log.Printf("kafka async producer error:[%s]", fail.Err.Error())

		case dt := <-text:
			msg := &sarama.ProducerMessage{
				Topic: kp.Topic,
				Value: sarama.ByteEncoder(dt),
				//Key:   //暂时不设分区key,按轮询写入分区
			}
			// send msg by channel
			kp.Cp.Input() <- msg
			//zlog.Debug("kafka async producer send msgs", "detail", string(dt))
		} //end select
	}

	return nil
}

func (kp *KafkaAsyncProducer) newAsyncProducer() error {
	config := sarama.NewConfig()
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = 5 * time.Second
	config.Producer.RequiredAcks = sarama.WaitForAll          //等待服务器所有副本都保存成功后的响应
	config.Producer.Partitioner = sarama.NewRandomPartitioner //随机向partition发送消息
	config.Version = sarama.V0_10_0_1                         //设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置

	var err error
	kp.Cp, err = sarama.NewAsyncProducer(kp.Brokers, config)
	return err
}
