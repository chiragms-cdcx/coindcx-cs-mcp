CREATE SCHEMA IF NOT EXISTS audit;

CREATE TABLE IF NOT EXISTS audit.mcp_requests (
    id              TEXT        PRIMARY KEY,
    timestamp       TIMESTAMPTZ NOT NULL,
    tool_name       TEXT        NOT NULL,
    request_id      TEXT        NOT NULL,
    status          TEXT        NOT NULL,
    error_message   TEXT,
    transport       TEXT        NOT NULL,
    latency_ms      BIGINT      NOT NULL,
    params          JSONB,
    agent_id        TEXT,
    agent_email     TEXT,
    target_user_id  TEXT,
    session_id      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_mcp_requests_timestamp
    ON audit.mcp_requests (timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_audit_mcp_requests_tool_name
    ON audit.mcp_requests (tool_name);

CREATE INDEX IF NOT EXISTS idx_audit_mcp_requests_agent_id
    ON audit.mcp_requests (agent_id);

CREATE INDEX IF NOT EXISTS idx_audit_mcp_requests_target_user_id
    ON audit.mcp_requests (target_user_id);

CREATE INDEX IF NOT EXISTS idx_audit_mcp_requests_session_id
    ON audit.mcp_requests (session_id);
