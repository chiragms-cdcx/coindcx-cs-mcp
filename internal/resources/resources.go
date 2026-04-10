package resources

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const scheme = "coindcx-cs"

func All() ([]*mcp.Resource, mcp.ResourceHandler) {
	handler := resourceHandler(contents)
	return list, handler
}

var list = []*mcp.Resource{
	{URI: scheme + "://docs/tools", Name: "tools", Title: "CS MCP Tool Reference", Description: "All available tools by category (admin, KB, market data, read-only trading).", MIMEType: "text/plain"},
	{URI: scheme + "://docs/auth", Name: "auth", Title: "Authentication Guide", Description: "Admin service credentials and auth token setup for CS agents.", MIMEType: "text/plain"},
	{URI: scheme + "://docs/cs-workflow", Name: "cs-workflow", Title: "CS Agent Workflow", Description: "Step-by-step workflow for investigating and resolving customer issues.", MIMEType: "text/plain"},
	{URI: scheme + "://docs/agent-guidelines", Name: "agent-guidelines", Title: "AI Agent Guidelines", Description: "Do's and don'ts for the CS AI agent.", MIMEType: "text/plain"},
	{URI: scheme + "://docs/common-issues", Name: "common-issues", Title: "Common Issue Patterns", Description: "Top CRT categories and resolution patterns.", MIMEType: "text/plain"},
}

