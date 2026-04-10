package prompts

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func All() []struct {
	Prompt  *mcp.Prompt
	Handler mcp.PromptHandler
} {
	return []struct {
		Prompt  *mcp.Prompt
		Handler mcp.PromptHandler
	}{
		{Prompt: &mcp.Prompt{Name: "cs_investigate_crt", Description: "Investigate a CRT ticket: look up user, check data, find root cause, suggest response.", Arguments: []*mcp.PromptArgument{{Name: "crt_id", Description: "CRT ticket ID (e.g. CRT-7455)", Required: true}, {Name: "user_id", Description: "User ID or email of the customer", Required: false}}}, Handler: csInvestigateCRTHandler},
		{Prompt: &mcp.Prompt{Name: "cs_user_overview", Description: "Get a complete overview of a customer: profile, KYC, balances, recent activity.", Arguments: []*mcp.PromptArgument{{Name: "user_identifier", Description: "Email, phone, user_id, or CoinDCX UID", Required: true}}}, Handler: csUserOverviewHandler},
		{Prompt: &mcp.Prompt{Name: "cs_order_investigation", Description: "Investigate an order issue for a customer (stuck, cancelled, not filled, etc.).", Arguments: []*mcp.PromptArgument{{Name: "user_id", Description: "User ID", Required: true}, {Name: "order_id", Description: "Order ID if known", Required: false}, {Name: "product", Description: "spot, futures, or options", Required: false}}}, Handler: csOrderInvestigationHandler},
		{Prompt: &mcp.Prompt{Name: "cs_deposit_investigation", Description: "Investigate a deposit not reflecting for a customer.", Arguments: []*mcp.PromptArgument{{Name: "user_id", Description: "User ID", Required: true}, {Name: "currency", Description: "Currency (e.g. BTC, USDT, INR)", Required: false}}}, Handler: csDepositInvestigationHandler},
		{Prompt: &mcp.Prompt{Name: "cs_withdrawal_investigation", Description: "Investigate a stuck or failed withdrawal.", Arguments: []*mcp.PromptArgument{{Name: "user_id", Description: "User ID", Required: true}, {Name: "currency", Description: "Currency", Required: false}}}, Handler: csWithdrawalInvestigationHandler},
		{Prompt: &mcp.Prompt{Name: "cs_futures_investigation", Description: "Investigate futures issues (stuck orders, liquidation, negative balance).", Arguments: []*mcp.PromptArgument{{Name: "user_id", Description: "User ID", Required: true}, {Name: "pair", Description: "Futures pair (e.g. KC-BTC_USDT)", Required: false}}}, Handler: csFuturesInvestigationHandler},
		{Prompt: &mcp.Prompt{Name: "cs_options_investigation", Description: "Investigate options issues (position not reflecting, liquidation).", Arguments: []*mcp.PromptArgument{{Name: "user_id", Description: "User ID", Required: true}, {Name: "symbol", Description: "Options symbol", Required: false}}}, Handler: csOptionsInvestigationHandler},
		{Prompt: &mcp.Prompt{Name: "cs_known_issue_check", Description: "Check if a customer's issue matches a known issue with an existing SOP.", Arguments: []*mcp.PromptArgument{{Name: "symptom", Description: "Customer's reported symptom or issue description", Required: true}}}, Handler: csKnownIssueCheckHandler},
		{Prompt: &mcp.Prompt{Name: "cs_compliance_check", Description: "Check compliance status, flags, and restrictions for a customer.", Arguments: []*mcp.PromptArgument{{Name: "user_id", Description: "User ID", Required: true}}}, Handler: csComplianceCheckHandler},
		{Prompt: &mcp.Prompt{Name: "cs_generate_response", Description: "Generate a customer-facing response based on investigation findings.", Arguments: []*mcp.PromptArgument{{Name: "issue_type", Description: "Type of issue (e.g. deposit_not_reflecting, order_not_filled)", Required: true}, {Name: "findings", Description: "Summary of what you found during investigation", Required: true}}}, Handler: csGenerateResponseHandler},
	}
}

