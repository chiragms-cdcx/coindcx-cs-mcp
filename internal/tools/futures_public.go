package tools

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type FuturesActiveInstrumentsInput struct {
	MarginCurrency string `json:"margin_currency_short_name,omitempty" jsonschema:"USDT or INR (default USDT)"`
}

func FuturesActiveInstruments(ctx context.Context, req *mcp.CallToolRequest, input FuturesActiveInstrumentsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	q := url.Values{}
	margin := input.MarginCurrency
	if margin == "" {
		margin = "USDT"
	}
	q.Set("margin_currency_short_name[]", margin)
	data, status, err := c.GetBase("/exchange/v1/derivatives/futures/data/active_instruments", q)
	return toolResultFromCall(data, status, err)
}

type FuturesInstrumentInput struct {
	Pair                    string `json:"pair" jsonschema:"Instrument pair e.g. KC-BTC_USDT"`
	MarginCurrencyShortName string `json:"margin_currency_short_name,omitempty" jsonschema:"USDT or INR (default USDT)"`
}

func FuturesInstrument(ctx context.Context, req *mcp.CallToolRequest, input FuturesInstrumentInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if input.Pair == "" {
		return toolError("pair is required")
	}
	if !ValidateFuturesPair(input.Pair) {
		return toolError("pair must use format B-* or KC-*")
	}
	q := url.Values{}
	q.Set("pair", input.Pair)
	margin := input.MarginCurrencyShortName
	if margin == "" {
		margin = "USDT"
	}
	q.Set("margin_currency_short_name", margin)
	data, status, err := c.GetBase("/exchange/v1/derivatives/futures/data/instrument", q)
	return toolResultFromCall(data, status, err)
}

type FuturesTradesInput struct {
	Pair string `json:"pair" jsonschema:"Instrument pair e.g. KC-BTC_USDT"`
}

func FuturesTrades(ctx context.Context, req *mcp.CallToolRequest, input FuturesTradesInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if input.Pair == "" {
		return toolError("pair is required")
	}
	if !ValidateFuturesPair(input.Pair) {
		return toolError("pair must use format B-* or KC-*")
	}
	q := url.Values{}
	q.Set("pair", input.Pair)
	data, status, err := c.GetBase("/exchange/v1/derivatives/futures/data/trades", q)
	return toolResultFromCall(data, status, err)
}

type FuturesCurrentPricesInput struct{}

func FuturesCurrentPrices(ctx context.Context, req *mcp.CallToolRequest, _ FuturesCurrentPricesInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	data, status, err := c.GetPublic("/market_data/v3/current_prices/futures/rt", nil)
	return toolResultFromCall(data, status, err)
}

type FuturesOrderbookInput struct {
	Pair  string `json:"pair" jsonschema:"Instrument pair e.g. KC-BTC_USDT"`
	Depth int    `json:"depth,omitempty" jsonschema:"Depth (default 50)"`
}

func FuturesOrderbook(ctx context.Context, req *mcp.CallToolRequest, input FuturesOrderbookInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if input.Pair == "" {
		return toolError("pair is required")
	}
	if !ValidateFuturesPair(input.Pair) {
		return toolError("pair must use format B-* or KC-*")
	}
	depth := input.Depth
	if depth <= 0 {
		depth = 50
	}
	path := "/market_data/v3/orderbook/" + input.Pair + "-futures/" + strconv.Itoa(depth)
	data, status, err := c.GetPublic(path, nil)
	return toolResultFromCall(data, status, err)
}
