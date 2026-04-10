package tools

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type OptionsBaseCurrenciesInput struct {
	IsActive *bool `json:"isActive,omitempty" jsonschema:"Filter by active (optional)"`
}

func OptionsGetBaseCurrencies(ctx context.Context, req *mcp.CallToolRequest, input OptionsBaseCurrenciesInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	q := url.Values{}
	if input.IsActive != nil {
		q.Set("isActive", strconv.FormatBool(*input.IsActive))
	}
	data, status, err := c.OptionsGetPublic("base-currencies", q)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsInstrumentsInput struct {
	BaseCurrency string `json:"baseCurrency,omitempty" jsonschema:"Base currency (optional)"`
	Symbol       string `json:"symbol,omitempty" jsonschema:"Symbol (optional)"`
}

func OptionsGetInstruments(ctx context.Context, req *mcp.CallToolRequest, input OptionsInstrumentsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	q := url.Values{}
	if input.BaseCurrency != "" {
		q.Set("baseCurrency", input.BaseCurrency)
	}
	if input.Symbol != "" {
		if !ValidateOptionsSymbol(input.Symbol) {
			return toolError("invalid options symbol: expected format " + OptionsSymbolFormat)
		}
		q.Set("symbol", input.Symbol)
	}
	data, status, err := c.OptionsGetPublic("instruments", q)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsTickerInput struct {
	BaseCurrency string `json:"baseCurrency,omitempty" jsonschema:"Base currency"`
	ExpiryTime   string `json:"expiryTime,omitempty" jsonschema:"Expiry time"`
	Symbol       string `json:"symbol,omitempty" jsonschema:"Option symbol"`
}

func OptionsGetTicker(ctx context.Context, req *mcp.CallToolRequest, input OptionsTickerInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	hasBaseAndExpiry := input.BaseCurrency != "" && input.ExpiryTime != ""
	hasSymbol := input.Symbol != ""
	if !hasBaseAndExpiry && !hasSymbol {
		return toolError("use either (baseCurrency and expiryTime) or symbol")
	}
	if hasSymbol && !ValidateOptionsSymbol(input.Symbol) {
		return toolError("invalid options symbol: expected format " + OptionsSymbolFormat)
	}
	q := url.Values{}
	if input.BaseCurrency != "" {
		q.Set("baseCurrency", input.BaseCurrency)
	}
	if input.ExpiryTime != "" {
		q.Set("expiryTime", input.ExpiryTime)
	}
	if input.Symbol != "" {
		q.Set("symbol", input.Symbol)
	}
	data, status, err := c.OptionsGetPublic("ticker", q)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsOrderbookInput struct {
	Symbol string `json:"symbol" jsonschema:"Option symbol"`
	Depth  int    `json:"depth,omitempty" jsonschema:"Order book depth (optional)"`
}

func OptionsGetOrderbook(ctx context.Context, req *mcp.CallToolRequest, input OptionsOrderbookInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if input.Symbol == "" {
		return toolError("symbol is required")
	}
	if !ValidateOptionsSymbol(input.Symbol) {
		return toolError("invalid options symbol: expected format " + OptionsSymbolFormat)
	}
	q := url.Values{}
	depth := input.Depth
	if depth <= 0 || depth > 25 {
		depth = 25
	}
	q.Set("depth", strconv.Itoa(depth))
	data, status, err := c.OptionsGetPublic("orderbook/"+input.Symbol, q)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsTradesInput struct {
	Symbol string `json:"symbol" jsonschema:"Option symbol"`
	Limit  int    `json:"limit,omitempty" jsonschema:"Number of trades (optional)"`
}

func OptionsGetTrades(ctx context.Context, req *mcp.CallToolRequest, input OptionsTradesInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if input.Symbol == "" {
		return toolError("symbol is required")
	}
	if !ValidateOptionsSymbol(input.Symbol) {
		return toolError("invalid options symbol: expected format " + OptionsSymbolFormat)
	}
	q := url.Values{}
	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}
	q.Set("limit", strconv.Itoa(limit))
	data, status, err := c.OptionsGetPublic("trades/"+input.Symbol, q)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsStatsInput struct {
	BaseCurrencySymbol string `json:"base_currency_symbol" jsonschema:"Base currency symbol (e.g. BTC)"`
}

func OptionsGetStats(ctx context.Context, req *mcp.CallToolRequest, input OptionsStatsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if input.BaseCurrencySymbol == "" {
		return toolError("base_currency_symbol is required")
	}
	data, status, err := c.OptionsGetPublic("stats/"+input.BaseCurrencySymbol, nil)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsGetCombinedInstrumentsInput struct{}

func OptionsGetCombinedInstruments(ctx context.Context, req *mcp.CallToolRequest, _ OptionsGetCombinedInstrumentsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	data, status, err := c.OptionsGetPublic("instrument-details/combined", nil)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}

type OptionsGetIndexPriceInput struct {
	BaseCurrency string `json:"baseCurrency" jsonschema:"Base currency (e.g. BTC, ETH)"`
}

func OptionsGetIndexPrice(ctx context.Context, req *mcp.CallToolRequest, input OptionsGetIndexPriceInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if input.BaseCurrency == "" {
		return toolError("baseCurrency is required")
	}
	data, status, err := c.OptionsGetPublic("price/"+input.BaseCurrency, nil)
	if err != nil {
		return toolError("request failed: " + err.Error())
	}
	return toolResult(data, status != 200)
}
