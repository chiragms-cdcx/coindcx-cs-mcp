package tools

import (
	"context"
	"encoding/json"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/logger"
	"github.com/chiragms-cdcx/coindcx-cs-mcp/pkg/audit"
)

const OptionsSymbolFormat = "BASE-DDMMMYY-STRIKE-C|P-USDT (e.g. BTC-14MAR26-73500-C-USDT)"

var optionsSymbolRe = regexp.MustCompile(`^[A-Z0-9]+-\d{2}(?:JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)\d{2}-\d+-[CP]-USDT$`)

func ValidateOptionsSymbol(s string) bool {
	if s == "" {
		return false
	}
	return optionsSymbolRe.MatchString(strings.ToUpper(s))
}

func parseFloat64(s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f, err == nil
}

func ValidateFuturesPair(pair string) bool {
	if pair == "" || len(pair) < 4 {
		return false
	}
	return strings.HasPrefix(pair, "B-") || strings.HasPrefix(pair, "KC-")
}

type BatchHandler func(ctx context.Context, args json.RawMessage, c CoinDCXClient) (*mcp.CallToolResult, any, error)

type BatchHandlerFactory func(c CoinDCXClient) BatchHandler

var (
	batchRegistry   = make(map[string]BatchHandlerFactory)
	batchRegistryMu sync.RWMutex
)

func RegisterBatchFactory(name string, factory BatchHandlerFactory) {
	batchRegistryMu.Lock()
	defer batchRegistryMu.Unlock()
	if _, exists := batchRegistry[name]; !exists {
		batchRegistry[name] = factory
	}
}

func GetBatchHandler(name string, c CoinDCXClient) BatchHandler {
	batchRegistryMu.RLock()
	defer batchRegistryMu.RUnlock()
	f := batchRegistry[name]
	if f == nil {
		return nil
	}
	return f(c)
}

func MakeBatchHandlerFactory[T any](
	fn func(context.Context, *mcp.CallToolRequest, T, CoinDCXClient) (*mcp.CallToolResult, any, error),
) BatchHandlerFactory {
	return func(c CoinDCXClient) BatchHandler {
		return func(ctx context.Context, args json.RawMessage, client CoinDCXClient) (*mcp.CallToolResult, any, error) {
			var input T
			if len(args) > 0 {
				if err := json.Unmarshal(args, &input); err != nil {
					res, _, _ := toolError("invalid arguments: " + err.Error())
					return res, nil, nil
				}
			}
			req := &mcp.CallToolRequest{Params: &mcp.CallToolParamsRaw{Name: "", Arguments: args}}
			return fn(ctx, req, input, client)
		}
	}
}

type CoinDCXClient interface {
	GetPublic(path string, query url.Values) ([]byte, int, error)
	GetBase(path string, query url.Values) ([]byte, int, error)
	PostSigned(path string, body any) ([]byte, int, error)
	GetSigned(path string, query url.Values, body any) ([]byte, int, error)
	HasCredentials() bool
	OptionsGetPublic(path string, query url.Values) ([]byte, int, error)
	OptionsGetPrivate(path string, query url.Values) ([]byte, int, error)
	OptionsPostPrivate(path string, body any) ([]byte, int, error)
	FuturesGetPrivate(path string, query url.Values) ([]byte, int, error)
	FuturesPostPrivate(path string, body any) ([]byte, int, error)
	SpotGetV2Private(path string, query url.Values) ([]byte, int, error)
	HasAuthToken() bool
	AdminGet(path string, query url.Values) ([]byte, int, error)
	AdminPost(path string, body any) ([]byte, int, error)
}

func toolResult(data []byte, isError bool) (*mcp.CallToolResult, any, error) {
	text := string(data)
	if text == "" {
		text = "{}"
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
		IsError: isError,
	}, nil, nil
}

func toolResultFromStruct(v any, isError bool) (*mcp.CallToolResult, any, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
			IsError: true,
		}, nil, nil
	}
	return toolResult(data, isError)
}

func toolError(msg string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}, nil, nil
}

func toolResultFromCall(data []byte, status int, err error) (*mcp.CallToolResult, any, error) {
	if err != nil {
		logger.Error("API request failed: %v", err)
		return toolError("request failed: " + err.Error())
	}
	if status == 401 {
		msg := "API returned 401 Unauthorized. "
		if len(data) > 0 {
			msg += "Response: " + string(data) + ". "
		}
		msg += "Check admin service credentials (X-ADMIN-SERVICE-KEY/SECRET) or auth token (X-COINDCX-AUTH-TOKEN)."
		return toolError(msg)
	}
	if status >= 400 {
		logger.Error("API returned status=%d", status)
	}
	return toolResult(data, status != 200)
}

