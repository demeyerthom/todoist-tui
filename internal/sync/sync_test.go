package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/demeyerthom/todoist-tui/internal/model"
	"github.com/demeyerthom/todoist-tui/internal/store"
)

// openTestStore creates a Store backed by a temporary database file.
// It returns the store and a cleanup function.
func openTestStore(t *testing.T) (*store.Store, func()) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	s, err := store.New(store.StoreConfig{Path: path})
	if err != nil {
		t.Fatalf("open test store: %v", err)
	}
	return s, func() { _ = s.Close() }
}

// newTestClient creates a Client pointing at the given server URL with a test token.
func newTestClient(serverURL string) *Client {
	return NewClient(ClientConfig{
		Token:    "test-token",
		Endpoint: serverURL,
	})
}

// mockHandler returns an http.HandlerFunc that validates the incoming request
// and responds with the given SyncResponse JSON.
func mockHandler(t *testing.T, resp SyncResponse) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		auth := r.Header.Get("Authorization")
		if auth == "" || auth[:7] != "Bearer " {
			t.Errorf("missing or malformed Authorization header: %q", auth)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// Decode request body to verify it's well-formed.
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		var req SyncRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Errorf("unmarshal request: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}
}

// mockHandlerWithRequestCapture returns a handler that also stores the
// decoded SyncRequest in reqOut for inspection.
func mockHandlerWithRequestCapture(t *testing.T, resp SyncResponse, reqOut *SyncRequest) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		auth := r.Header.Get("Authorization")
		if auth == "" || auth[:7] != "Bearer " {
			t.Errorf("missing or malformed Authorization header: %q", auth)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		var req SyncRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Errorf("unmarshal request: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		*reqOut = req

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}
}

func TestFullSync(t *testing.T) {
	resp := SyncResponse{
		SyncToken: "full-sync-token-1",
		Items: []model.Task{
			{ID: "item-1", Content: "Buy groceries", ProjectID: "proj-1", Priority: 3},
			{ID: "item-2", Content: "Write report", ProjectID: "proj-2", Priority: 4},
		},
		Projects: []model.Project{
			{ID: "proj-1", Name: "Personal", Color: "red"},
			{ID: "proj-2", Name: "Work", Color: "blue"},
		},
		Sections: []model.Section{
			{ID: "sec-1", Name: "Morning", ProjectID: "proj-1", SectionOrder: 1},
		},
		Labels: []model.Label{
			{ID: "lbl-1", Name: "urgent", Color: "red", ItemOrder: 1},
		},
		Filters: []model.Filter{
			{ID: "flt-1", Name: "Today", Query: "today", Color: "green"},
		},
		TempIDMapping: map[string]string{},
	}

	srv := httptest.NewServer(mockHandler(t, resp))
	defer srv.Close()

	s, cleanup := openTestStore(t)
	defer cleanup()

	client := newTestClient(srv.URL)
	if err := client.FullSync(context.Background(), s); err != nil {
		t.Fatalf("FullSync: %v", err)
	}

	// Verify tasks.
	tasks, err := s.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
	taskByID := make(map[string]model.Task)
	for _, tk := range tasks {
		taskByID[tk.ID] = tk
	}
	if tk, ok := taskByID["item-1"]; !ok || tk.Content != "Buy groceries" {
		t.Errorf("task item-1 not found or content mismatch: %+v", tk)
	}
	if tk, ok := taskByID["item-2"]; !ok || tk.Content != "Write report" {
		t.Errorf("task item-2 not found or content mismatch: %+v", tk)
	}

	// Verify projects.
	projects, err := s.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects: %v", err)
	}
	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projects))
	}
	projByID := make(map[string]model.Project)
	for _, p := range projects {
		projByID[p.ID] = p
	}
	if p, ok := projByID["proj-1"]; !ok || p.Name != "Personal" {
		t.Errorf("project proj-1 not found or name mismatch: %+v", p)
	}

	// Verify sections.
	sections, err := s.ListSections()
	if err != nil {
		t.Fatalf("ListSections: %v", err)
	}
	if len(sections) != 1 {
		t.Errorf("expected 1 section, got %d", len(sections))
	}
	if len(sections) > 0 && sections[0].Name != "Morning" {
		t.Errorf("section name: got=%q, want=%q", sections[0].Name, "Morning")
	}

	// Verify labels.
	labels, err := s.ListLabels()
	if err != nil {
		t.Fatalf("ListLabels: %v", err)
	}
	if len(labels) != 1 {
		t.Errorf("expected 1 label, got %d", len(labels))
	}
	if len(labels) > 0 && labels[0].Name != "urgent" {
		t.Errorf("label name: got=%q, want=%q", labels[0].Name, "urgent")
	}

	// Verify filters.
	filters, err := s.ListFilters()
	if err != nil {
		t.Fatalf("ListFilters: %v", err)
	}
	if len(filters) != 1 {
		t.Errorf("expected 1 filter, got %d", len(filters))
	}
	if len(filters) > 0 && filters[0].Name != "Today" {
		t.Errorf("filter name: got=%q, want=%q", filters[0].Name, "Today")
	}

	// Verify sync token.
	token, err := s.GetSyncToken()
	if err != nil {
		t.Fatalf("GetSyncToken: %v", err)
	}
	if token != "full-sync-token-1" {
		t.Errorf("sync token: got=%q, want=%q", token, "full-sync-token-1")
	}
}

