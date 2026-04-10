package auditstore

import (
	"context"

	"github.com/chiragms-cdcx/coindcx-cs-mcp/pkg/audit"
)

type Store interface {
	Insert(ctx context.Context, event audit.Event) error
	Close() error
}