func Wrap[T any](
	fn func(context.Context, *mcp.CallToolRequest, T, CoinDCXClient) (*mcp.CallToolResult, any, error),
	c CoinDCXClient,
) func(context.Context, *mcp.CallToolRequest, T) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input T) (*mcp.CallToolResult, any, error) {
		toolName := req.Params.Name
		logger.Info("tool call: %s", toolName)
		logger.Debug("tool %s arguments:\n%s", toolName, logger.PrettyJSON(string(req.Params.Arguments)))
		res, data, err := fn(ctx, req, input, c)
		if err != nil {
			logger.Error("tool %s failed: %v", toolName, err)
		} else if res != nil && res.IsError {
			logger.Error("tool %s returned error", toolName)
		} else {
			logger.Debug("tool %s completed successfully", toolName)
		}
		return res, data, err
	}
}

func firstTextContent(res *mcp.CallToolResult) string {
	if res == nil || len(res.Content) == 0 {
		return ""
	}
	if tc, ok := res.Content[0].(*mcp.TextContent); ok {
		return strings.TrimSpace(tc.Text)
	}
	return ""
}

type BatchCall struct {
	Name      string         `json:"name" jsonschema:"Tool name"`
	Arguments map[string]any `json:"arguments,omitempty" jsonschema:"JSON object of arguments for the tool"`
}

type BatchInput struct {
	Calls []BatchCall `json:"calls" jsonschema:"List of tool calls to run concurrently"`
}

type BatchResultItem struct {
	Name    string `json:"name"`
	IsError bool   `json:"is_error"`
	Content string `json:"content"`
}

func Batch(ctx context.Context, req *mcp.CallToolRequest, input BatchInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if len(input.Calls) == 0 {
		res, _, _ := toolError("calls must be a non-empty array")
		return res, nil, nil
	}
	if len(input.Calls) > 20 {
		res, _, _ := toolError("at most 20 concurrent calls per batch")
		return res, nil, nil
	}

	type result struct {
		index int
		name  string
		res   *mcp.CallToolResult
	}
	results := make([]result, len(input.Calls))
	var wg sync.WaitGroup
	for i, call := range input.Calls {
		if call.Name == "" {
			res, _, _ := toolError("each call must have a non-empty name")
			return res, nil, nil
		}
		handler := GetBatchHandler(call.Name, c)
		if handler == nil {
			res, _, _ := toolError("unknown tool: " + call.Name)
			return res, nil, nil
		}
		var argsRaw json.RawMessage
		if len(call.Arguments) > 0 {
			argsRaw, _ = json.Marshal(call.Arguments)
		}
		wg.Add(1)
		go func(idx int, name string, args json.RawMessage) {
			defer wg.Done()
			res, _, _ := handler(ctx, args, c)
			results[idx] = result{index: idx, name: name, res: res}
		}(i, call.Name, argsRaw)
	}
	wg.Wait()
	items := make([]BatchResultItem, len(results))
	for i, r := range results {
		content := ""
		if r.res != nil && len(r.res.Content) > 0 {
			if tc, ok := r.res.Content[0].(*mcp.TextContent); ok {
				content = tc.Text
			}
		}
		items[i] = BatchResultItem{Name: r.name, IsError: r.res != nil && r.res.IsError, Content: content}
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		res, _, _ := toolError("failed to marshal batch result: " + err.Error())
		return res, nil, nil
	}
	res, _, _ := toolResult(data, false)
	return res, nil, nil
}

func WrapWithAudit[T any](
	fn func(context.Context, *mcp.CallToolRequest, T, CoinDCXClient) (*mcp.CallToolResult, any, error),
	c CoinDCXClient,
	pub audit.Publisher,
	transport audit.Transport,
) func(context.Context, *mcp.CallToolRequest, T) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input T) (*mcp.CallToolResult, any, error) {
		start := time.Now()
		result, extra, err := fn(ctx, req, input, c)
		latencyMs := time.Since(start).Milliseconds()

		var status audit.Status
		var errorMessage string
		switch {
		case err != nil:
			status = audit.StatusError
			errorMessage = audit.TruncateError(err.Error())
		case result != nil && result.IsError:
			status = audit.StatusError
			if len(result.Content) > 0 {
				if tc, ok := result.Content[0].(*mcp.TextContent); ok {
					errorMessage = audit.TruncateError(tc.Text)
				}
			}
		default:
			status = audit.StatusSuccess
		}

		requestID := uuid.New().String()
		if token := req.Params.GetProgressToken(); token != nil {
			if s, ok := token.(string); ok && s != "" {
				requestID = s
			}
		}

		event := audit.Event{
			ID:           uuid.New().String(),
			Timestamp:    time.Now().UTC(),
			ToolName:     req.Params.Name,
			RequestID:    requestID,
			Status:       status,
			ErrorMessage: errorMessage,
			Transport:    transport,
			LatencyMs:    latencyMs,
			Params:       req.Params.Arguments,
		}
		pub.Publish(ctx, event)
		return result, extra, err
	}
}
