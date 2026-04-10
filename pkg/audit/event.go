package audit

import (
	"encoding/json"
	"time"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type Transport string

const (
	TransportHTTP Transport = "http"
)

type Event struct {
	ID           string          `json:"id"`
	Timestamp    time.Time       `json:"timestamp"`
	ToolName     string          `json:"tool_name"`
	RequestID    string          `json:"request_id"`
	Status       Status          `json:"status"`
	ErrorMessage string          `json:"error_message,omitempty"`
	Transport    Transport       `json:"transport"`
	LatencyMs    int64           `json:"latency_ms"`
	Params       json.RawMessage `json:"params,omitempty"`
	AgentID      string          `json:"agent_id,omitempty"`
	AgentEmail   string          `json:"agent_email,omitempty"`
	TargetUserID string          `json:"target_user_id,omitempty"`
	SessionID    string          `json:"session_id,omitempty"`
}

const maxErrorLen = 512

func TruncateError(msg string) string {
	runes := []rune(msg)
	if len(runes) <= maxErrorLen {
		return msg
	}
	return string(runes[:maxErrorLen]) + "\u2026"
}
