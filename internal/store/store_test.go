package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/demeyerthom/todoist-tui/internal/model"
)

// openTestStore creates a Store backed by a temporary database file.
// It returns the store and a cleanup function that closes the store.
func openTestStore(t *testing.T) (*Store, func()) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	s, err := New(StoreConfig{Path: path})
	if err != nil {
		t.Fatalf("open test store: %v", err)
	}
	return s, func() { _ = s.Close() }
}

func TestOpen(t *testing.T) {
	t.Run("creates_db_file_and_all_buckets", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.db")

		s, err := New(StoreConfig{Path: path})
		if err != nil {
			t.Fatalf("New(): %v", err)
		}
		t.Cleanup(func() { _ = s.Close() })

		// Verify the database file was created.
		if _, err := os.Stat(path); err != nil {
			t.Errorf("db file not created: %v", err)
		}

		// Verify all expected buckets exist by writing and reading from each.
		for _, name := range buckets {
			err := s.Put(name, "__probe__", "test")
			if err != nil {
				t.Errorf("bucket %q not accessible: %v", name, err)
			}
		}
	})

	t.Run("creates_parent_directories", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "nested", "dir", "test.db")

		s, err := New(StoreConfig{Path: path})
		if err != nil {
			t.Fatalf("New() with nested path: %v", err)
		}
		t.Cleanup(func() { _ = s.Close() })

		if _, err := os.Stat(path); err != nil {
			t.Errorf("db file not created at nested path: %v", err)
		}
	})
}

func TestClose(t *testing.T) {
	t.Run("is_idempotent", func(t *testing.T) {
		s, _ := openTestStore(t)

		if err := s.Close(); err != nil {
			t.Fatalf("first Close(): %v", err)
		}
		if err := s.Close(); !errors.Is(err, ErrStoreNotOpen) {
			t.Errorf("second Close(): got err=%v, want ErrStoreNotOpen", err)
		}
	})

	t.Run("returns_ErrStoreNotOpen_on_uninitialized", func(t *testing.T) {
		s := &Store{}
		if err := s.Close(); !errors.Is(err, ErrStoreNotOpen) {
			t.Errorf("Close() on nil db: got err=%v, want ErrStoreNotOpen", err)
		}
	})
}

