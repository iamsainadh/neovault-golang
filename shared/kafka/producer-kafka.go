package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer interface {
	Publish(ctx context.Context, topic string, key string, message []byte) error
}

type producer struct {
	client *kgo.Client
}

func NewProducer(client *kgo.Client) Producer {
	return &producer{
		client: client,
	}
}

func (p *producer) Publish(ctx context.Context, topic string, key string, message []byte) error {
	record := &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: message,
	}

	return p.client.ProduceSync(ctx, record).FirstErr()
}
