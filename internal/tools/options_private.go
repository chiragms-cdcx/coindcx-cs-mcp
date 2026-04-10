package tools

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const authTokenRequired = "Auth token required for options: set X-COINDCX-AUTH-TOKEN header"

type OptionsGetOrdersInput struct {
	BaseCurrency string `json:"baseCurrency,omitempty" jsonschema:"Filter by base currency"`
	Symbol       string `json:"symbol,omitempty" jsonschema:"Filter by symbol"`
	Size         int    `json:"size,omitempty" jsonschema:"Page size (optional)"`
}

func OptionsGetOrders(ctx context.Context, req *mcp.CallToolRequest, input OptionsGetOrdersInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(authTokenRequired)
	}
	q := url.Values{}
	if input.BaseCurrency != "" {
		q.Set("baseCurrency", input.BaseCurrency)
	}
	if input.Symbol != "" {
		q.Set("symbol", input.Symbol)
	}
	size := input.Size
	if size <= 0 {
		size = 10
	}
	q.Set("size", strconv.Itoa(size))
	data, status, err := c.OptionsGetPrivate("orders", q)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsOrderHistoryInput struct {
	Size         int    `json:"size,omitempty" jsonschema:"Page size"`
	Status       string `json:"status,omitempty" jsonschema:"Filter by status"`
	BaseCurrency string `json:"baseCurrency,omitempty" jsonschema:"Filter by base currency"`
	Cursor       string `json:"cursor,omitempty" jsonschema:"Pagination cursor"`
	Symbol       string `json:"symbol,omitempty" jsonschema:"Filter by symbol"`
}

func OptionsGetOrderHistory(ctx context.Context, req *mcp.CallToolRequest, input OptionsOrderHistoryInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(authTokenRequired)
	}
	q := url.Values{}
	size := input.Size
	if size <= 0 {
		size = 10
	}
	q.Set("size", strconv.Itoa(size))
	if input.Status != "" {
		q.Set("status", input.Status)
	}
	if input.BaseCurrency != "" {
		q.Set("baseCurrency", input.BaseCurrency)
	}
	if input.Cursor != "" {
		q.Set("cursor", input.Cursor)
	}
	if input.Symbol != "" {
		q.Set("symbol", input.Symbol)
	}
	data, status, err := c.OptionsGetPrivate("order-history", q)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsGetTradesByOrderIDInput struct {
	OrderID string `json:"orderId" jsonschema:"Order ID"`
}

func OptionsGetTradesByOrderID(ctx context.Context, req *mcp.CallToolRequest, input OptionsGetTradesByOrderIDInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(authTokenRequired)
	}
	if input.OrderID == "" {
		return toolError("orderId is required")
	}
	data, status, err := c.OptionsGetPrivate("trades/"+input.OrderID, nil)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsGetPositionsInput struct {
	Size int `json:"size,omitempty" jsonschema:"Page size"`
}

func OptionsGetPositions(ctx context.Context, req *mcp.CallToolRequest, input OptionsGetPositionsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(authTokenRequired)
	}
	q := url.Values{}
	size := input.Size
	if size <= 0 {
		size = 10
	}
	q.Set("size", strconv.Itoa(size))
	data, status, err := c.OptionsGetPrivate("positions", q)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsGetWalletBalanceInput struct{}

func OptionsGetWalletBalance(ctx context.Context, req *mcp.CallToolRequest, _ OptionsGetWalletBalanceInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(authTokenRequired)
	}
	data, status, err := c.OptionsGetPrivate("wallet/balance", nil)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsClosedPositionsInput struct {
	Symbol    string `json:"symbol,omitempty" jsonschema:"Filter by symbol"`
	StartTime int64  `json:"startTime,omitempty" jsonschema:"Start time Unix ms"`
	EndTime   int64  `json:"endTime,omitempty" jsonschema:"End time Unix ms"`
	Limit     int    `json:"limit,omitempty" jsonschema:"Limit"`
	Cursor    string `json:"cursor,omitempty" jsonschema:"Pagination cursor"`
}

func OptionsGetClosedPositions(ctx context.Context, req *mcp.CallToolRequest, input OptionsClosedPositionsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(authTokenRequired)
	}
	q := url.Values{}
	if input.Symbol != "" {
		q.Set("symbol", input.Symbol)
	}
	if input.StartTime > 0 {
		q.Set("startTime", strconv.FormatInt(input.StartTime, 10))
	}
	if input.EndTime > 0 {
		q.Set("endTime", strconv.FormatInt(input.EndTime, 10))
	}
	if input.Limit > 0 {
		q.Set("limit", strconv.Itoa(input.Limit))
	}
	if input.Cursor != "" {
		q.Set("cursor", input.Cursor)
	}
	data, status, err := c.OptionsGetPrivate("closed-positions", q)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsGetMarginModeInput struct{}

func OptionsGetMarginMode(ctx context.Context, req *mcp.CallToolRequest, _ OptionsGetMarginModeInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(authTokenRequired)
	}
	data, status, err := c.OptionsGetPrivate("margin-mode/current", nil)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsGetConsumedLimitInput struct{}

func OptionsGetWalletConsumedLimit(ctx context.Context, req *mcp.CallToolRequest, _ OptionsGetConsumedLimitInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(authTokenRequired)
	}
	data, status, err := c.OptionsGetPrivate("wallet/consumed-limit", nil)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsWalletTransactionsInput struct {
	Size   int    `json:"size,omitempty" jsonschema:"Page size (default 20)"`
	Cursor string `json:"cursor,omitempty" jsonschema:"Pagination cursor"`
}

func OptionsGetWalletTransactions(ctx context.Context, req *mcp.CallToolRequest, input OptionsWalletTransactionsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasAuthToken() {
		return toolError(authTokenRequired)
	}
	q := url.Values{}
	size := input.Size
	if size <= 0 {
		size = 20
	}
	q.Set("size", strconv.Itoa(size))
	if input.Cursor != "" {
		q.Set("cursor", input.Cursor)
	}
	data, status, err := c.OptionsGetPrivate("wallet/transactions", q)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}