func TestPutGet(t *testing.T) {
	s, cleanup := openTestStore(t)
	defer cleanup()

	t.Run("task_round_trip", func(t *testing.T) {
		want := &model.Task{
			ID:        "task-1",
			Content:   "Buy groceries",
			ProjectID: "proj-1",
			Priority:  3,
			Completed: false,
			Labels:    []string{"lbl-1"},
			Due:       &model.DueDate{Date: "2025-01-15", IsRecurring: true},
		}
		if err := s.PutTask(want); err != nil {
			t.Fatalf("PutTask: %v", err)
		}
		got, err := s.GetTask(want.ID)
		if err != nil {
			t.Fatalf("GetTask: %v", err)
		}
		if got.ID != want.ID || got.Content != want.Content || got.Priority != want.Priority {
			t.Errorf("round trip mismatch: got=%+v, want=%+v", got, want)
		}
		if len(got.Labels) != len(want.Labels) {
			t.Errorf("labels length: got=%d, want=%d", len(got.Labels), len(want.Labels))
		}
		if got.Due == nil || got.Due.Date != want.Due.Date {
			t.Errorf("due date: got=%v, want=%v", got.Due, want.Due)
		}
	})

	t.Run("project_round_trip", func(t *testing.T) {
		want := &model.Project{
			ID:         "proj-1",
			Name:       "Inbox",
			Color:      "grey",
			IsFavorite: false,
			IsInbox:    true,
		}
		if err := s.PutProject(want); err != nil {
			t.Fatalf("PutProject: %v", err)
		}
		got, err := s.GetProject(want.ID)
		if err != nil {
			t.Fatalf("GetProject: %v", err)
		}
		if got.ID != want.ID || got.Name != want.Name || got.IsInbox != want.IsInbox {
			t.Errorf("round trip mismatch: got=%+v, want=%+v", got, want)
		}
	})

	t.Run("section_round_trip", func(t *testing.T) {
		want := &model.Section{
			ID:           "sec-1",
			Name:         "Morning",
			ProjectID:    "proj-1",
			SectionOrder: 1,
		}
		if err := s.PutSection(want); err != nil {
			t.Fatalf("PutSection: %v", err)
		}
		got, err := s.GetSection(want.ID)
		if err != nil {
			t.Fatalf("GetSection: %v", err)
		}
		if got.ID != want.ID || got.Name != want.Name || got.SectionOrder != want.SectionOrder {
			t.Errorf("round trip mismatch: got=%+v, want=%+v", got, want)
		}
	})

	t.Run("label_round_trip", func(t *testing.T) {
		want := &model.Label{
			ID:         "lbl-1",
			Name:       "urgent",
			Color:      "red",
			IsFavorite: true,
			ItemOrder:  5,
		}
		if err := s.PutLabel(want); err != nil {
			t.Fatalf("PutLabel: %v", err)
		}
		got, err := s.GetLabel(want.ID)
		if err != nil {
			t.Fatalf("GetLabel: %v", err)
		}
		if got.ID != want.ID || got.Name != want.Name || got.ItemOrder != want.ItemOrder {
			t.Errorf("round trip mismatch: got=%+v, want=%+v", got, want)
		}
	})

	t.Run("filter_round_trip", func(t *testing.T) {
		want := &model.Filter{
			ID:    "flt-1",
			Name:  "Today",
			Query: "today",
			Color: "green",
		}
		if err := s.PutFilter(want); err != nil {
			t.Fatalf("PutFilter: %v", err)
		}
		got, err := s.GetFilter(want.ID)
		if err != nil {
			t.Fatalf("GetFilter: %v", err)
		}
		if got.ID != want.ID || got.Name != want.Name || got.Query != want.Query {
			t.Errorf("round trip mismatch: got=%+v, want=%+v", got, want)
		}
	})

	t.Run("sync_meta_round_trip", func(t *testing.T) {
		want := &model.SyncMeta{
			Key:   "some_key",
			Value: "some_value",
		}
		if err := s.PutSyncMeta(want); err != nil {
			t.Fatalf("PutSyncMeta: %v", err)
		}
		got, err := s.GetSyncMeta(want.Key)
		if err != nil {
			t.Fatalf("GetSyncMeta: %v", err)
		}
		if got.Key != want.Key || got.Value != want.Value {
			t.Errorf("round trip mismatch: got=%+v, want=%+v", got, want)
		}
	})

	t.Run("get_nonexistent_key_returns_error", func(t *testing.T) {
		_, err := s.GetTask("does-not-exist")
		if err == nil {
			t.Error("expected error for nonexistent key, got nil")
		}
	})
}

func TestList(t *testing.T) {
	s, cleanup := openTestStore(t)
	defer cleanup()

	t.Run("returns_all_items_in_bucket", func(t *testing.T) {
		tasks := []model.Task{
			{ID: "t1", Content: "First", Priority: 1},
			{ID: "t2", Content: "Second", Priority: 2},
			{ID: "t3", Content: "Third", Priority: 3},
		}
		for _, task := range tasks {
			if err := s.PutTask(&task); err != nil {
				t.Fatalf("PutTask(%s): %v", task.ID, err)
			}
		}

		got, err := s.ListTasks()
		if err != nil {
			t.Fatalf("ListTasks: %v", err)
		}
		if len(got) != len(tasks) {
			t.Fatalf("ListTasks: got %d items, want %d", len(got), len(tasks))
		}

		gotByID := make(map[string]model.Task)
		for _, t := range got {
			gotByID[t.ID] = t
		}
		for _, want := range tasks {
			g, ok := gotByID[want.ID]
			if !ok {
				t.Errorf("task %q not found in list results", want.ID)
				continue
			}
			if g.Content != want.Content || g.Priority != want.Priority {
				t.Errorf("task %q: got=%+v, want=%+v", want.ID, g, want)
			}
		}
	})

	t.Run("empty_bucket_returns_empty_slice", func(t *testing.T) {
		fresh, freshCleanup := openTestStore(t)
		defer freshCleanup()

		got, err := fresh.ListProjects()
		if err != nil {
			t.Fatalf("ListProjects on empty: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected 0 projects, got %d", len(got))
		}
	})
}

