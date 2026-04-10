package tools

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func spotPairToExchangeSymbol(pair string) string {
	if pair == "" {
		return ""
	}
	if idx := strings.Index(pair, "-"); idx >= 0 {
		pair = pair[idx+1:]
	}
	return strings.ReplaceAll(pair, "_", "")
}

func spotPairToMarketDataPair(pair string) string {
	if pair == "" || len(pair) < 6 {
		return ""
	}
	if strings.HasSuffix(pair, "INR") {
		base := strings.TrimSuffix(pair, "INR")
		return "I-" + base + "_INR"
	}
	return ""
}

type GetTickerInput struct{}

func GetTicker(ctx context.Context, req *mcp.CallToolRequest, _ GetTickerInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	data, status, err := c.GetPublic("/exchange/ticker", nil)
	return toolResultFromCall(data, status, err)
}

type GetMarketsInput struct{}

func GetMarkets(ctx context.Context, req *mcp.CallToolRequest, _ GetMarketsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	data, status, err := c.GetPublic("/exchange/v1/markets", nil)
	return toolResultFromCall(data, status, err)
}

type GetMarketsDetailsInput struct{}

func GetMarketsDetails(ctx context.Context, req *mcp.CallToolRequest, _ GetMarketsDetailsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	data, status, err := c.GetPublic("/exchange/v1/markets_details", nil)
	return toolResultFromCall(data, status, err)
}

type GetOrderbookInput struct {
	Pair string `json:"pair" jsonschema:"Market pair (e.g. KC-BTC_USDT)"`
}

func GetOrderbook(ctx context.Context, req *mcp.CallToolRequest, input GetOrderbookInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if input.Pair == "" {
		return toolError("pair is required")
	}
	q := url.Values{}
	q.Set("pair", input.Pair)
	data, status, err := c.GetPublic("/market_data/orderbook", q)
	if err != nil {
		return toolResultFromCall(data, status, err)
	}
	if status == 422 {
		if alt := spotPairToMarketDataPair(input.Pair); alt != "" {
			q.Set("pair", alt)
			data, status, err = c.GetPublic("/market_data/orderbook", q)
		}
		if status == 422 {
			symbol := spotPairToExchangeSymbol(input.Pair)
			if symbol != "" && symbol != input.Pair {
				q.Set("pair", symbol)
				data, status, err = c.GetBase("/exchange/orderbook", q)
			}
		}
	}
	return toolResultFromCall(data, status, err)
}

type GetCandlesInput struct {
	Pair      string `json:"pair" jsonschema:"Market pair (e.g. KC-BTC_USDT)"`
	Interval  string `json:"interval" jsonschema:"One of: 1m, 15m, 1h, 1d"`
	Limit     int    `json:"limit,omitempty" jsonschema:"Default 500, max 1000"`
	StartTime int64  `json:"startTime,omitempty" jsonschema:"Start time in ms"`
	EndTime   int64  `json:"endTime,omitempty" jsonschema:"End time in ms"`
}

var candlesIntervalAllowed = map[string]bool{"1m": true, "15m": true, "1h": true, "1d": true}
var candlesIntervalMap = map[string]string{
	"5m": "15m", "30m": "15m", "2h": "1h", "4h": "1h", "6h": "1h", "8h": "1h", "3d": "1d", "1w": "1d", "1M": "1d",
}

func candlesIntervalForAPI(interval string) string {
	interval = strings.TrimSpace(strings.ToLower(interval))
	if candlesIntervalAllowed[interval] {
		return interval
	}
	if mapped, ok := candlesIntervalMap[interval]; ok {
		return mapped
	}
	return interval
}

func GetCandles(ctx context.Context, req *mcp.CallToolRequest, input GetCandlesInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if input.Pair == "" || input.Interval == "" {
		return toolError("pair and interval are required")
	}
	q := url.Values{}
	q.Set("pair", input.Pair)
	q.Set("interval", candlesIntervalForAPI(input.Interval))
	if input.Limit > 0 {
		q.Set("limit", fmtInt(input.Limit))
	}
	if input.StartTime > 0 {
		q.Set("startTime", fmtInt64(input.StartTime))
	}
	if input.EndTime > 0 {
		q.Set("endTime", fmtInt64(input.EndTime))
	}
	data, status, err := c.GetPublic("/market_data/candles", q)
	if err != nil {
		return toolResultFromCall(data, status, err)
	}
	if status == 422 {
		if alt := spotPairToMarketDataPair(input.Pair); alt != "" {
			q.Set("pair", alt)
			data, status, err = c.GetPublic("/market_data/candles", q)
		}
	}
	return toolResultFromCall(data, status, err)
}

func fmtInt(i int) string     { return strconv.Itoa(i) }
func fmtInt64(i int64) string { return strconv.FormatInt(i, 10) }
