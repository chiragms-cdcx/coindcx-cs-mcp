package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	kafka "github.com/segmentio/kafka-go"

	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/config"
	"github.com/chiragms-cdcx/coindcx-cs-mcp/pkg/audit"
	"github.com/chiragms-cdcx/coindcx-cs-mcp/pkg/auditstore"
)

func main() {
	cfg := config.LoadAuditLogs()
	if !cfg.Enabled {
		log.Fatal("audit logs not enabled")
	}
	if cfg.KafkaBrokers == "" {
		log.Fatal("AUDIT_KAFKA_BROKERS required")
	}
	if cfg.DBDsn == "" {
		log.Fatal("AUDIT_DB_DSN required")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	store, err := auditstore.NewPostgresStore(ctx, cfg.DBDsn)
	if err != nil {
		log.Fatalf("audit-consumer: connect to postgres: %v", err)
	}
	defer store.Close()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(cfg.KafkaBrokers, ","), Topic: cfg.KafkaTopic, GroupID: cfg.ConsumerGroup,
		MinBytes: 10e3, MaxBytes: 10e6, ReadBackoffMin: 500 * time.Millisecond, ReadBackoffMax: 5 * time.Second,
	})
	defer reader.Close()
	log.Printf("audit-consumer: consuming topic=%s group=%s", cfg.KafkaTopic, cfg.ConsumerGroup)
	const retryDelayMin, retryDelayMax = 2 * time.Second, 30 * time.Second
	retryDelay := retryDelayMin
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("audit-consumer: shutting down")
				return
			}
			log.Printf("audit-consumer: read: %v — retrying in %s", err, retryDelay)
			select {
			case <-time.After(retryDelay):
			case <-ctx.Done():
				return
			}
			retryDelay *= 2
			if retryDelay > retryDelayMax {
				retryDelay = retryDelayMax
			}
			continue
		}
		retryDelay = retryDelayMin
		var event audit.Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("audit-consumer: unmarshal offset=%d: %v", msg.Offset, err)
			continue
		}
		if err := store.Insert(ctx, event); err != nil {
			log.Printf("audit-consumer: insert id=%s: %v", event.ID, err)
		}
	}
}
