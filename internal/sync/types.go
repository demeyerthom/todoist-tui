package sync

import (
	"errors"

	"github.com/demeyerthom/todoist-tui/internal/model"
)

// SyncEndpoint is the Todoist Sync API v9 endpoint.
const SyncEndpoint = "https://api.todoist.com/api/v1/sync"

// ResourceTypes lists the resource types requested in every sync call.
var ResourceTypes = []string{"items", "projects", "sections", "labels", "filters"}

// Sentinel errors for sync operations.
var (
	// ErrAuthFailed is returned when the API responds with 401 Unauthorized.
	ErrAuthFailed = errors.New("sync: authentication failed (401)")

	// ErrSyncGone is returned when the API responds with 410 Gone,
	// indicating the stored sync token is no longer valid and a full
	// sync is required.
	ErrSyncGone = errors.New("sync: sync token expired (410)")

	// ErrSyncFailed is returned for non-200 responses or network errors.
	ErrSyncFailed = errors.New("sync: request failed")
)

// SyncRequest is the payload sent to the Todoist Sync API v9 endpoint.
type SyncRequest struct {
	SyncToken     string    `json:"sync_token"`
	ResourceTypes []string  `json:"resource_types"`
	Commands      []Command `json:"commands,omitempty"`
}

// SyncResponse is the payload received from the Todoist Sync API v9 endpoint.
type SyncResponse struct {
	SyncToken     string            `json:"sync_token"`
	Items         []model.Task      `json:"items"`
	Projects      []model.Project   `json:"projects"`
	Sections      []model.Section   `json:"sections"`
	Labels        []model.Label     `json:"labels"`
	Filters       []model.Filter    `json:"filters"`
	TempIDMapping map[string]string `json:"temp_id_mapping"`
}

// Command represents a single command sent to the Sync API for optimistic updates.
type Command struct {
	Type   string         `json:"type"`
	Args   map[string]any `json:"args"`
	UUID   string         `json:"uuid"`
	TempID string         `json:"temp_id,omitempty"`
}
