package config

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/logger"
)

const (
	DefaultBaseURL     = "https://api.coindcx.com"
	DefaultPublicURL   = "https://public.coindcx.com"
	DefaultAdminURL    = "https://admin-api.coindcx.com"
	DefaultMCPHTTPHost = "0.0.0.0"

	DefaultAuditKafkaTopic    = "coindcx-cs-mcp-audit"
	DefaultAuditConsumerGroup = "coindcx-cs-mcp-audit-consumer"
)

const (
	HeaderAdminServiceKey    = "X-ADMIN-SERVICE-KEY"
	HeaderAdminServiceSecret = "X-ADMIN-SERVICE-SECRET"
	HeaderAdminAgentID       = "X-ADMIN-AGENT-ID"
	HeaderAdminAgentEmail    = "X-ADMIN-AGENT-EMAIL"
	HeaderAdminTargetUserID  = "X-ADMIN-TARGET-USER-ID"
	HeaderAdminBaseURL       = "X-ADMIN-BASE-URL"
	HeaderAuthToken          = "X-COINDCX-AUTH-TOKEN"
	HeaderBaseURL            = "X-COINDCX-BASE-URL"
	HeaderPublicBaseURL      = "X-COINDCX-PUBLIC-BASE-URL"
)

type Config struct {
	BaseURL      string
	PublicURL    string
	AdminURL     string
	AdminKey     string
	AdminSecret  string
	AgentID      string
	AgentEmail   string
	TargetUserID string
	AuthToken    string
	MCPHTTPPort  int
	MCPHTTPHost  string
}

type AuditLogs struct {
	Enabled       bool
	KafkaBrokers  string
	KafkaTopic    string
	KafkaClientID string
	ConsumerGroup string
	DBDsn         string
}

func FromRequest(r *http.Request) *Config {
	baseURL := strings.TrimSpace(r.Header.Get(HeaderBaseURL))
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	publicURL := strings.TrimSpace(r.Header.Get(HeaderPublicBaseURL))
	if publicURL == "" {
		publicURL = DefaultPublicURL
	}
	adminURL := strings.TrimSpace(r.Header.Get(HeaderAdminBaseURL))
	if adminURL == "" {
		adminURL = DefaultAdminURL
	}
	authToken := strings.TrimSpace(r.Header.Get(HeaderAuthToken))
	if authToken == "" {
		authToken = bearerTokenFromHeader(r.Header.Get("Authorization"))
	}
	cfg := &Config{
		BaseURL:      baseURL,
		PublicURL:    publicURL,
		AdminURL:     adminURL,
		AdminKey:     strings.TrimSpace(r.Header.Get(HeaderAdminServiceKey)),
		AdminSecret:  strings.TrimSpace(r.Header.Get(HeaderAdminServiceSecret)),
		AgentID:      strings.TrimSpace(r.Header.Get(HeaderAdminAgentID)),
		AgentEmail:   strings.TrimSpace(r.Header.Get(HeaderAdminAgentEmail)),
		TargetUserID: strings.TrimSpace(r.Header.Get(HeaderAdminTargetUserID)),
		AuthToken:    authToken,
	}
	logger.Debug("config from request: baseURL=%s adminURL=%s agentID=%s hasAdminKey=%v hasAuthToken=%v",
		baseURL, adminURL, cfg.AgentID, cfg.AdminKey != "", cfg.AuthToken != "")
	return cfg
}

func bearerTokenFromHeader(auth string) string {
	auth = strings.TrimSpace(auth)
	const prefix = "Bearer "
	if len(auth) > len(prefix) && strings.EqualFold(auth[:len(prefix)], prefix) {
		return strings.TrimSpace(auth[len(prefix):])
	}
	return ""
}

func Load() *Config {
	mcpHost := os.Getenv("MCP_HTTP_HOST")
	if mcpHost == "" {
		mcpHost = DefaultMCPHTTPHost
	}
	mcpPort := 0
	if p := os.Getenv("MCP_HTTP_PORT"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 && n < 65536 {
			mcpPort = n
		} else {
			logger.Error("invalid MCP_HTTP_PORT=%q, ignoring", p)
		}
	}
	if mcpPort == 0 {
		mcpPort = 8080
	}
	logger.Debug("config load: MCP_HTTP_HOST=%s MCP_HTTP_PORT=%d", mcpHost, mcpPort)
	return &Config{MCPHTTPPort: mcpPort, MCPHTTPHost: mcpHost}
}

func (c *Config) HasAdminCredentials() bool {
	return c.AdminKey != "" && c.AdminSecret != ""
}

func (c *Config) HasAuthToken() bool {
	return c.AuthToken != ""
}

func LoadAuditLogs() *AuditLogs {
	topic := os.Getenv("AUDIT_KAFKA_TOPIC")
	if topic == "" {
		topic = DefaultAuditKafkaTopic
	}
	group := os.Getenv("AUDIT_CONSUMER_GROUP")
	if group == "" {
		group = DefaultAuditConsumerGroup
	}
	return &AuditLogs{
		Enabled:       strings.EqualFold(os.Getenv("AUDIT_LOGS_ENABLED"), "true"),
		KafkaBrokers:  os.Getenv("AUDIT_KAFKA_BROKERS"),
		KafkaTopic:    topic,
		KafkaClientID: os.Getenv("AUDIT_KAFKA_CLIENT_ID"),
		ConsumerGroup: group,
		DBDsn:         os.Getenv("AUDIT_DB_DSN"),
	}
}

func (a *AuditLogs) IsReady() bool {
	return a.Enabled && a.KafkaBrokers != ""
}
