---
# todoist-tui-c4lv
title: HTTP client with auth and sync endpoint
status: todo
type: task
created_at: 2026-05-03T15:09:41Z
updated_at: 2026-05-03T15:09:41Z
parent: todoist-tui-efa6
blocked_by:
    - todoist-tui-7gco
---

## Description

Implement the Sync API HTTP client with Bearer auth, context support, and the core DoSync method.

## Requirements

- ClientConfig: Token string, Timeout time.Duration, Endpoint string (with default)
- NewClient(cfg ClientConfig) *Client creates http.Client with timeout and auth header
- DoSync(ctx context.Context, req SyncRequest) (*SyncResponse, error) sends POST to endpoint
- Auth: 'Authorization: Bearer <token>' header
- Request body: JSON-encoded SyncRequest
- Response: JSON-decoded SyncResponse
- Context-aware: respect ctx cancellation
- Timeout via http.Client.Timeout, default 30s
- Parse API error responses and return meaningful errors
- Sentinel errors: ErrAuthFailed, ErrSyncFailed

## Acceptance Criteria

- Client sends authenticated POST requests
- Context cancellation is respected
- Timeouts and network errors return wrapped errors
- 'go build ./...' and 'go vet ./...' pass