func TestIncrementalSync(t *testing.T) {
	// Phase 1: FullSync to populate the store.
	fullResp := SyncResponse{
		SyncToken: "token-v1",
		Items: []model.Task{
			{ID: "item-1", Content: "Buy groceries", ProjectID: "proj-1", Priority: 3},
		},
		Projects: []model.Project{
			{ID: "proj-1", Name: "Personal", Color: "red"},
		},
		Sections:    []model.Section{},
		Labels:     []model.Label{},
		Filters:    []model.Filter{},
		TempIDMapping: map[string]string{},
	}

	srv := httptest.NewServer(mockHandler(t, fullResp))
	defer srv.Close()

	s, cleanup := openTestStore(t)
	defer cleanup()

	client := newTestClient(srv.URL)
	if err := client.FullSync(context.Background(), s); err != nil {
		t.Fatalf("FullSync: %v", err)
	}

	// Verify initial state.
	tasks, _ := s.ListTasks()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task after full sync, got %d", len(tasks))
	}

	// Phase 2: IncrementalSync with updated entities.
	// Close the first server and create a new one with updated data.
	srv.Close()

	incrResp := SyncResponse{
		SyncToken: "token-v2",
		Items: []model.Task{
			{ID: "item-1", Content: "Buy groceries and cook", ProjectID: "proj-1", Priority: 4},
			{ID: "item-3", Content: "New task", ProjectID: "proj-1", Priority: 1},
		},
		Projects: []model.Project{
			{ID: "proj-1", Name: "Personal (updated)", Color: "blue"},
		},
		Sections:    []model.Section{},
		Labels:     []model.Label{},
		Filters:    []model.Filter{},
		TempIDMapping: map[string]string{},
	}

	srv2 := httptest.NewServer(mockHandler(t, incrResp))
	defer srv2.Close()

	client2 := newTestClient(srv2.URL)
	if err := client2.IncrementalSync(context.Background(), s); err != nil {
		t.Fatalf("IncrementalSync: %v", err)
	}

	// Verify updated task.
	task, err := s.GetTask("item-1")
	if err != nil {
		t.Fatalf("GetTask item-1: %v", err)
	}
	if task.Content != "Buy groceries and cook" {
		t.Errorf("task content: got=%q, want=%q", task.Content, "Buy groceries and cook")
	}
	if task.Priority != 4 {
		t.Errorf("task priority: got=%d, want=%d", task.Priority, 4)
	}

	// Verify new task.
	task3, err := s.GetTask("item-3")
	if err != nil {
		t.Fatalf("GetTask item-3: %v", err)
	}
	if task3.Content != "New task" {
		t.Errorf("task3 content: got=%q, want=%q", task3.Content, "New task")
	}

	// Verify updated project.
	proj, err := s.GetProject("proj-1")
	if err != nil {
		t.Fatalf("GetProject proj-1: %v", err)
	}
	if proj.Name != "Personal (updated)" {
		t.Errorf("project name: got=%q, want=%q", proj.Name, "Personal (updated)")
	}

	// Verify sync token updated.
	token, err := s.GetSyncToken()
	if err != nil {
		t.Fatalf("GetSyncToken: %v", err)
	}
	if token != "token-v2" {
		t.Errorf("sync token: got=%q, want=%q", token, "token-v2")
	}
}

