# CoinDCX CS AI Agent — MCP Server

An MCP (Model Context Protocol) server in Go that powers the CoinDCX Customer Support AI Agent. Built for the Admin Dashboard, it lets CS agents query customer data in natural language, diagnose issues, and resolve CRT tickets without engineering involvement.

**Based on:** [CoinDCX MCP Hackathon Server](https://github.com/coindcx-hackathon-2026/team10-Growth-Momemtum-Capture-Protocol)  
**PRD:** [CoinDCX CS AI Agent — MCP-Powered Customer Support Intelligence](https://coindcx.atlassian.net/wiki/spaces/CefiOption/pages/4204429674)

---

## Architecture

```
CS Agent (browser)
    ↓
Admin Dashboard (React)
    ↓
LLM Gateway (AI Proxy — Claude Opus 4.5 via AWS Bedrock)
    ↓
MCP Server (this repo — Go, HTTP mode)
    ├── CoinDCX Admin APIs (user lookup, balances, orders, KYC, compliance)
    ├── CoinDCX Public APIs (market data, prices, orderbook)
    ├── Knowledge Base (known issues, SOPs, past CRTs)
    └── Kafka → Postgres (audit log)
```

## Tools Summary (77 total, all read-only)

| Category | Count | Description |
|----------|-------|-------------|
| Admin User Lookup | 13 | User profile, balances, orders, deposits, withdrawals, KYC, compliance, positions, login history, referrals |
| Knowledge Base | 5 | Known issues search, SOPs, standard responses, past CRT search |
| Spot Public | 5 | Ticker, markets, orderbook, candles |
| Spot User Read | 8 | Balances, order status, active orders, trade history |
| Futures Public | 5 | Instruments, prices, orderbook, trades |
| Futures Private Read | 8 | Orders, positions, wallets, trades, cross margin |
| Options Public | 8 | Instruments, ticker, orderbook, stats |
| Options Private Read | 9 | Orders, positions, wallet, closed positions, margin mode |
| News | 2 | Crypto news from CoinDesk, CoinTelegraph, Decrypt |
| Batch | 1 | Run up to 20 tools concurrently |

**No write tools.** CS agents cannot place orders, cancel orders, transfer funds, or modify any user state through this server.

## CS-Specific Prompts (10)

| Prompt | Description |
|--------|-------------|
| `cs_investigate_crt` | End-to-end CRT ticket investigation |
| `cs_user_overview` | Complete customer overview (profile, KYC, balances, compliance) |
| `cs_order_investigation` | Investigate order issues (stuck, cancelled, not filled) |
| `cs_deposit_investigation` | Investigate deposits not reflecting |
| `cs_withdrawal_investigation` | Investigate stuck/failed withdrawals |
| `cs_futures_investigation` | Investigate futures issues (stuck orders, liquidation) |
| `cs_options_investigation` | Investigate options issues |
| `cs_known_issue_check` | Check if symptom matches known issue with SOP |
| `cs_compliance_check` | Check compliance and KYC status |
| `cs_generate_response` | Generate customer-facing response |

---

## Quick Start

### Prerequisites

- Go 1.24+
- (Optional) Kafka + PostgreSQL for audit logs

### Build and Run

```bash
# Clone the repo
git clone https://github.com/chiragms-cdcx/coindcx-cs-mcp.git
cd coindcx-cs-mcp

# Copy env file
cp .env.example .env

# Build
make build

# Run (HTTP on port 8080)
make run

# Or run without building
make go-run
```

The server starts on `http://0.0.0.0:8080`.

---

## MCP Integration Steps

### 1. Add to Cursor IDE

Add this to your `.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "coindcx-cs": {
      "url": "http://localhost:8080",
      "headers": {
        "X-ADMIN-SERVICE-KEY": "your-admin-service-key",
        "X-ADMIN-SERVICE-SECRET": "your-admin-service-secret",
        "X-ADMIN-AGENT-ID": "agent-001",
        "X-ADMIN-AGENT-EMAIL": "agent@coindcx.com",
        "X-COINDCX-AUTH-TOKEN": "bearer-token-for-user-data"
      }
    }
  }
}
```

### 2. Add to Claude Desktop

Use [mcp-remote](https://www.npmjs.com/package/mcp-remote) to bridge HTTP to stdio:

```json
{
  "mcpServers": {
    "coindcx-cs": {
      "command": "npx",
      "args": [
        "mcp-remote",
        "http://localhost:8080",
        "--header", "X-ADMIN-SERVICE-KEY: your-key",
        "--header", "X-ADMIN-SERVICE-SECRET: your-secret",
        "--header", "X-ADMIN-AGENT-ID: agent-001",
        "--header", "X-ADMIN-AGENT-EMAIL: agent@coindcx.com",
        "--header", "X-COINDCX-AUTH-TOKEN: your-token"
      ]
    }
  }
}
```

### 3. Integrate with Admin Dashboard (Production)

The Admin Dashboard frontend connects to the MCP server via the LLM Gateway:

```
Admin Dashboard → LLM Gateway → MCP Server (this)
```

The LLM Gateway (AI Proxy) should:
1. Accept chat messages from the Admin Dashboard
2. Forward them to Claude Opus 4.5 (via AWS Bedrock) with MCP tool definitions
3. Execute tool calls against this MCP server
4. Return responses to the dashboard

Headers are set by the Admin Dashboard backend (not the browser):
- `X-ADMIN-SERVICE-KEY` / `X-ADMIN-SERVICE-SECRET`: Admin service account credentials
- `X-ADMIN-AGENT-ID` / `X-ADMIN-AGENT-EMAIL`: CS agent identity (from SSO)
- `X-ADMIN-TARGET-USER-ID`: Customer being looked up
- `X-COINDCX-AUTH-TOKEN`: Auth token for CoinDCX private APIs

---

## Adding APIs to the MCP Server

### Adding a New Admin Tool

1. **Define the input struct** in `internal/tools/admin.go`:

```go
type AdminGetUserNewDataInput struct {
    UserID string `json:"user_id" jsonschema:"User ID"`
}
```

2. **Implement the handler**:

```go
func AdminGetUserNewData(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserNewDataInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
    if !c.HasCredentials() {
        return toolError(adminCredsRequired)
    }
    if input.UserID == "" {
        return toolError("user_id is required")
    }
    q := url.Values{}
    q.Set("user_id", input.UserID)
    data, status, err := c.AdminGet("/admin/api/v1/users/new-data", q)
    return toolResultFromCall(data, status, err)
}
```

3. **Register in `main.go`**:

```go
addTool(server, "admin_get_user_new_data", "Description of the tool. (Admin)", tools.AdminGetUserNewData, c, pub)
```

### Adding a New Knowledge Base Tool

Same pattern as admin tools, but use the `/admin/api/v1/kb/` endpoint prefix.

### Adding a New Public Market Data Tool

For public (no-auth) tools, use `c.GetPublic()` or `c.GetBase()` and follow the pattern in `internal/tools/public.go`.

---

## Audit Logging

All tool calls are audit-logged with:
- Tool name and parameters
- CS agent identity (agent_id, agent_email)
- Target user being looked up
- Session ID for correlation
- Latency and status

### Setup Audit Logging

```bash
# 1. Set up PostgreSQL database
createdb auditdb

# 2. Run migration
AUDIT_DB_DSN=postgres://user:pass@localhost/auditdb make migrate

# 3. Enable in .env
AUDIT_LOGS_ENABLED=true
AUDIT_KAFKA_BROKERS=localhost:9092

# 4. Start audit consumer
make run-audit-consumer
```

---

## Project Structure

```
.
├── cmd/
│   ├── coindcx-cs-mcp/    # MCP server binary (HTTP mode)
│   ├── audit-consumer/     # Kafka → Postgres audit consumer
│   └── audit-migrate/      # Database schema migration
├── internal/
│   ├── config/             # Configuration (admin auth, audit, server)
│   ├── client/             # HTTP client for CoinDCX + Admin APIs
│   ├── tools/              # MCP tool handlers
│   │   ├── admin.go        # Admin user lookup tools (13)
│   │   ├── knowledge.go    # Knowledge base tools (5)
│   │   ├── public.go       # Spot public tools
│   │   ├── user_read.go    # Spot user read-only tools
│   │   ├── futures_*.go    # Futures tools (public + private read)
│   │   ├── options_*.go    # Options tools (public + private read)
│   │   ├── news.go         # Crypto news tools
│   │   └── helpers.go      # Batch, wrap, client interface
│   ├── prompts/            # CS-specific MCP prompts (10)
│   ├── resources/          # MCP resources (docs, guidelines)
│   └── logger/             # Structured logger
├── pkg/
│   ├── audit/              # Audit event types + Kafka producer
│   └── auditstore/         # Postgres consumer + migrations
├── Makefile
├── .env.example
└── README.md
```

---

## Key Differences from Hackathon Server

| Aspect | Hackathon Server | CS MCP Server |
|--------|-----------------|---------------|
| Transport | HTTP + Stdio | HTTP only |
| Auth | User API keys | Admin service credentials |
| Write tools | 25+ order/transfer tools | None (read-only) |
| Admin tools | None | 13 user lookup tools |
| Knowledge base | None | 5 KB tools |
| Prompts | Trading-focused (12) | CS-focused (10) |
| Audit | Basic tool logging | Agent identity + target user tracking |
| Elicitation | Write tool confirmation | Removed (no write tools) |

---

## License

Internal — CoinDCX
