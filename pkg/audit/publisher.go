package audit

import "context"

type Publisher interface {
	Publish(ctx context.Context, event Event)
	Close() error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(_ context.Context, _ Event) {}
func (NoopPublisher) Close() error                       { return nil }