func getArg(args map[string]string, key, defaultVal string) string {
	if args == nil {
		return defaultVal
	}
	if v, ok := args[key]; ok && v != "" {
		return v
	}
	return defaultVal
}

func csInvestigateCRTHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	crtID := getArg(req.Params.Arguments, "crt_id", "")
	userID := getArg(req.Params.Arguments, "user_id", "")
	text := fmt.Sprintf(`Investigate CRT ticket %s.

1. First, fetch the CRT details: call kb_get_crt_details with crt_id "%s".
2. Check if this matches a known issue: call kb_search_known_issues with keywords from the CRT description.
3. If a user_id is available (%s), look up the customer:
   - Call admin_lookup_user (if needed) then admin_get_user_profile to get account status.
   - Based on the CRT category, fetch relevant data (orders, deposits, withdrawals, positions).
4. Check for past similar CRTs: call kb_search_past_crts with the symptom.
5. Summarize findings: what happened, root cause, and whether this needs engineering or can be resolved by CS.
6. If resolvable, get the standard response: call kb_get_standard_response.
7. Present a clear summary with the suggested customer response.`, crtID, crtID, userID)
	return &mcp.GetPromptResult{
		Description: "Investigate a CRT ticket end-to-end.",
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}, nil
}

func csUserOverviewHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	identifier := getArg(req.Params.Arguments, "user_identifier", "")
	text := fmt.Sprintf(`Get a complete overview for customer: %s

1. Look up the user: call admin_lookup_user with the identifier (email/phone/user_id/uid).
2. Use the returned user_id to fetch in parallel via coindcx_request:
   - admin_get_user_profile (KYC, account tier, verification status)
   - admin_get_user_balances (all wallets)
   - admin_get_user_compliance (flags, restrictions)
3. Summarize concisely:
   - Account status, KYC status, account tier
   - Non-zero balances across spot, futures, options
   - Any compliance flags or restrictions
   - Recent login info if relevant`, identifier)
	return &mcp.GetPromptResult{
		Description: "Complete customer overview for CS agents.",
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}, nil
}

func csOrderInvestigationHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	userID := getArg(req.Params.Arguments, "user_id", "")
	orderID := getArg(req.Params.Arguments, "order_id", "")
	product := getArg(req.Params.Arguments, "product", "spot")
	text := fmt.Sprintf(`Investigate order issue for user %s (product: %s, order: %s).

1. Fetch the user's recent orders: call admin_get_user_orders with user_id and product filter.
2. If order_id is provided, check its specific status.
3. For spot: check market conditions (coindcx_get_orderbook, coindcx_get_ticker).
4. For futures: check positions and market data.
5. Check known issues: call kb_search_known_issues with relevant keywords.
6. Explain why the order behaved as it did (not filled due to price, cancelled, etc.).
7. Suggest resolution or standard response.`, userID, product, orderID)
	return &mcp.GetPromptResult{
		Description: "Investigate order issues.",
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}, nil
}

func csDepositInvestigationHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	userID := getArg(req.Params.Arguments, "user_id", "")
	currency := getArg(req.Params.Arguments, "currency", "")
	text := fmt.Sprintf(`Investigate deposit not reflecting for user %s (currency: %s).

1. Call admin_get_user_deposits with user_id and currency filter.
2. Call admin_get_user_profile to check KYC and compliance status.
3. Check known issues: call kb_search_known_issues with "deposit not reflecting".
4. Look at deposit status: pending, confirmed, failed, etc.
5. Summarize: deposit status, any known issues, suggested next steps.
6. Get standard response if available.`, userID, currency)
	return &mcp.GetPromptResult{
		Description: "Investigate deposit issues.",
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}, nil
}

func csWithdrawalInvestigationHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	userID := getArg(req.Params.Arguments, "user_id", "")
	currency := getArg(req.Params.Arguments, "currency", "")
	text := fmt.Sprintf(`Investigate withdrawal issue for user %s (currency: %s).

1. Call admin_get_user_withdrawals with user_id and currency filter.
2. Check compliance: call admin_get_user_compliance.
3. Check known issues: call kb_search_known_issues with "withdrawal stuck" or "withdrawal failed".
4. Identify the withdrawal status and any failure reasons.
5. Summarize findings and suggest resolution.`, userID, currency)
	return &mcp.GetPromptResult{
		Description: "Investigate withdrawal issues.",
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}, nil
}

func csFuturesInvestigationHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	userID := getArg(req.Params.Arguments, "user_id", "")
	pair := getArg(req.Params.Arguments, "pair", "")
	text := fmt.Sprintf(`Investigate futures issue for user %s (pair: %s).

1. Call admin_get_user_futures_positions with user_id.
2. Call admin_get_user_orders with user_id and product=futures.
3. If pair is provided, check market data: coindcx_futures_stats, coindcx_futures_current_prices.
4. Check known issues: call kb_search_known_issues with "futures" keyword.
5. Look for common patterns: GS stuck orders, liquidation, negative balance.
6. Summarize findings and suggest resolution.`, userID, pair)
	return &mcp.GetPromptResult{
		Description: "Investigate futures issues.",
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}, nil
}

func csOptionsInvestigationHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	userID := getArg(req.Params.Arguments, "user_id", "")
	symbol := getArg(req.Params.Arguments, "symbol", "")
	text := fmt.Sprintf(`Investigate options issue for user %s (symbol: %s).

1. Call admin_get_user_options_positions with user_id.
2. Check wallet: coindcx_options_get_wallet_balance.
3. If symbol is provided, check market data: coindcx_options_get_ticker.
4. Check known issues: call kb_search_known_issues with "options" keyword.
5. Summarize findings and suggest resolution.`, userID, symbol)
	return &mcp.GetPromptResult{
		Description: "Investigate options issues.",
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}, nil
}

func csKnownIssueCheckHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	symptom := getArg(req.Params.Arguments, "symptom", "")
	text := fmt.Sprintf(`Check if this customer symptom matches a known issue: "%s"

1. Call kb_search_known_issues with the symptom as keyword.
2. If matches found, call kb_get_sop for the matching issue type.
3. Call kb_get_standard_response for the issue type.
4. Call kb_search_past_crts to find similar resolved tickets.
5. Summarize: is this a known issue? What is the SOP? What is the standard response?
6. If not a known issue, suggest escalation path.`, symptom)
	return &mcp.GetPromptResult{
		Description: "Check if symptom matches a known issue.",
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}, nil
}

func csComplianceCheckHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	userID := getArg(req.Params.Arguments, "user_id", "")
	text := fmt.Sprintf(`Check compliance status for user %s.

1. Call admin_get_user_compliance to check flags, restrictions, AML status.
2. Call admin_get_user_kyc_details for KYC verification status.
3. Call admin_get_user_profile for account tier and verification level.
4. Summarize: compliance status, any restrictions, KYC status, and what the CS agent should communicate to the customer.`, userID)
	return &mcp.GetPromptResult{
		Description: "Check compliance and KYC status.",
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}, nil
}

func csGenerateResponseHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	issueType := getArg(req.Params.Arguments, "issue_type", "")
	findings := getArg(req.Params.Arguments, "findings", "")
	text := fmt.Sprintf(`Generate a customer-facing response for issue type: %s

Findings from investigation: %s

1. Call kb_get_standard_response with the issue_type to get a template.
2. Customize the template with the specific findings.
3. The response should be:
   - Professional and empathetic
   - Clear about what happened and why
   - Specific about next steps (if any)
   - Avoid technical jargon the customer won't understand
4. Present the response ready for the CS agent to copy and send.`, issueType, findings)
	return &mcp.GetPromptResult{
		Description: "Generate customer-facing response.",
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}, nil
}
