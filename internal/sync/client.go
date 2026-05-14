package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ClientConfig holds the configuration for constructing a Sync API client.
type ClientConfig struct {
	// Token is the Todoist API authentication token.
	Token string
	// Timeout is the HTTP client timeout. Defaults to 30 seconds if zero.
	Timeout time.Duration
	// Endpoint is the Sync API URL. Defaults to SyncEndpoint if empty.
	Endpoint string
}

// Client is an HTTP client for the Todoist Sync API v9.
type Client struct {
	httpClient *http.Client
	token      string
	endpoint   string
}

// NewClient creates a new Sync API client from the given configuration.
// If Timeout is zero, it defaults to 30 seconds.
// If Endpoint is empty, it defaults to SyncEndpoint.
func NewClient(cfg ClientConfig) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = SyncEndpoint
	}

	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		token:      cfg.Token,
		endpoint:   endpoint,
	}
}

// DoSync sends a sync request to the Todoist API and returns the response.
// The provided context is used for request cancellation and timeouts.
// Request parameters are sent as URL query parameters.
// A 401 response yields ErrAuthFailed; other non-2xx responses and
// network errors yield ErrSyncFailed wrapping the underlying cause.
func (c *Client) DoSync(ctx context.Context, req SyncRequest) (*SyncResponse, error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing endpoint: %v", ErrSyncFailed, err)
	}

	params := url.Values{}
	params.Set("sync_token", req.SyncToken)

	resourceTypesJSON, err := json.Marshal(req.ResourceTypes)
	if err != nil {
		return nil, fmt.Errorf("%w: encoding resource_types: %v", ErrSyncFailed, err)
	}
	params.Set("resource_types", string(resourceTypesJSON))

	if len(req.Commands) > 0 {
		commandsJSON, err := json.Marshal(req.Commands)
		if err != nil {
			return nil, fmt.Errorf("%w: encoding commands: %v", ErrSyncFailed, err)
		}
		params.Set("commands", string(commandsJSON))
	}
	u.RawQuery = params.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: building request: %v", ErrSyncFailed, err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSyncFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrAuthFailed
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: status %d", ErrSyncFailed, resp.StatusCode)
	}

	var syncResp SyncResponse
	if err := json.NewDecoder(resp.Body).Decode(&syncResp); err != nil {
		return nil, fmt.Errorf("%w: decoding response: %v", ErrSyncFailed, err)
	}

	return &syncResp, nil
}
