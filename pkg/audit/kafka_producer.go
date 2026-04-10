package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	ch     chan Event
	done   chan struct{}
}

const kafkaProducerBufSize = 1024

func NewKafkaProducer(brokers []string, topic, clientID string) *KafkaProducer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		BatchTimeout:           10 * time.Millisecond,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
	}
	if clientID != "" {
		w.Transport = &kafka.Transport{ClientID: clientID}
	}
	p := &KafkaProducer{writer: w, ch: make(chan Event, kafkaProducerBufSize), done: make(chan struct{})}
	go p.loop()
	return p
}

func (p *KafkaProducer) Publish(_ context.Context, event Event) {
	select {
	case p.ch <- event:
	default:
		log.Printf("[audit] producer buffer full, dropping event for tool %q", event.ToolName)
	}
}

func (p *KafkaProducer) Close() error {
	close(p.ch)
	<-p.done
	return p.writer.Close()
}

func (p *KafkaProducer) loop() {
	defer close(p.done)
	for event := range p.ch {
		p.send(event)
	}
}

func (p *KafkaProducer) send(event Event) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("[audit] failed to marshal event for tool %q: %v", event.ToolName, err)
		return
	}
	msg := kafka.Message{Key: []byte(event.ID), Value: payload, Time: event.Timestamp}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("[audit] failed to write event for tool %q: %v", event.ToolName, fmt.Errorf("kafka write: %w", err))
	}
}
