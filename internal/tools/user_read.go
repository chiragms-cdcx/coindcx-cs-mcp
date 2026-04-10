package tools

import (
	"context"
	"net/url"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetBalancesInput struct{}

func GetBalances(ctx context.Context, req *mcp.CallToolRequest, _ GetBalancesInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError("admin service credentials must be set for this tool")
	}
	body := map[string]int64{"timestamp": time.Now().UnixMilli()}
	data, status, err := c.PostSigned("/exchange/v1/users/balances", body)
	return toolResultFromCall(data, status, err)
}

type GetUserInfoInput struct{}

func GetUserInfo(ctx context.Context, req *mcp.CallToolRequest, _ GetUserInfoInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError("admin service credentials must be set for this tool")
	}
	body := map[string]int64{"timestamp": time.Now().UnixMilli()}
	data, status, err := c.PostSigned("/exchange/v1/users/info", body)
	return toolResultFromCall(data, status, err)
}

type GetOrderStatusInput struct {
	OrderID       string `json:"order_id,omitempty" jsonschema:"System order ID"`
	ClientOrderID string `json:"client_order_id,omitempty" jsonschema:"Client order ID"`
}

func GetOrderStatus(ctx context.Context, req *mcp.CallToolRequest, input GetOrderStatusInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError("admin service credentials must be set for this tool")
	}
	if input.OrderID == "" && input.ClientOrderID == "" {
		return toolError("order_id or client_order_id is required")
	}
	body := map[string]any{"timestamp": time.Now().UnixMilli()}
	if input.OrderID != "" {
		body["id"] = input.OrderID
	}
	if input.ClientOrderID != "" {
		body["client_order_id"] = input.ClientOrderID
	}
	data, status, err := c.PostSigned("/exchange/v1/orders/status", body)
	return toolResultFromCall(data, status, err)
}

type GetActiveOrdersInput struct {
	Market string `json:"market,omitempty" jsonschema:"Filter by market (optional)"`
}

func GetActiveOrders(ctx context.Context, req *mcp.CallToolRequest, input GetActiveOrdersInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() && !c.HasCredentials() {
		return toolError("auth token or admin service credentials required")
	}
	body := map[string]any{"timestamp": time.Now().UnixMilli()}
	if input.Market != "" {
		body["market"] = input.Market
	}
	data, status, err := c.PostSigned("/exchange/v1/orders/active_orders", body)
	return toolResultFromCall(data, status, err)
}

type GetSpotOrdersV2Input struct {
	Tab string `json:"tab" jsonschema:"history for filled/cancelled, pending for open orders (required)"`
}

func GetSpotOrdersV2(ctx context.Context, req *mcp.CallToolRequest, input GetSpotOrdersV2Input, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError("auth token required for this tool")
	}
	tab := input.Tab
	if tab != "history" && tab != "pending" {
		return toolError("tab must be history or pending")
	}
	q := url.Values{}
	q.Set("tab", tab)
	data, status, err := c.SpotGetV2Private("orders", q)
	return toolResultFromCall(data, status, err)
}

type GetOrderTradeHistoryInput struct {
	Market string `json:"market,omitempty" jsonschema:"Filter by market (optional)"`
	Limit  int    `json:"limit,omitempty" jsonschema:"Limit (optional)"`
}

func GetOrderTradeHistory(ctx context.Context, req *mcp.CallToolRequest, input GetOrderTradeHistoryInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError("admin service credentials must be set for this tool")
	}
	body := map[string]any{"timestamp": time.Now().UnixMilli()}
	if input.Market != "" {
		body["market"] = input.Market
	}
	if input.Limit > 0 {
		body["limit"] = input.Limit
	}
	data, status, err := c.PostSigned("/exchange/v1/orders/trade_history", body)
	return toolResultFromCall(data, status, err)
}

type GetOrderStatusMultipleInput struct {
	OrderIDs       []string `json:"order_ids,omitempty" jsonschema:"Array of order IDs"`
	ClientOrderIDs []string `json:"client_order_ids,omitempty" jsonschema:"Array of client order IDs"`
}

func GetOrderStatusMultiple(ctx context.Context, req *mcp.CallToolRequest, input GetOrderStatusMultipleInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError("admin service credentials must be set for this tool")
	}
	if len(input.OrderIDs) == 0 && len(input.ClientOrderIDs) == 0 {
		return toolError("order_ids or client_order_ids is required")
	}
	body := map[string]any{"timestamp": time.Now().UnixMilli()}
	if len(input.OrderIDs) > 0 {
		body["ids"] = input.OrderIDs
	}
	if len(input.ClientOrderIDs) > 0 {
		body["client_order_ids"] = input.ClientOrderIDs
	}
	data, status, err := c.PostSigned("/exchange/v1/orders/status_multiple", body)
	return toolResultFromCall(data, status, err)
}

type GetActiveOrdersCountInput struct {
	Market string `json:"market,omitempty" jsonschema:"Filter by market (optional)"`
}

func GetActiveOrdersCount(ctx context.Context, req *mcp.CallToolRequest, input GetActiveOrdersCountInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError("admin service credentials must be set for this tool")
	}
	body := map[string]any{"timestamp": time.Now().UnixMilli()}
	if input.Market != "" {
		body["market"] = input.Market
	}
	data, status, err := c.PostSigned("/exchange/v1/orders/active_orders_count", body)
	return toolResultFromCall(data, status, err)
}