func TestIncrementalSync_EmptyToken(t *testing.T) {
	// When the store has no sync token, IncrementalSync should delegate to FullSync
	// which sends sync_token="*".
	var capturedReq SyncRequest
	resp := SyncResponse{
		SyncToken:    "fresh-token",
		Items:        []model.Task{},
		Projects:     []model.Project{},
		Sections:     []model.Section{},
		Labels:       []model.Label{},
		Filters:      []model.Filter{},
		TempIDMapping: map[string]string{},
	}

	srv := httptest.NewServer(mockHandlerWithRequestCapture(t, resp, &capturedReq))
	defer srv.Close()

	s, cleanup := openTestStore(t)
	defer cleanup()

	// Fresh store — no sync token set.
	token, err := s.GetSyncToken()
	if err != nil {
		t.Fatalf("GetSyncToken: %v", err)
	}
	if token != "" {
		t.Fatalf("expected empty token on fresh store, got %q", token)
	}

	client := newTestClient(srv.URL)
	if err := client.IncrementalSync(context.Background(), s); err != nil {
		t.Fatalf("IncrementalSync: %v", err)
	}

	// Verify that a full sync request (sync_token="*") was sent.
	if capturedReq.SyncToken != "*" {
		t.Errorf("expected sync_token=%q, got %q", "*", capturedReq.SyncToken)
	}

	// Verify the sync token was persisted.
	stored, err := s.GetSyncToken()
	if err != nil {
		t.Fatalf("GetSyncToken after sync: %v", err)
	}
	if stored != "fresh-token" {
		t.Errorf("stored sync token: got=%q, want=%q", stored, "fresh-token")
	}
}

func TestAuthFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	s, cleanup := openTestStore(t)
	defer cleanup()

	client := newTestClient(srv.URL)
	err := client.FullSync(context.Background(), s)
	if err == nil {
		t.Fatal("expected error for 401, got nil")
	}
	if !errors.Is(err, ErrAuthFailed) {
		t.Errorf("expected ErrAuthFailed, got: %v", err)
	}
}

func TestNetworkError(t *testing.T) {
	s, cleanup := openTestStore(t)
	defer cleanup()

	// Create a client pointing to a URL that will fail to connect.
	// Use a port that nothing is listening on.
	client := NewClient(ClientConfig{
		Token:    "test-token",
		Endpoint: "http://127.0.0.1:1/sync",
	})

	err := client.FullSync(context.Background(), s)
	if err == nil {
		t.Fatal("expected error for network failure, got nil")
	}
	if !errors.Is(err, ErrSyncFailed) {
		t.Errorf("expected ErrSyncFailed, got: %v", err)
	}
}

