package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/config"
	"github.com/chiragms-cdcx/coindcx-cs-mcp/pkg/auditstore"
)

func main() {
	cfg := config.LoadAuditLogs()
	if cfg.DBDsn == "" {
		log.Fatal("audit-migrate: AUDIT_DB_DSN is required")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	log.Printf("audit-migrate: connecting to database...")
	store, err := auditstore.NewPostgresStore(ctx, cfg.DBDsn)
	if err != nil {
		log.Fatalf("audit-migrate: %v", err)
	}
	defer store.Close()
	log.Println("audit-migrate: migration complete")
}
