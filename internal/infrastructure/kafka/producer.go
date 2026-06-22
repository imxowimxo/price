package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaNotifier struct {
	writer *kafka.Writer
}

func NewKafkaNotifier(writer *kafka.Writer) *KafkaNotifier {
	return &KafkaNotifier{
		writer: writer,
	}
}

func (k *KafkaNotifier) SendMessage(ctx context.Context, topic string, key string, payload []byte) error {

	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: payload,
		Topic: topic,
	})
}
