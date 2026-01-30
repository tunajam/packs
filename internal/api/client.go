package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	packsv1 "github.com/tunajam/packs/gen/packs/v1"
	"github.com/tunajam/packs/gen/packs/v1/packsv1connect"
)

const (
	DefaultBaseURL = "https://packs-api.fly.dev"
	EnvBaseURL     = "PACKS_API_URL"
)

// Client wraps the Packs API
type Client struct {
	client packsv1connect.PacksServiceClient
}

// New creates a new API client
func New() *Client {
	baseURL := os.Getenv(EnvBaseURL)
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &Client{
		client: packsv1connect.NewPacksServiceClient(httpClient, baseURL),
	}
}

// SearchOpts are options for searching packs
type SearchOpts struct {
	Query  string
	Type   string // skill, context, prompt
	Tags   []string
	Author string
	Limit  int32
	Offset int32
	Sort   string // relevance, stars, newest, name
}

// PackSummary represents a pack in search results
type PackSummary struct {
	Name        string
	Version     string
	Type        string
	Description string
	Author      string
	Stars       int32
	Tags        []string
}

// Pack represents a full pack with content
type Pack struct {
	PackSummary
	Content   string
	GithubRef string
}

// Search searches for packs
func (c *Client) Search(ctx context.Context, opts SearchOpts) ([]PackSummary, int32, error) {
	req := &packsv1.SearchRequest{
		Query:  opts.Query,
		Author: opts.Author,
		Limit:  opts.Limit,
		Offset: opts.Offset,
		Sort:   opts.Sort,
	}

	if opts.Type != "" {
		switch opts.Type {
		case "skill":
			req.Type = packsv1.PackType_PACK_TYPE_SKILL
		case "context":
			req.Type = packsv1.PackType_PACK_TYPE_CONTEXT
		case "prompt":
			req.Type = packsv1.PackType_PACK_TYPE_PROMPT
		}
	}

	if len(opts.Tags) > 0 {
		req.Tags = opts.Tags
	}

	resp, err := c.client.Search(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, 0, err
	}

	var packs []PackSummary
	for _, p := range resp.Msg.Packs {
		packs = append(packs, PackSummary{
			Name:        p.Name,
			Version:     p.Version,
			Type:        packTypeToString(p.Type),
			Description: p.Description,
			Author:      p.Author,
			Stars:       p.Stars,
			Tags:        p.Tags,
		})
	}

	return packs, resp.Msg.Total, nil
}

// Get fetches a pack by name and optional version
func (c *Client) Get(ctx context.Context, name, version string) (*Pack, error) {
	req := &packsv1.GetRequest{
		Name:    name,
		Version: version,
	}

	resp, err := c.client.Get(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	p := resp.Msg.Pack
	return &Pack{
		PackSummary: PackSummary{
			Name:        p.Name,
			Version:     p.Version,
			Type:        packTypeToString(p.Type),
			Description: p.Description,
			Author:      p.Author,
			Stars:       p.Stars,
			Tags:        p.Tags,
		},
		Content:   p.Content,
		GithubRef: p.GithubRef,
	}, nil
}

// Submit submits a GitHub pack for indexing
func (c *Client) Submit(ctx context.Context, githubRef string) (name, version, message string, err error) {
	req := &packsv1.SubmitRequest{
		GithubRef: githubRef,
	}

	resp, err := c.client.Submit(ctx, connect.NewRequest(req))
	if err != nil {
		return "", "", "", err
	}

	return resp.Msg.Name, resp.Msg.Version, resp.Msg.Message, nil
}

// Telemetry sends a telemetry event (fire and forget)
func (c *Client) Telemetry(ctx context.Context, pack, source, version, cliVersion, os, arch string) {
	req := &packsv1.TelemetryEvent{
		Pack:       pack,
		Source:     source,
		Version:    version,
		CliVersion: cliVersion,
		Os:         os,
		Arch:       arch,
	}

	// Fire and forget
	go c.client.Telemetry(ctx, connect.NewRequest(req))
}

func packTypeToString(t packsv1.PackType) string {
	switch t {
	case packsv1.PackType_PACK_TYPE_SKILL:
		return "skill"
	case packsv1.PackType_PACK_TYPE_CONTEXT:
		return "context"
	case packsv1.PackType_PACK_TYPE_PROMPT:
		return "prompt"
	default:
		return "unknown"
	}
}
