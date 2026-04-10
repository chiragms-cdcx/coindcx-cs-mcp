package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type KBSearchKnownIssuesInput struct {
	Keyword  string `json:"keyword" jsonschema:"Search keyword or symptom description"`
	Category string `json:"category,omitempty" jsonschema:"Category filter (e.g. futures, spot, deposits, withdrawals)"`
}

func KBSearchKnownIssues(ctx context.Context, req *mcp.CallToolRequest, input KBSearchKnownIssuesInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.Keyword == "" {
		return toolError("keyword is required")
	}
	q := url.Values{}
	q.Set("keyword", input.Keyword)
	if input.Category != "" {
		q.Set("category", input.Category)
	}
	data, status, err := c.AdminGet("/admin/api/v1/kb/known-issues", q)
	return toolResultFromCall(data, status, err)
}

type KBGetSOPInput struct {
	IssueType string `json:"issue_type" jsonschema:"Issue type (e.g. gs_futures_stuck, chart_loading, deposit_not_reflecting)"`
}

func KBGetSOP(ctx context.Context, req *mcp.CallToolRequest, input KBGetSOPInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.IssueType == "" {
		return toolError("issue_type is required")
	}
	q := url.Values{}
	q.Set("issue_type", input.IssueType)
	data, status, err := c.AdminGet("/admin/api/v1/kb/sop", q)
	return toolResultFromCall(data, status, err)
}

type KBGetStandardResponseInput struct {
	IssueType string `json:"issue_type" jsonschema:"Issue type for the response template"`
	Language  string `json:"language,omitempty" jsonschema:"Response language (default: en)"`
}

func KBGetStandardResponse(ctx context.Context, req *mcp.CallToolRequest, input KBGetStandardResponseInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.IssueType == "" {
		return toolError("issue_type is required")
	}
	q := url.Values{}
	q.Set("issue_type", input.IssueType)
	if input.Language != "" {
		q.Set("language", input.Language)
	}
	data, status, err := c.AdminGet("/admin/api/v1/kb/standard-response", q)
	return toolResultFromCall(data, status, err)
}

type KBSearchPastCRTsInput struct {
	Keyword  string `json:"keyword" jsonschema:"Search keyword, symptom, or error description"`
	Category string `json:"category,omitempty" jsonschema:"CRT category label filter"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Max results (default 10)"`
}

func KBSearchPastCRTs(ctx context.Context, req *mcp.CallToolRequest, input KBSearchPastCRTsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.Keyword == "" {
		return toolError("keyword is required")
	}
	q := url.Values{}
	q.Set("keyword", input.Keyword)
	if input.Category != "" {
		q.Set("category", input.Category)
	}
	if input.Limit > 0 {
		q.Set("limit", fmtInt(input.Limit))
	}
	data, status, err := c.AdminGet("/admin/api/v1/kb/past-crts", q)
	return toolResultFromCall(data, status, err)
}

type KBGetCRTDetailsInput struct {
	CRTID string `json:"crt_id" jsonschema:"CRT ticket ID (e.g. CRT-7455)"`
}

func KBGetCRTDetails(ctx context.Context, req *mcp.CallToolRequest, input KBGetCRTDetailsInput, c CoinDCXClient) (*mcp.CallToolResult, any, error) {
	if !c.HasCredentials() {
		return toolError(adminCredsRequired)
	}
	if input.CRTID == "" {
		return toolError("crt_id is required")
	}
	q := url.Values{}
	q.Set("crt_id", input.CRTID)
	data, status, err := c.AdminGet("/admin/api/v1/kb/crt-details", q)
	return toolResultFromCall(data, status, err)
}