func TestTempIDMapping(t *testing.T) {
	s, cleanup := openTestStore(t)
	defer cleanup()

	// Create a project with a temp ID.
	tmpProjID := "tmp-proj-1"
	realProjID := "proj-real-1"
	if err := s.PutProject(&model.Project{
		ID:   tmpProjID,
		Name: "My Project",
	}); err != nil {
		t.Fatalf("PutProject: %v", err)
	}

	// Create a section referencing the temp project ID.
	tmpSecID := "tmp-sec-1"
	realSecID := "sec-real-1"
	if err := s.PutSection(&model.Section{
		ID:        tmpSecID,
		Name:      "Morning",
		ProjectID: tmpProjID,
	}); err != nil {
		t.Fatalf("PutSection: %v", err)
	}

	// Create a task with a temp ID referencing the temp project and section.
	tmpTaskID := "tmp-task-1"
	realTaskID := "task-real-1"
	if err := s.PutTask(&model.Task{
		ID:        tmpTaskID,
		Content:   "Do something",
		ProjectID: tmpProjID,
		SectionID: tmpSecID,
	}); err != nil {
		t.Fatalf("PutTask: %v", err)
	}

	// Create another task that references the temp task as parent.
	if err := s.PutTask(&model.Task{
		ID:       "subtask-1",
		Content:  "Sub-task",
		ParentID: tmpTaskID,
	}); err != nil {
		t.Fatalf("PutTask subtask: %v", err)
	}

	// Resolve temp IDs.
	mapping := map[string]string{
		tmpProjID: realProjID,
		tmpSecID:  realSecID,
		tmpTaskID: realTaskID,
	}
	if err := resolveTempIDs(s, mapping); err != nil {
		t.Fatalf("resolveTempIDs: %v", err)
	}

	// Verify the project ID was updated.
	proj, err := s.GetProject(realProjID)
	if err != nil {
		t.Fatalf("GetProject real ID: %v", err)
	}
	if proj.Name != "My Project" {
		t.Errorf("project name: got=%q, want=%q", proj.Name, "My Project")
	}

	// Verify the old temp project ID is gone.
	_, err = s.GetProject(tmpProjID)
	if err == nil {
		t.Error("expected error for temp project ID, but it still exists")
	}

	// Verify the section ID was updated and its ProjectID was remapped.
	sec, err := s.GetSection(realSecID)
	if err != nil {
		t.Fatalf("GetSection real ID: %v", err)
	}
	if sec.ProjectID != realProjID {
		t.Errorf("section ProjectID: got=%q, want=%q", sec.ProjectID, realProjID)
	}

	// Verify the task ID was updated and references were remapped.
	task, err := s.GetTask(realTaskID)
	if err != nil {
		t.Fatalf("GetTask real ID: %v", err)
	}
	if task.Content != "Do something" {
		t.Errorf("task content: got=%q, want=%q", task.Content, "Do something")
	}
	if task.ProjectID != realProjID {
		t.Errorf("task ProjectID: got=%q, want=%q", task.ProjectID, realProjID)
	}
	if task.SectionID != realSecID {
		t.Errorf("task SectionID: got=%q, want=%q", task.SectionID, realSecID)
	}

	// Verify the sub-task's ParentID was remapped.
	subtask, err := s.GetTask("subtask-1")
	if err != nil {
		t.Fatalf("GetTask subtask: %v", err)
	}
	if subtask.ParentID != realTaskID {
		t.Errorf("subtask ParentID: got=%q, want=%q", subtask.ParentID, realTaskID)
	}

	// Verify the old temp task ID is gone.
	_, err = s.GetTask(tmpTaskID)
	if err == nil {
		t.Error("expected error for temp task ID, but it still exists")
	}
}