func TestDelete(t *testing.T) {
	s, cleanup := openTestStore(t)
	defer cleanup()

	t.Run("removes_item_from_bucket", func(t *testing.T) {
		task := &model.Task{ID: "del-1", Content: "Delete me"}
		if err := s.PutTask(task); err != nil {
			t.Fatalf("PutTask: %v", err)
		}

		// Confirm it exists.
		got, err := s.GetTask(task.ID)
		if err != nil {
			t.Fatalf("GetTask before delete: %v", err)
		}
		if got.ID != task.ID {
			t.Fatalf("task ID mismatch before delete: got=%s, want=%s", got.ID, task.ID)
		}

		// Delete it.
		if err := s.DeleteTask(task.ID); err != nil {
			t.Fatalf("DeleteTask: %v", err)
		}

		// Confirm it's gone.
		_, err = s.GetTask(task.ID)
		if err == nil {
			t.Error("expected error after deleting task, got nil")
		}
	})

	t.Run("delete_nonexistent_key_no_error", func(t *testing.T) {
		// bbolt Delete on a missing key is a no-op.
		if err := s.DeleteTask("ghost"); err != nil {
			t.Errorf("DeleteTask on nonexistent key: %v", err)
		}
	})
}

func TestSyncToken(t *testing.T) {
	t.Run("returns_empty_string_when_not_set", func(t *testing.T) {
		fresh, cleanup := openTestStore(t)
		defer cleanup()

		token, err := fresh.GetSyncToken()
		if err != nil {
			t.Fatalf("GetSyncToken on empty: %v", err)
		}
		if token != "" {
			t.Errorf("expected empty token, got %q", token)
		}
	})

	t.Run("round_trip", func(t *testing.T) {
		s, cleanup := openTestStore(t)
		defer cleanup()

		want := "abc123def456"
		if err := s.SetSyncToken(want); err != nil {
			t.Fatalf("SetSyncToken: %v", err)
		}
		got, err := s.GetSyncToken()
		if err != nil {
			t.Fatalf("GetSyncToken: %v", err)
		}
		if got != want {
			t.Errorf("sync token: got=%q, want=%q", got, want)
		}
	})

	t.Run("overwrites_previous_value", func(t *testing.T) {
		s, cleanup := openTestStore(t)
		defer cleanup()

		if err := s.SetSyncToken("first"); err != nil {
			t.Fatalf("SetSyncToken first: %v", err)
		}
		if err := s.SetSyncToken("second"); err != nil {
			t.Fatalf("SetSyncToken second: %v", err)
		}
		got, err := s.GetSyncToken()
		if err != nil {
			t.Fatalf("GetSyncToken: %v", err)
		}
		if got != "second" {
			t.Errorf("expected overwritten token %q, got %q", "second", got)
		}
	})
}

func TestLastSyncTime(t *testing.T) {
	t.Run("returns_empty_string_when_not_set", func(t *testing.T) {
		fresh, cleanup := openTestStore(t)
		defer cleanup()

		ts, err := fresh.GetLastSyncTime()
		if err != nil {
			t.Fatalf("GetLastSyncTime on empty: %v", err)
		}
		if ts != "" {
			t.Errorf("expected empty timestamp, got %q", ts)
		}
	})

	t.Run("round_trip", func(t *testing.T) {
		s, cleanup := openTestStore(t)
		defer cleanup()

		want := "2025-01-15T10:30:00Z"
		tm, err := time.Parse(time.RFC3339, want)
		if err != nil {
			t.Fatalf("parse time: %v", err)
		}
		if err := s.SetLastSyncTime(tm); err != nil {
			t.Fatalf("SetLastSyncTime: %v", err)
		}
		got, err := s.GetLastSyncTime()
		if err != nil {
			t.Fatalf("GetLastSyncTime: %v", err)
		}
		if got != want {
			t.Errorf("last sync time: got=%q, want=%q", got, want)
		}
	})
}