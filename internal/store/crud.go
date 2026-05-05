package store

import (
	"fmt"

	"github.com/demeyerthom/todoist-tui/internal/model"
)

// PutProject stores a project by its ID.
func (s *Store) PutProject(p *model.Project) error {
	return s.Put(p.Bucket(), p.ID, p)
}

// GetProject retrieves a project by ID.
func (s *Store) GetProject(id string) (*model.Project, error) {
	var p model.Project
	if err := s.Get(model.Project{}.Bucket(), id, &p); err != nil {
		return nil, fmt.Errorf("store: get project: %w", err)
	}
	return &p, nil
}

// DeleteProject removes a project by ID.
func (s *Store) DeleteProject(id string) error {
	return s.Delete(model.Project{}.Bucket(), id)
}

// ListProjects returns all projects.
func (s *Store) ListProjects() ([]model.Project, error) {
	var projects []model.Project
	if err := s.List(model.Project{}.Bucket(), &projects); err != nil {
		return nil, fmt.Errorf("store: list projects: %w", err)
	}
	return projects, nil
}

// PutTask stores a task by its ID.
func (s *Store) PutTask(t *model.Task) error {
	return s.Put(t.Bucket(), t.ID, t)
}

// GetTask retrieves a task by ID.
func (s *Store) GetTask(id string) (*model.Task, error) {
	var t model.Task
	if err := s.Get(model.Task{}.Bucket(), id, &t); err != nil {
		return nil, fmt.Errorf("store: get task: %w", err)
	}
	return &t, nil
}

// DeleteTask removes a task by ID.
func (s *Store) DeleteTask(id string) error {
	return s.Delete(model.Task{}.Bucket(), id)
}

// ListTasks returns all tasks.
func (s *Store) ListTasks() ([]model.Task, error) {
	var tasks []model.Task
	if err := s.List(model.Task{}.Bucket(), &tasks); err != nil {
		return nil, fmt.Errorf("store: list tasks: %w", err)
	}
	return tasks, nil
}

// PutSection stores a section by its ID.
func (s *Store) PutSection(sec *model.Section) error {
	return s.Put(sec.Bucket(), sec.ID, sec)
}

// GetSection retrieves a section by ID.
func (s *Store) GetSection(id string) (*model.Section, error) {
	var sec model.Section
	if err := s.Get(model.Section{}.Bucket(), id, &sec); err != nil {
		return nil, fmt.Errorf("store: get section: %w", err)
	}
	return &sec, nil
}

// DeleteSection removes a section by ID.
func (s *Store) DeleteSection(id string) error {
	return s.Delete(model.Section{}.Bucket(), id)
}

// ListSections returns all sections.
func (s *Store) ListSections() ([]model.Section, error) {
	var sections []model.Section
	if err := s.List(model.Section{}.Bucket(), &sections); err != nil {
		return nil, fmt.Errorf("store: list sections: %w", err)
	}
	return sections, nil
}

// PutLabel stores a label by its ID.
func (s *Store) PutLabel(l *model.Label) error {
	return s.Put(l.Bucket(), l.ID, l)
}

// GetLabel retrieves a label by ID.
func (s *Store) GetLabel(id string) (*model.Label, error) {
	var l model.Label
	if err := s.Get(model.Label{}.Bucket(), id, &l); err != nil {
		return nil, fmt.Errorf("store: get label: %w", err)
	}
	return &l, nil
}

// DeleteLabel removes a label by ID.
func (s *Store) DeleteLabel(id string) error {
	return s.Delete(model.Label{}.Bucket(), id)
}

// ListLabels returns all labels.
func (s *Store) ListLabels() ([]model.Label, error) {
	var labels []model.Label
	if err := s.List(model.Label{}.Bucket(), &labels); err != nil {
		return nil, fmt.Errorf("store: list labels: %w", err)
	}
	return labels, nil
}

// PutFilter stores a filter by its ID.
func (s *Store) PutFilter(f *model.Filter) error {
	return s.Put(f.Bucket(), f.ID, f)
}

// GetFilter retrieves a filter by ID.
func (s *Store) GetFilter(id string) (*model.Filter, error) {
	var f model.Filter
	if err := s.Get(model.Filter{}.Bucket(), id, &f); err != nil {
		return nil, fmt.Errorf("store: get filter: %w", err)
	}
	return &f, nil
}

// DeleteFilter removes a filter by ID.
func (s *Store) DeleteFilter(id string) error {
	return s.Delete(model.Filter{}.Bucket(), id)
}

// ListFilters returns all filters.
func (s *Store) ListFilters() ([]model.Filter, error) {
	var filters []model.Filter
	if err := s.List(model.Filter{}.Bucket(), &filters); err != nil {
		return nil, fmt.Errorf("store: list filters: %w", err)
	}
	return filters, nil
}

// PutSyncMeta stores a sync metadata entry by its key.
func (s *Store) PutSyncMeta(m *model.SyncMeta) error {
	return s.Put(m.Bucket(), m.Key, m)
}

// GetSyncMeta retrieves a sync metadata entry by key.
func (s *Store) GetSyncMeta(key string) (*model.SyncMeta, error) {
	var m model.SyncMeta
	if err := s.Get(model.SyncMeta{}.Bucket(), key, &m); err != nil {
		return nil, fmt.Errorf("store: get sync meta: %w", err)
	}
	return &m, nil
}

// DeleteSyncMeta removes a sync metadata entry by key.
func (s *Store) DeleteSyncMeta(key string) error {
	return s.Delete(model.SyncMeta{}.Bucket(), key)
}

// ListSyncMeta returns all sync metadata entries.
func (s *Store) ListSyncMeta() ([]model.SyncMeta, error) {
	var items []model.SyncMeta
	if err := s.List(model.SyncMeta{}.Bucket(), &items); err != nil {
		return nil, fmt.Errorf("store: list sync meta: %w", err)
	}
	return items, nil
}