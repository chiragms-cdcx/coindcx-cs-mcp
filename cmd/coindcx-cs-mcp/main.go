package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/client"
	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/config"
	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/logger"
	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/prompts"
	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/resources"
	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/tools"
	"github.com/chiragms-cdcx/coindcx-cs-mcp/pkg/audit"
)

var schemaCache = mcp.NewSchemaCache()

func main() {
	serverCfg := config.Load()

	auditCfg := config.LoadAuditLogs()
	var pub audit.Publisher = audit.NoopPublisher{}
	if auditCfg.IsReady() {
		brokers := strings.Split(auditCfg.KafkaBrokers, ",")
		pub = audit.NewKafkaProducer(brokers, auditCfg.KafkaTopic, auditCfg.KafkaClientID)
		log.Printf("audit logs enabled: kafka brokers=%s topic=%s", auditCfg.KafkaBrokers, auditCfg.KafkaTopic)
	}
	defer pub.Close() //nolint:errcheck

	logger.Info("starting CoinDCX CS MCP server in HTTP mode")
	serverOpts := &mcp.ServerOptions{SchemaCache: schemaCache}
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		cfg := config.FromRequest(r)
		c := client.New(cfg.BaseURL, cfg.PublicURL, cfg.AdminURL, cfg.AdminKey, cfg.AdminSecret, cfg.AuthToken)
		return newServer(c, serverOpts, pub)
	}, nil)

	_ = context.Background()
	addr := net.JoinHostPort(serverCfg.MCPHTTPHost, fmt.Sprintf("%d", serverCfg.MCPHTTPPort))
	logger.Info("CS MCP server listening on http://%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		logger.Error("HTTP server failed: %v", err, serverCfg.MCPHTTPPort)
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func addTool[T any](
	server *mcp.Server,
	name, desc string,
	fn func(context.Context, *mcp.CallToolRequest, T, tools.CoinDCXClient) (*mcp.CallToolResult, any, error),
	c tools.CoinDCXClient,
	pub audit.Publisher,
) {
	tools.RegisterBatchFactory(name, tools.MakeBatchHandlerFactory(fn))
	mcp.AddTool(server, &mcp.Tool{Name: name, Description: desc}, tools.Wrap(fn, c))
}

func addToolNoBatch[T any](
	server *mcp.Server,
	name, desc string,
	fn func(context.Context, *mcp.CallToolRequest, T, tools.CoinDCXClient) (*mcp.CallToolResult, any, error),
	c tools.CoinDCXClient,
	pub audit.Publisher,
) {
	mcp.AddTool(server, &mcp.Tool{Name: name, Description: desc}, tools.WrapWithAudit(fn, c, pub, audit.TransportHTTP))
}

func newServer(c tools.CoinDCXClient, opts *mcp.ServerOptions, pub audit.Publisher) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "coindcx-cs-mcp",
		Version: "1.0.0",
	}, opts)

	for _, p := range prompts.All() {
		server.AddPrompt(p.Prompt, p.Handler)
	}

	resourceList, resourceHandler := resources.All()
	for _, r := range resourceList {
		server.AddResource(r, resourceHandler)
	}

	// ----- Spot public (5 tools) -----
	addTool(server, "coindcx_get_ticker", "Get 24h ticker for all markets. (Public)", tools.GetTicker, c, pub)
	addTool(server, "coindcx_get_markets", "List active market symbols. (Public)", tools.GetMarkets, c, pub)
	addTool(server, "coindcx_get_markets_details", "Get full market details. (Public)", tools.GetMarketsDetails, c, pub)
	addTool(server, "coindcx_get_orderbook", "Get order book for a pair. (Public)", tools.GetOrderbook, c, pub)
	addTool(server, "coindcx_get_candles", "Get OHLCV candles for a pair and interval. (Public)", tools.GetCandles, c, pub)

	// ----- Spot user read-only (8 tools) -----
	addTool(server, "coindcx_get_balances", "Get user spot balances. (Read-only)", tools.GetBalances, c, pub)
	addTool(server, "coindcx_get_user_info", "Get user profile info. (Read-only)", tools.GetUserInfo, c, pub)
	addTool(server, "coindcx_get_order_status", "Get status of one order. (Read-only)", tools.GetOrderStatus, c, pub)
	addTool(server, "coindcx_get_order_status_multiple", "Get status of multiple orders. (Read-only)", tools.GetOrderStatusMultiple, c, pub)
	addTool(server, "coindcx_get_active_orders", "Get user's active orders. (Read-only)", tools.GetActiveOrders, c, pub)
	addTool(server, "coindcx_get_spot_orders_v2", "Get spot orders (history or pending). (Read-only)", tools.GetSpotOrdersV2, c, pub)
	addTool(server, "coindcx_get_order_trade_history", "Get user's order trade history. (Read-only)", tools.GetOrderTradeHistory, c, pub)
	addTool(server, "coindcx_get_active_orders_count", "Get count of active orders. (Read-only)", tools.GetActiveOrdersCount, c, pub)

	// ----- Batch (async multi-tool call) -----
	addToolNoBatch(server, "coindcx_request", "Run multiple read-only tool calls concurrently. Params: calls (array of {name, arguments}). Max 20.", tools.Batch, c, pub)

	// ----- Futures public (5 tools) -----
	addTool(server, "coindcx_futures_active_instruments", "List active futures instruments. (Public)", tools.FuturesActiveInstruments, c, pub)
	addTool(server, "coindcx_futures_instrument", "Get details for one futures instrument. (Public)", tools.FuturesInstrument, c, pub)
	addTool(server, "coindcx_futures_trades", "Recent public trades for a futures pair. (Public)", tools.FuturesTrades, c, pub)
	addTool(server, "coindcx_futures_current_prices", "Real-time prices for all futures. (Public)", tools.FuturesCurrentPrices, c, pub)
	addTool(server, "coindcx_futures_orderbook", "Order book for a futures pair. (Public)", tools.FuturesOrderbook, c, pub)

	// ----- Futures private read-only (8 tools) -----
	addTool(server, "coindcx_futures_orders", "List futures orders by status/side. (Read-only)", tools.FuturesOrders, c, pub)
	addTool(server, "coindcx_futures_positions", "List futures positions. (Read-only)", tools.FuturesPositions, c, pub)
	addTool(server, "coindcx_futures_position_transactions", "Position transactions. (Read-only)", tools.FuturesPositionTransactions, c, pub)
	addTool(server, "coindcx_futures_user_trades", "User's futures trade history. (Read-only)", tools.FuturesUserTrades, c, pub)
	addTool(server, "coindcx_futures_stats", "Stats for a futures pair. (Read-only)", tools.FuturesStats, c, pub)
	addTool(server, "coindcx_futures_cross_margin_details", "Cross margin details. (Read-only)", tools.FuturesCrossMarginDetails, c, pub)
	addTool(server, "coindcx_futures_wallets", "Futures wallet balances. (Read-only)", tools.FuturesWallets, c, pub)
	addTool(server, "coindcx_futures_wallet_transactions", "Futures wallet transaction history. (Read-only)", tools.FuturesWalletTransactions, c, pub)

	// ----- Options public (8 tools) -----
	addTool(server, "coindcx_options_get_base_currencies", "Options base currencies. (Public)", tools.OptionsGetBaseCurrencies, c, pub)
	addTool(server, "coindcx_options_get_instruments", "Options instruments metadata. (Public)", tools.OptionsGetInstruments, c, pub)
	addTool(server, "coindcx_options_get_ticker", "Options ticker data. (Public)", tools.OptionsGetTicker, c, pub)
	addTool(server, "coindcx_options_get_orderbook", "Options orderbook. (Public)", tools.OptionsGetOrderbook, c, pub)
	addTool(server, "coindcx_options_get_trades", "Options recent trades. (Public)", tools.OptionsGetTrades, c, pub)
	addTool(server, "coindcx_options_get_stats", "Options statistics by base currency. (Public)", tools.OptionsGetStats, c, pub)
	addTool(server, "coindcx_options_get_combined_instruments", "Combined instrument details. (Public)", tools.OptionsGetCombinedInstruments, c, pub)
	addTool(server, "coindcx_options_get_index_price", "Index price for a base currency. (Public)", tools.OptionsGetIndexPrice, c, pub)

	// ----- Options private read-only (9 tools) -----
	addTool(server, "coindcx_options_get_orders", "Open options orders. (Read-only)", tools.OptionsGetOrders, c, pub)
	addTool(server, "coindcx_options_get_order_history", "Options order history. (Read-only)", tools.OptionsGetOrderHistory, c, pub)
	addTool(server, "coindcx_options_get_trades_by_order_id", "Trades for a specific order. (Read-only)", tools.OptionsGetTradesByOrderID, c, pub)
	addTool(server, "coindcx_options_get_positions", "Open options positions. (Read-only)", tools.OptionsGetPositions, c, pub)
	addTool(server, "coindcx_options_get_wallet_balance", "Options wallet balance. (Read-only)", tools.OptionsGetWalletBalance, c, pub)
	addTool(server, "coindcx_options_get_closed_positions", "Closed options positions. (Read-only)", tools.OptionsGetClosedPositions, c, pub)
	addTool(server, "coindcx_options_get_margin_mode", "Current margin mode. (Read-only)", tools.OptionsGetMarginMode, c, pub)
	addTool(server, "coindcx_options_get_wallet_consumed_limit", "Daily consumed transfer limit. (Read-only)", tools.OptionsGetWalletConsumedLimit, c, pub)
	addTool(server, "coindcx_options_get_wallet_transactions", "Options wallet transactions. (Read-only)", tools.OptionsGetWalletTransactions, c, pub)

	// ----- Crypto news (2 tools) -----
	addTool(server, "coindcx_get_crypto_news", "Fetch crypto news from credible sources. (Public)", tools.GetCryptoNews, c, pub)
	addTool(server, "coindcx_get_crypto_news_summary", "Brief crypto headlines. (Public)", tools.GetCryptoNewsSummary, c, pub)

	// ----- Admin user lookup (13 tools) -----
	addTool(server, "admin_lookup_user", "Search user by email, phone, user_id, or UID. (Admin)", tools.AdminLookupUser, c, pub)
	addTool(server, "admin_get_user_profile", "Full user profile including KYC status. (Admin)", tools.AdminGetUserProfile, c, pub)
	addTool(server, "admin_get_user_balances", "All wallet balances for a user. (Admin)", tools.AdminGetUserBalances, c, pub)
	addTool(server, "admin_get_user_orders", "Order history across spot, futures, options. (Admin)", tools.AdminGetUserOrders, c, pub)
	addTool(server, "admin_get_user_deposits", "Deposit history with status. (Admin)", tools.AdminGetUserDeposits, c, pub)
	addTool(server, "admin_get_user_withdrawals", "Withdrawal history with status and failure reasons. (Admin)", tools.AdminGetUserWithdrawals, c, pub)
	addTool(server, "admin_get_user_transactions", "Full transaction log for a user. (Admin)", tools.AdminGetUserTransactions, c, pub)
	addTool(server, "admin_get_user_kyc_details", "KYC verification status and rejection reasons. (Admin)", tools.AdminGetUserKYCDetails, c, pub)
	addTool(server, "admin_get_user_compliance", "Compliance flags, restrictions, AML status. (Admin)", tools.AdminGetUserCompliance, c, pub)
	addTool(server, "admin_get_user_futures_positions", "Current futures positions for a user. (Admin)", tools.AdminGetUserFuturesPositions, c, pub)
	addTool(server, "admin_get_user_options_positions", "Current options positions for a user. (Admin)", tools.AdminGetUserOptionsPositions, c, pub)
	addTool(server, "admin_get_user_login_history", "Recent login sessions, devices, IPs. (Admin)", tools.AdminGetUserLoginHistory, c, pub)
	addTool(server, "admin_get_user_referrals", "Referral status and linked accounts. (Admin)", tools.AdminGetUserReferrals, c, pub)

	// ----- Knowledge base (5 tools) -----
	addTool(server, "kb_search_known_issues", "Search known-issues database by symptom or keyword. (KB)", tools.KBSearchKnownIssues, c, pub)
	addTool(server, "kb_get_sop", "Get SOP for a specific issue type. (KB)", tools.KBGetSOP, c, pub)
	addTool(server, "kb_get_standard_response", "Get pre-approved customer response template. (KB)", tools.KBGetStandardResponse, c, pub)
	addTool(server, "kb_search_past_crts", "Search past CRT tickets for similar issues. (KB)", tools.KBSearchPastCRTs, c, pub)
	addTool(server, "kb_get_crt_details", "Fetch full details of a CRT ticket. (KB)", tools.KBGetCRTDetails, c, pub)

	return server
}