var contents = map[string]string{
	scheme + "://docs/tools": `CoinDCX CS MCP – Tool Reference (59 read-only tools + 13 admin + 5 KB = 77 total)

Admin User Lookup (13 tools, require X-ADMIN-SERVICE-KEY/SECRET):
admin_lookup_user, admin_get_user_profile, admin_get_user_balances, admin_get_user_orders,
admin_get_user_deposits, admin_get_user_withdrawals, admin_get_user_transactions,
admin_get_user_kyc_details, admin_get_user_compliance, admin_get_user_futures_positions,
admin_get_user_options_positions, admin_get_user_login_history, admin_get_user_referrals.

Knowledge Base (5 tools):
kb_search_known_issues, kb_get_sop, kb_get_standard_response, kb_search_past_crts, kb_get_crt_details.

Market Data – Public (18 tools, no auth):
Spot: coindcx_get_ticker, coindcx_get_markets, coindcx_get_markets_details, coindcx_get_orderbook, coindcx_get_candles.
Futures: coindcx_futures_active_instruments, coindcx_futures_instrument, coindcx_futures_trades, coindcx_futures_current_prices, coindcx_futures_orderbook.
Options: coindcx_options_get_base_currencies, get_instruments, get_ticker, get_orderbook, get_trades, get_stats, get_combined_instruments, get_index_price.

User Data – Read-Only (25 tools, require auth token or admin creds):
Spot: get_balances, get_user_info, get_order_status, get_order_status_multiple, get_active_orders, get_spot_orders_v2, get_order_trade_history, get_active_orders_count.
Futures: futures_orders, futures_positions, futures_position_transactions, futures_user_trades, futures_stats, futures_cross_margin_details, futures_wallets, futures_wallet_transactions.
Options: options_get_orders, get_order_history, get_trades_by_order_id, get_positions, get_wallet_balance, get_closed_positions, get_margin_mode, get_wallet_consumed_limit, get_wallet_transactions.

News (2 tools): coindcx_get_crypto_news, coindcx_get_crypto_news_summary.
Batch: coindcx_request (up to 20 concurrent read-only calls).

NO WRITE TOOLS. This server is strictly read-only. CS agents cannot place orders, cancel orders, transfer funds, or modify any user state through this MCP server.`,

	scheme + "://docs/auth": `CoinDCX CS MCP – Authentication

This MCP server runs in HTTP mode only. All credentials come from request headers.

Admin API access (required for admin_* and kb_* tools):
- X-ADMIN-SERVICE-KEY: Admin service account key
- X-ADMIN-SERVICE-SECRET: Admin service account secret
- X-ADMIN-AGENT-ID: CS agent's employee ID (for audit)
- X-ADMIN-AGENT-EMAIL: CS agent's email (for audit)

Trading data access (required for user read-only and futures/options private tools):
- X-COINDCX-AUTH-TOKEN: Bearer token for the target user's data

Public tools (market data, news) require no authentication.

All tool calls are audit-logged with the agent's identity and the target user.`,

	scheme + "://docs/cs-workflow": `CoinDCX CS MCP – CS Agent Workflow

1. Customer files a complaint.
2. CS agent opens the AI Chat in Admin Dashboard.
3. Use cs_user_overview prompt to get customer context.
4. Use cs_known_issue_check to see if the symptom is a known issue.
5. If known issue: follow the SOP and use the standard response template.
6. If unknown: use specific investigation prompts (cs_order_investigation, cs_deposit_investigation, etc.).
7. Use cs_generate_response to draft the customer reply.
8. Resolve the ticket. No CRT needed for engineering.

Target: resolve Q&A issues in under 5 minutes (down from 1-3 days).`,

	scheme + "://docs/agent-guidelines": `CoinDCX CS MCP – AI Agent Guidelines

Don'ts:
- NEVER expose customer PII (phone numbers, email addresses, account details) in logs or responses beyond what's needed.
- NEVER dump raw API JSON to the CS agent. Always summarize.
- NEVER attempt to modify user state. This server has NO write tools.
- NEVER make up data. If an API call fails, say so clearly.
- NEVER share internal system details with customers (API paths, error codes, internal IDs).

Do's:
- Always cite the data source (which tool returned the data).
- Present findings in clear, concise summaries with relevant details only.
- Handle errors gracefully: "Could not fetch deposit history – check admin credentials."
- Use the knowledge base tools to check for known issues before deep investigation.
- Generate customer-friendly responses that are empathetic and clear.
- Log all lookups via audit for compliance.`,

	scheme + "://docs/common-issues": `CoinDCX CS MCP – Common Issue Patterns (from CRT analysis Jan-Apr 2026)

Top categories by volume:
1. Futures Trading (~30%): Stuck orders, negative balance, liquidation RCA. Check: admin_get_user_futures_positions, futures_orders.
2. Spot Trading (15-20%): Order cancellation stuck, locked balance. Check: admin_get_user_orders, get_active_orders.
3. Charts Issues (~10%): Charts not loading. Usually self-resolves. Check kb_search_known_issues.
4. Deposits (~10%): Not reflecting. Check: admin_get_user_deposits, admin_get_user_compliance.
5. Withdrawals (~5%): Rejected or stuck. Check: admin_get_user_withdrawals, admin_get_user_compliance.
6. KYC (~5%): Stuck, selfie issues. Check: admin_get_user_kyc_details.
7. Web3/DeFi (~5%): Reset requests, PnL display. Check admin tools.
8. Options (~3%): Liquidation, position not reflecting. Check: admin_get_user_options_positions.
9. API Issues (~5%): Unable to cancel via API. Check: admin_get_user_orders.

Systemic issues:
- GS Futures Stuck: 16+ identical tickets with same resolution. Use kb_search_known_issues.
- Charts Loading: 18+ tickets, many self-resolve. Usually Agent_Error.
- Common push-back from tech leads: "This can be checked via Admin Dashboard."`,
}

func resourceHandler(contentByURI map[string]string) mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		uri := req.Params.URI
		u, err := url.Parse(uri)
		if err != nil || u.Scheme != scheme {
			return nil, mcp.ResourceNotFoundError(uri)
		}
		text, ok := contentByURI[uri]
		if !ok {
			return nil, mcp.ResourceNotFoundError(uri)
		}
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{URI: uri, MIMEType: "text/plain", Text: text},
			},
		}, nil
	}
}
