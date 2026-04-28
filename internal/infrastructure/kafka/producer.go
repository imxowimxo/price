package kafka

import (
	"Price/internal/domain/price_drop_event"
	"context"
	"encoding/json"
	"strconv"

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

func (k *KafkaNotifier) SendPriceDrop(ctx context.Context, event price_drop_event.PriceDropEvent) error {

	byt, err := json.Marshal(event)
	if err != nil {
		return err
	}

	str := strconv.FormatInt(event.UserID, 10)

	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(str),
		Value: byt,
	})
}
