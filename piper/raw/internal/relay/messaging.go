package relay

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"golang.org/x/net/proxy"

	"piper/internal/tail"
)

// Producer is an abstraction over a sarama.SyncProducer that exposes a limited API surface.
type Producer struct {
	client     sarama.SyncProducer
	serializer MessageSerializer
}

// NewProducer creates a new Kafka producer for a particular target broker.
func NewProducer(
	serializer MessageSerializer,
	socks5 string,
	address string,
	acks int,
	retries int,
	timeout time.Duration,
) (*Producer, error) {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Producer.RequiredAcks = sarama.RequiredAcks(acks)
	kafkaCfg.Producer.Retry.Backoff = 5 * time.Second
	kafkaCfg.Producer.Return.Successes = true
	kafkaCfg.Producer.Return.Errors = true
	if retries != 0 {
		kafkaCfg.Producer.Retry.Max = retries
	}
	if timeout != 0 {
		kafkaCfg.Producer.Timeout = timeout
	}
	if socks5 != "" {
		dialer, err := proxy.SOCKS5("tcp", socks5, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}

		kafkaCfg.Net.Proxy.Enable = true
		kafkaCfg.Net.Proxy.Dialer = dialer
	}

	client, err := sarama.NewSyncProducer([]string{address}, kafkaCfg)
	if err != nil {
		return nil, fmt.Errorf("relay: failed to create client: err=%v", err)
	}

	return &Producer{
		client:     client,
		serializer: serializer,
	}, nil
}

// WriteMessage synchronously writes a tail.Message to a particular Kafka topic, blocking until the
// write has been sent and has reached the Kafka acknowledgement level specified in configuration.
// It returns the length of the produced content in bytes and an error as applicable.
func (p *Producer) WriteMessage(topic string, message tail.Message) (int, error) {
	serialized, err := p.serializer.Serialize(message)
	if err != nil {
		return 0, err
	}

	value := sarama.ByteEncoder(serialized)

	if _, _, err := p.client.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Value:     value,
		Timestamp: message.Timestamp,
	}); err != nil {
		return 0, fmt.Errorf("relay: failed to send message: err=%v", err)
	}

	return value.Length(), nil
}

// Close closes the underlying Kafka client.
func (p *Producer) Close() error {
	return p.client.Close()
}
