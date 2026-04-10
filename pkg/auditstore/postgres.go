package auditstore

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/chiragms-cdcx/coindcx-cs-mcp/pkg/audit"
)

//go:embed migrations/001_create_audit_schema.sql
var schemaMigration string

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, dsn string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("auditstore: parse DSN: %w", err)
	}
	log.Println("auditstore: pinging database...")
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		panic(fmt.Sprintf("auditstore: database unreachable: %v", err))
	}
	log.Println("auditstore: database ping OK")
	s := &PostgresStore{pool: pool}
	if err := s.EnsureSchema(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return s, nil
}

func (s *PostgresStore) Insert(ctx context.Context, event audit.Event) error {
	const q = `
INSERT INTO audit.mcp_requests (id, timestamp, tool_name, request_id, status, error_message, transport, latency_ms, params, agent_id, agent_email, target_user_id, session_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (id) DO NOTHING`
	var params any
	if len(event.Params) > 0 {
		params = []byte(event.Params)
	}
	_, err := s.pool.Exec(ctx, q, event.ID, event.Timestamp, event.ToolName, event.RequestID, string(event.Status), event.ErrorMessage, string(event.Transport), event.LatencyMs, params, event.AgentID, event.AgentEmail, event.TargetUserID, event.SessionID)
	if err != nil {
		return fmt.Errorf("auditstore insert: %w", err)
	}
	return nil
}

func (s *PostgresStore) Close() error {
	s.pool.Close()
	return nil
}

func (s *PostgresStore) EnsureSchema(ctx context.Context) error {
	log.Println("auditstore: running schema migration...")
	if _, err := s.pool.Exec(ctx, schemaMigration); err != nil {
		return fmt.Errorf("auditstore: ensure schema: %w", err)
	}
	log.Println("auditstore: schema migration OK")
	return nil
}
