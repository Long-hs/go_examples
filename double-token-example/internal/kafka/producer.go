package kafka

import (
	"double-token-example/internal/config"
	"github.com/IBM/sarama"
	"log"
)

type Producer struct {
	syncProducer sarama.SyncProducer
}

var producer *Producer

func init() {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Retry.Max = 5
	cfg.Net.MaxOpenRequests = 1
	cfg.Producer.Idempotent = true
	// 创建生产者
	p, err := sarama.NewSyncProducer(config.Cfg.Kafka.Brokers, cfg)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	producer = &Producer{
		syncProducer: p,
	}
}

func GetProducer() *Producer {
	return producer
}

func (p *Producer) Send(msg *sarama.ProducerMessage) error {
	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}
	log.Printf("Message sent to partition %d at offset %d", partition, offset)
	return nil
}
