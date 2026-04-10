package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const adminCredsRequired = "Admin service credentials required: ensure X-ADMIN-SERVICE-KEY and X-ADMIN-SERVICE-SECRET headers are set"

type AdminLookupUserInput struct {
	Email  string `json:"email,omitempty" jsonschema:"Search by email"`
	Phone  string `json:"phone,omitempty" jsonschema:"Search by phone number"`
	UserID string `json:"user_id,omitempty" jsonschema:"Search by internal user ID"`
	UID    string `json:"uid,omitempty" jsonschema:"Search by CoinDCX UID"`
}

func AdminLookupUser(ctx context.Context, req *mcp.CallToolRequest, input AdminLookupUserInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.Email == "" && input.Phone == "" && input.UserID == "" && input.UID == "" {
		return toolError("at least one of email, phone, user_id, or uid is required")
	}
	q := url.Values{}
	if input.Email != "" {
		q.Set("email", input.Email)
	}
	if input.Phone != "" {
		q.Set("phone", input.Phone)
	}
	if input.UserID != "" {
		q.Set("user_id", input.UserID)
	}
	if input.UID != "" {
		q.Set("uid", input.UID)
	}
	data, status, err := c.AdminGet("/admin/api/v1/users/lookup", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserProfileInput struct {
	UserID string `json:"user_id" jsonschema:"User ID to look up"`
}

func AdminGetUserProfile(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserProfileInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	data, status, err := c.AdminGet("/admin/api/v1/users/profile", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserBalancesInput struct {
	UserID string `json:"user_id" jsonschema:"User ID"`
}

func AdminGetUserBalances(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserBalancesInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	data, status, err := c.AdminGet("/admin/api/v1/users/balances", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserOrdersInput struct {
	UserID  string `json:"user_id" jsonschema:"User ID"`
	Market  string `json:"market,omitempty" jsonschema:"Filter by market"`
	Status  string `json:"status,omitempty" jsonschema:"Filter by status (open, filled, cancelled)"`
	Product string `json:"product,omitempty" jsonschema:"spot, futures, or options"`
}

func AdminGetUserOrders(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserOrdersInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	if input.Market != "" {
		q.Set("market", input.Market)
	}
	if input.Status != "" {
		q.Set("status", input.Status)
	}
	if input.Product != "" {
		q.Set("product", input.Product)
	}
	data, status, err := c.AdminGet("/admin/api/v1/users/orders", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserDepositsInput struct {
	UserID   string `json:"user_id" jsonschema:"User ID"`
	Currency string `json:"currency,omitempty" jsonschema:"Filter by currency (e.g. BTC, INR)"`
	Status   string `json:"status,omitempty" jsonschema:"Filter by status"`
}

func AdminGetUserDeposits(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserDepositsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	if input.Currency != "" {
		q.Set("currency", input.Currency)
	}
	if input.Status != "" {
		q.Set("status", input.Status)
	}
	data, status, err := c.AdminGet("/admin/api/v1/users/deposits", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserWithdrawalsInput struct {
	UserID   string `json:"user_id" jsonschema:"User ID"`
	Currency string `json:"currency,omitempty" jsonschema:"Filter by currency"`
	Status   string `json:"status,omitempty" jsonschema:"Filter by status"`
}

func AdminGetUserWithdrawals(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserWithdrawalsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	if input.Currency != "" {
		q.Set("currency", input.Currency)
	}
	if input.Status != "" {
		q.Set("status", input.Status)
	}
	data, status, err := c.AdminGet("/admin/api/v1/users/withdrawals", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserTransactionsInput struct {
	UserID string `json:"user_id" jsonschema:"User ID"`
	Limit  int    `json:"limit,omitempty" jsonschema:"Max results (default 50)"`
}

func AdminGetUserTransactions(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserTransactionsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	if input.Limit > 0 {
		q.Set("limit", fmtInt(input.Limit))
	}
	data, status, err := c.AdminGet("/admin/api/v1/users/transactions", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserKYCInput struct {
	UserID string `json:"user_id" jsonschema:"User ID"`
}

func AdminGetUserKYCDetails(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserKYCInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	data, status, err := c.AdminGet("/admin/api/v1/users/kyc", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserComplianceInput struct {
	UserID string `json:"user_id" jsonschema:"User ID"`
}

func AdminGetUserCompliance(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserComplianceInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	data, status, err := c.AdminGet("/admin/api/v1/users/compliance", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserFuturesPositionsInput struct {
	UserID string `json:"user_id" jsonschema:"User ID"`
	Pair   string `json:"pair,omitempty" jsonschema:"Filter by pair (e.g. KC-BTC_USDT)"`
}

func AdminGetUserFuturesPositions(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserFuturesPositionsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	if input.Pair != "" {
		q.Set("pair", input.Pair)
	}
	data, status, err := c.AdminGet("/admin/api/v1/users/futures/positions", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserOptionsPositionsInput struct {
	UserID string `json:"user_id" jsonschema:"User ID"`
	Symbol string `json:"symbol,omitempty" jsonschema:"Filter by symbol"`
}

func AdminGetUserOptionsPositions(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserOptionsPositionsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	if input.Symbol != "" {
		q.Set("symbol", input.Symbol)
	}
	data, status, err := c.AdminGet("/admin/api/v1/users/options/positions", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserLoginHistoryInput struct {
	UserID string `json:"user_id" jsonschema:"User ID"`
	Limit  int    `json:"limit,omitempty" jsonschema:"Max results (default 20)"`
}

func AdminGetUserLoginHistory(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserLoginHistoryInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	if input.Limit > 0 {
		q.Set("limit", fmtInt(input.Limit))
	}
	data, status, err := c.AdminGet("/admin/api/v1/users/login-history", q)
	return toolResultFromCall(data, status, err)
}

type AdminGetUserReferralsInput struct {
	UserID string `json:"user_id" jsonschema:"User ID"`
}

func AdminGetUserReferrals(ctx context.Context, req *mcp.CallToolRequest, input AdminGetUserReferralsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.UserID == "" {
		return toolError("user_id is required")
	}
	q := url.Values{}
	q.Set("user_id", input.UserID)
	data, status, err := c.AdminGet("/admin/api/v1/users/referrals", q)
	return toolResultFromCall(data, status, err)
}
