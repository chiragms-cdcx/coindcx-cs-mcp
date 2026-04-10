package tools

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const futuresAuthTokenRequired = "Auth token required for futures: set X-COINDCX-AUTH-TOKEN header"

func futuresBody(extra map[string]any) map[string]any {
	b := map[string]any{"timestamp": time.Now().UnixMilli()}
	for k, v := range extra {
		b[k] = v
	}
	return b
}

type FuturesOrdersInput struct {
	Status         string `json:"status" jsonschema:"Comma-separated: open, filled, partially_filled, cancelled, etc."`
	Side           string `json:"side" jsonschema:"buy or sell"`
	Page           string `json:"page" jsonschema:"Page number"`
	Size           string `json:"size" jsonschema:"Page size"`
	MarginCurrency string `json:"margin_currency_short_name,omitempty" jsonschema:"USDT or INR (optional)"`
}

func FuturesOrders(ctx context.Context, req *mcp.CallToolRequest, input FuturesOrdersInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(futuresAuthTokenRequired)
	}
	if input.Status == "" || input.Side == "" || input.Page == "" || input.Size == "" {
		return toolError("status, side, page and size are required")
	}
	side := strings.ToLower(strings.TrimSpace(input.Side))
	if side != "buy" && side != "sell" {
		return toolError("side must be buy or sell")
	}
	body := futuresBody(map[string]any{"status": input.Status, "side": side, "page": input.Page, "size": input.Size})
	marginCurrency := input.MarginCurrency
	if marginCurrency == "" {
		marginCurrency = "USDT"
	}
	body["margin_currency_short_name"] = []string{marginCurrency}
	data, status, err := c.FuturesPostPrivate("orders", body)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResultFromCall(data, status, nil)
}

type FuturesPositionsInput struct {
	Page           string `json:"page,omitempty" jsonschema:"Page number (default 1)"`
	Size           string `json:"size,omitempty" jsonschema:"Page size (default 10)"`
	Pairs          string `json:"pairs,omitempty" jsonschema:"Comma-separated pairs"`
	PositionIDs    string `json:"position_ids,omitempty" jsonschema:"Comma-separated position IDs"`
	MarginCurrency string `json:"margin_currency_short_name,omitempty" jsonschema:"USDT or INR"`
}

func FuturesPositions(ctx context.Context, req *mcp.CallToolRequest, input FuturesPositionsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(futuresAuthTokenRequired)
	}
	body := futuresBody(nil)
	if input.Page != "" {
		body["page"] = input.Page
	} else {
		body["page"] = "1"
	}
	if input.Size != "" {
		body["size"] = input.Size
	} else {
		body["size"] = "10"
	}
	if input.Pairs != "" {
		body["pairs"] = input.Pairs
	}
	if input.PositionIDs != "" {
		body["position_ids"] = input.PositionIDs
	}
	if input.MarginCurrency != "" {
		body["margin_currency_short_name"] = []string{input.MarginCurrency}
	} else {
		body["margin_currency_short_name"] = []string{"USDT"}
	}
	data, status, err := c.FuturesPostPrivate("positions", body)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	if status == 404 {
		return toolResult([]byte("[]"), false)
	}
	return toolResultFromCall(data, status, nil)
}

type FuturesPositionTransactionsInput struct {
	Stage          string `json:"stage" jsonschema:"all, default, funding, exit, tpsl_exit, liquidation"`
	Page           string `json:"page" jsonschema:"Page number"`
	Size           string `json:"size" jsonschema:"Page size"`
	MarginCurrency string `json:"margin_currency_short_name,omitempty" jsonschema:"USDT or INR (optional)"`
}

func FuturesPositionTransactions(ctx context.Context, req *mcp.CallToolRequest, input FuturesPositionTransactionsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(futuresAuthTokenRequired)
	}
	if input.Stage == "" || input.Page == "" || input.Size == "" {
		return toolError("stage, page and size are required")
	}
	body := futuresBody(map[string]any{"stage": input.Stage, "page": input.Page, "size": input.Size})
	if input.MarginCurrency != "" {
		body["margin_currency_short_name"] = []string{input.MarginCurrency}
	}
	data, status, err := c.FuturesPostPrivate("positions/transactions", body)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResultFromCall(data, status, nil)
}

type FuturesUserTradesInput struct {
	Pair           string `json:"pair" jsonschema:"e.g. KC-BTC_USDT"`
	FromDate       string `json:"from_date" jsonschema:"YYYY-MM-DD"`
	ToDate         string `json:"to_date" jsonschema:"YYYY-MM-DD"`
	Page           string `json:"page" jsonschema:"Page number"`
	Size           string `json:"size" jsonschema:"Page size"`
	MarginCurrency string `json:"margin_currency_short_name,omitempty" jsonschema:"USDT or INR (optional)"`
}

func FuturesUserTrades(ctx context.Context, req *mcp.CallToolRequest, input FuturesUserTradesInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(futuresAuthTokenRequired)
	}
	if input.Pair == "" || input.FromDate == "" || input.ToDate == "" || input.Page == "" || input.Size == "" {
		return toolError("pair, from_date, to_date, page and size are required")
	}
	body := futuresBody(map[string]any{"pair": input.Pair, "from_date": input.FromDate, "to_date": input.ToDate, "page": input.Page, "size": input.Size})
	if input.MarginCurrency != "" {
		body["margin_currency_short_name"] = []string{input.MarginCurrency}
	} else {
		body["margin_currency_short_name"] = []string{"USDT"}
	}
	data, status, err := c.FuturesPostPrivate("trades", body)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResultFromCall(data, status, nil)
}

type FuturesStatsInput struct {
	Pair string `json:"pair" jsonschema:"e.g. KC-ETH_USDT"`
}

func FuturesStats(ctx context.Context, req *mcp.CallToolRequest, input FuturesStatsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(futuresAuthTokenRequired)
	}
	if input.Pair == "" {
		return toolError("pair is required")
	}
	body := futuresBody(nil)
	data, status, err := c.FuturesPostPrivate("data/stats?pair="+url.QueryEscape(input.Pair), body)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResultFromCall(data, status, nil)
}

type FuturesCrossMarginDetailsInput struct{}

func FuturesCrossMarginDetails(ctx context.Context, req *mcp.CallToolRequest, _ FuturesCrossMarginDetailsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(futuresAuthTokenRequired)
	}
	body := futuresBody(nil)
	data, status, err := c.FuturesPostPrivate("positions/cross_margin_details", body)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	if status == 404 {
		return toolResult([]byte("{}"), false)
	}
	return toolResultFromCall(data, status, nil)
}

type FuturesWalletsInput struct{}

func FuturesWallets(ctx context.Context, req *mcp.CallToolRequest, _ FuturesWalletsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(futuresAuthTokenRequired)
	}
	data, status, err := c.FuturesGetPrivate("wallets", nil)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResultFromCall(data, status, nil)
}

type FuturesWalletTransactionsInput struct {
	Page string `json:"page,omitempty" jsonschema:"Page number (default 1)"`
	Size string `json:"size,omitempty" jsonschema:"Page size (default 1000)"`
}

func FuturesWalletTransactions(ctx context.Context, req *mcp.CallToolRequest, input FuturesWalletTransactionsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(futuresAuthTokenRequired)
	}
	q := url.Values{}
	if input.Page != "" {
		q.Set("page", input.Page)
	} else {
		q.Set("page", "1")
	}
	if input.Size != "" {
		q.Set("size", input.Size)
	} else {
		q.Set("size", "1000")
	}
	data, status, err := c.FuturesGetPrivate("wallets/transactions", q)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResultFromCall(data, status, nil)
}
