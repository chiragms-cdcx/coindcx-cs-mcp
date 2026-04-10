package tools

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/logger"
)

var cryptoNewsFeeds = []struct {
	Name string
	URL  string
}{
	{"CoinDesk", "https://www.coindesk.com/arc/outboundfeeds/rss/"},
	{"CoinTelegraph", "https://cointelegraph.com/rss"},
	{"Decrypt", "https://decrypt.co/feed"},
}

const (
	defaultNewsLimit = 10
	maxNewsLimit     = 25
	feedTimeout      = 12 * time.Second
)

type rssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Item []rssItem `xml:"item"`
	} `xml:"channel"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type CryptoNewsArticle struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Source      string `json:"source"`
	PublishedAt string `json:"published_at"`
	Snippet     string `json:"snippet,omitempty"`
}

type GetCryptoNewsInput struct {
	Limit    int    `json:"limit,omitempty" jsonschema:"Max articles (default 10, max 25)"`
	Keyword  string `json:"keyword,omitempty" jsonschema:"Filter by keyword"`
	DaysBack int    `json:"days_back,omitempty" jsonschema:"Articles from last N days (1-90)"`
}

func fetchCryptoNewsArticles(ctx context.Context, limit int, keyword string, daysBack int) []CryptoNewsArticle {
	if limit <= 0 {
		limit = defaultNewsLimit
	}
	if limit > maxNewsLimit {
		limit = maxNewsLimit
	}
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	client := &http.Client{Timeout: feedTimeout}
	type feedResult struct {
		source string
		items  []rssItem
	}
	results := make(chan feedResult, len(cryptoNewsFeeds))
	for _, feed := range cryptoNewsFeeds {
		go func(name, url string) {
			items := fetchRSS(ctx, client, url)
			results <- feedResult{source: name, items: items}
		}(feed.Name, feed.URL)
	}
	var cutoff time.Time
	if daysBack > 0 {
		if daysBack > 90 {
			daysBack = 90
		}
		cutoff = time.Now().AddDate(0, 0, -daysBack)
	}
	var all []CryptoNewsArticle
	seen := make(map[string]bool)
	for i := 0; i < len(cryptoNewsFeeds); i++ {
		r := <-results
		for _, it := range r.items {
			if daysBack > 0 {
				if t := parsePubDate(it.PubDate); !t.IsZero() && t.Before(cutoff) {
					continue
				}
			}
			title := stripHTML(it.Title)
			snippet := stripHTML(it.Description)
			if len(snippet) > 200 {
				snippet = snippet[:200] + "..."
			}
			if keyword != "" {
				lower := strings.ToLower(title + " " + snippet)
				if !strings.Contains(lower, keyword) {
					continue
				}
			}
			link := strings.TrimSpace(it.Link)
			if link == "" || seen[link] {
				continue
			}
			seen[link] = true
			all = append(all, CryptoNewsArticle{Title: title, Link: link, Source: r.source, PublishedAt: it.PubDate, Snippet: snippet})
			if len(all) >= limit {
				break
			}
		}
		if len(all) >= limit {
			break
		}
	}
	if len(all) > limit {
		all = all[:limit]
	}
	return all
}

func GetCryptoNews(ctx context.Context, req *mcp.CallToolRequest, input GetCryptoNewsInput, _ CoinDCXClient) (*mcp.CallToolResult, any, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = defaultNewsLimit
	}
	if limit > maxNewsLimit {
		limit = maxNewsLimit
	}
	all := fetchCryptoNewsArticles(ctx, limit, input.Keyword, input.DaysBack)
	out := map[string]any{"articles": all, "sources": []string{"CoinDesk", "CoinTelegraph", "Decrypt"}, "count": len(all)}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		logger.Error("marshal crypto news: %v", err)
		return toolError("failed to format news: " + err.Error())
	}
	return toolResult(data, false)
}

func fetchRSS(ctx context.Context, client *http.Client, url string) []rssItem {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "CoinDCX-CS-MCP/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	var feed rssFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil
	}
	return feed.Channel.Item
}

func parsePubDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC1123Z, time.RFC1123} {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func stripHTML(s string) string {
	s = html.UnescapeString(s)
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(b.String())
}

type GetCryptoNewsSummaryInput struct {
	Limit    int    `json:"limit,omitempty" jsonschema:"Max headlines (default 5, max 15)"`
	Keyword  string `json:"keyword,omitempty" jsonschema:"Filter keyword"`
	DaysBack int    `json:"days_back,omitempty" jsonschema:"Headlines from last N days (1-90)"`
}

func GetCryptoNewsSummary(ctx context.Context, req *mcp.CallToolRequest, input GetCryptoNewsSummaryInput, _ CoinDCXClient) (*mcp.CallToolResult, any, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 5
	}
	if limit > 15 {
		limit = 15
	}
	articles := fetchCryptoNewsArticles(ctx, limit, strings.TrimSpace(input.Keyword), input.DaysBack)
	headlines := make([]map[string]string, 0, len(articles))
	for _, a := range articles {
		headlines = append(headlines, map[string]string{"title": a.Title, "source": a.Source, "link": a.Link})
	}
	out := map[string]any{"headlines": headlines, "count": len(headlines)}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return toolError("failed to format summary: " + err.Error())
	}
	return toolResult(data, false)
}
