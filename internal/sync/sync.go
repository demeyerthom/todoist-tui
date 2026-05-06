package sync

import (
	"context"
	"fmt"

	"github.com/demeyerthom/todoist-tui/internal/store"
)

// FullSync performs a full sync against the Todoist Sync API v9.
// It sends sync_token="*" with all resource types, writes every entity
// collection into the bbolt store, and persists the returned sync token.
//
// Auth errors (401) yield ErrAuthFailed; other non-2xx responses and
// network errors yield ErrSyncFailed. Store errors are returned as-is.
func (c *Client) FullSync(ctx context.Context, s *store.Store) error {
	req := SyncRequest{
		SyncToken:     "*",
		ResourceTypes: ResourceTypes,
	}

	resp, err := c.DoSync(ctx, req)
	if err != nil {
		return err
	}

	for i := range resp.Items {
		if err := s.PutTask(&resp.Items[i]); err != nil {
			return fmt.Errorf("sync: store task %s: %w", resp.Items[i].ID, err)
		}
	}

	for i := range resp.Projects {
		if err := s.PutProject(&resp.Projects[i]); err != nil {
			return fmt.Errorf("sync: store project %s: %w", resp.Projects[i].ID, err)
		}
	}

	for i := range resp.Sections {
		if err := s.PutSection(&resp.Sections[i]); err != nil {
			return fmt.Errorf("sync: store section %s: %w", resp.Sections[i].ID, err)
		}
	}

	for i := range resp.Labels {
		if err := s.PutLabel(&resp.Labels[i]); err != nil {
			return fmt.Errorf("sync: store label %s: %w", resp.Labels[i].ID, err)
		}
	}

	for i := range resp.Filters {
		if err := s.PutFilter(&resp.Filters[i]); err != nil {
			return fmt.Errorf("sync: store filter %s: %w", resp.Filters[i].ID, err)
		}
	}

	if err := s.SetSyncToken(resp.SyncToken); err != nil {
		return fmt.Errorf("sync: persist sync token: %w", err)
	}

	return nil
}

// IncrementalSync performs an incremental sync using the stored sync token.
// If no sync token has been stored (empty string), it delegates to FullSync.
// Otherwise it sends the stored token with all resource types, merges updated
// entities into the store, removes deleted entities, and persists the new
// sync token.
func (c *Client) IncrementalSync(ctx context.Context, s *store.Store) error {
	token, err := s.GetSyncToken()
	if err != nil {
		return fmt.Errorf("sync: read sync token: %w", err)
	}

	if token == "" {
		return c.FullSync(ctx, s)
	}

	req := SyncRequest{
		SyncToken:     token,
		ResourceTypes: ResourceTypes,
	}

	resp, err := c.DoSync(ctx, req)
	if err != nil {
		return err
	}

	for i := range resp.Items {
		item := &resp.Items[i]
		if item.IsDeleted {
			if err := s.DeleteTask(item.ID); err != nil {
				return fmt.Errorf("sync: delete task %s: %w", item.ID, err)
			}
		} else {
			if err := s.PutTask(item); err != nil {
				return fmt.Errorf("sync: store task %s: %w", item.ID, err)
			}
		}
	}

	for i := range resp.Projects {
		p := &resp.Projects[i]
		if p.IsDeleted {
			if err := s.DeleteProject(p.ID); err != nil {
				return fmt.Errorf("sync: delete project %s: %w", p.ID, err)
			}
		} else {
			if err := s.PutProject(p); err != nil {
				return fmt.Errorf("sync: store project %s: %w", p.ID, err)
			}
		}
	}

	for i := range resp.Sections {
		sec := &resp.Sections[i]
		if sec.IsDeleted {
			if err := s.DeleteSection(sec.ID); err != nil {
				return fmt.Errorf("sync: delete section %s: %w", sec.ID, err)
			}
		} else {
			if err := s.PutSection(sec); err != nil {
				return fmt.Errorf("sync: store section %s: %w", sec.ID, err)
			}
		}
	}

	for i := range resp.Labels {
		l := &resp.Labels[i]
		if l.IsDeleted {
			if err := s.DeleteLabel(l.ID); err != nil {
				return fmt.Errorf("sync: delete label %s: %w", l.ID, err)
			}
		} else {
			if err := s.PutLabel(l); err != nil {
				return fmt.Errorf("sync: store label %s: %w", l.ID, err)
			}
		}
	}

	for i := range resp.Filters {
		f := &resp.Filters[i]
		if f.IsDeleted {
			if err := s.DeleteFilter(f.ID); err != nil {
				return fmt.Errorf("sync: delete filter %s: %w", f.ID, err)
			}
		} else {
			if err := s.PutFilter(f); err != nil {
				return fmt.Errorf("sync: store filter %s: %w", f.ID, err)
			}
		}
	}

	if err := s.SetSyncToken(resp.SyncToken); err != nil {
		return fmt.Errorf("sync: persist sync token: %w", err)
	}

	return nil
}