func TestDeletedEntities(t *testing.T) {
	// Phase 1: FullSync to populate the store with entities.
	fullResp := SyncResponse{
		SyncToken: "token-v1",
		Items: []model.Task{
			{ID: "item-1", Content: "Keep me", ProjectID: "proj-1"},
			{ID: "item-2", Content: "Delete me", ProjectID: "proj-1"},
		},
		Projects: []model.Project{
			{ID: "proj-1", Name: "Keep"},
			{ID: "proj-2", Name: "Delete"},
		},
		Sections: []model.Section{
			{ID: "sec-1", Name: "Keep section", ProjectID: "proj-1"},
			{ID: "sec-2", Name: "Delete section", ProjectID: "proj-2"},
		},
		Labels: []model.Label{
			{ID: "lbl-1", Name: "keep-label"},
			{ID: "lbl-2", Name: "delete-label"},
		},
		Filters: []model.Filter{
			{ID: "flt-1", Name: "Keep filter", Query: "today"},
			{ID: "flt-2", Name: "Delete filter", Query: "p1"},
		},
		TempIDMapping: map[string]string{},
	}

	srv := httptest.NewServer(mockHandler(t, fullResp))
	s, cleanup := openTestStore(t)
	defer cleanup()

	client := newTestClient(srv.URL)
	if err := client.FullSync(context.Background(), s); err != nil {
		t.Fatalf("FullSync: %v", err)
	}
	srv.Close()

	// Verify initial state: 2 of each entity type.
	if tasks, _ := s.ListTasks(); len(tasks) != 2 {
		t.Fatalf("expected 2 tasks after full sync, got %d", len(tasks))
	}
	if projects, _ := s.ListProjects(); len(projects) != 2 {
		t.Fatalf("expected 2 projects after full sync, got %d", len(projects))
	}
	if sections, _ := s.ListSections(); len(sections) != 2 {
		t.Fatalf("expected 2 sections after full sync, got %d", len(sections))
	}
	if labels, _ := s.ListLabels(); len(labels) != 2 {
		t.Fatalf("expected 2 labels after full sync, got %d", len(labels))
	}
	if filters, _ := s.ListFilters(); len(filters) != 2 {
		t.Fatalf("expected 2 filters after full sync, got %d", len(filters))
	}

	// Phase 2: IncrementalSync with some entities marked as deleted.
	incrResp := SyncResponse{
		SyncToken: "token-v2",
		Items: []model.Task{
			{ID: "item-2", Content: "Delete me", IsDeleted: true},
		},
		Projects: []model.Project{
			{ID: "proj-2", Name: "Delete", IsDeleted: true},
		},
		Sections: []model.Section{
			{ID: "sec-2", Name: "Delete section", ProjectID: "proj-2", IsDeleted: true},
		},
		Labels: []model.Label{
			{ID: "lbl-2", Name: "delete-label", IsDeleted: true},
		},
		Filters: []model.Filter{
			{ID: "flt-2", Name: "Delete filter", Query: "p1", IsDeleted: true},
		},
		TempIDMapping: map[string]string{},
	}

	srv2 := httptest.NewServer(mockHandler(t, incrResp))
	defer srv2.Close()

	client2 := newTestClient(srv2.URL)
	if err := client2.IncrementalSync(context.Background(), s); err != nil {
		t.Fatalf("IncrementalSync: %v", err)
	}

	// Verify deleted task is gone, kept task remains.
	if _, err := s.GetTask("item-2"); err == nil {
		t.Error("expected item-2 to be deleted, but it still exists")
	}
	if task, err := s.GetTask("item-1"); err != nil {
		t.Errorf("expected item-1 to exist: %v", err)
	} else if task.Content != "Keep me" {
		t.Errorf("item-1 content: got=%q, want=%q", task.Content, "Keep me")
	}

	// Verify deleted project is gone, kept project remains.
	if _, err := s.GetProject("proj-2"); err == nil {
		t.Error("expected proj-2 to be deleted, but it still exists")
	}
	if proj, err := s.GetProject("proj-1"); err != nil {
		t.Errorf("expected proj-1 to exist: %v", err)
	} else if proj.Name != "Keep" {
		t.Errorf("proj-1 name: got=%q, want=%q", proj.Name, "Keep")
	}

	// Verify deleted section is gone, kept section remains.
	if _, err := s.GetSection("sec-2"); err == nil {
		t.Error("expected sec-2 to be deleted, but it still exists")
	}
	if sec, err := s.GetSection("sec-1"); err != nil {
		t.Errorf("expected sec-1 to exist: %v", err)
	} else if sec.Name != "Keep section" {
		t.Errorf("sec-1 name: got=%q, want=%q", sec.Name, "Keep section")
	}

	// Verify deleted label is gone, kept label remains.
	if _, err := s.GetLabel("lbl-2"); err == nil {
		t.Error("expected lbl-2 to be deleted, but it still exists")
	}
	if lbl, err := s.GetLabel("lbl-1"); err != nil {
		t.Errorf("expected lbl-1 to exist: %v", err)
	} else if lbl.Name != "keep-label" {
		t.Errorf("lbl-1 name: got=%q, want=%q", lbl.Name, "keep-label")
	}

	// Verify deleted filter is gone, kept filter remains.
	if _, err := s.GetFilter("flt-2"); err == nil {
		t.Error("expected flt-2 to be deleted, but it still exists")
	}
	if flt, err := s.GetFilter("flt-1"); err != nil {
		t.Errorf("expected flt-1 to exist: %v", err)
	} else if flt.Name != "Keep filter" {
		t.Errorf("flt-1 name: got=%q, want=%q", flt.Name, "Keep filter")
	}

	// Verify sync token was updated.
	token, err := s.GetSyncToken()
	if err != nil {
		t.Fatalf("GetSyncToken: %v", err)
	}
	if token != "token-v2" {
		t.Errorf("sync token: got=%q, want=%q", token, "token-v2")
	}
}

// TestIncrementalSync_Non401HTTPError verifies that a non-401 HTTP error
// results in ErrSyncFailed.
func TestIncrementalSync_Non401HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error": "internal"}`)
	}))
	defer srv.Close()

	s, cleanup := openTestStore(t)
	defer cleanup()

	// Set a sync token so IncrementalSync doesn't delegate to FullSync.
	if err := s.SetSyncToken("some-token"); err != nil {
		t.Fatalf("SetSyncToken: %v", err)
	}

	client := newTestClient(srv.URL)
	err := client.IncrementalSync(context.Background(), s)
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
	if !errors.Is(err, ErrSyncFailed) {
		t.Errorf("expected ErrSyncFailed, got: %v", err)
	}
}