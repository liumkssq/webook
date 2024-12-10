package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

var addrs = []string{"127.0.0.1:9094"}

func TestSyncProduce(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("hello world!!!"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("key"),
				Value: []byte("value"),
			},
		},
		Metadata: "metadata",
	})
	assert.NoError(t, err)
	producer.Close()
}

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addrs, "test_group", cfg)
	require.NoError(t, err)
	err = consumer.Consume(context.Background(), []string{"test_topic"}, testConsumerGroupHandler{})
	if err != nil {
		panic(err)
	}
}

type testConsumerGroupHandler struct {
}

func (t testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("setup consumer group")
	return nil
}

func (t testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (t testConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		log.Println("receive msg", string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

type MyBizMsg struct {
	Name string
}
