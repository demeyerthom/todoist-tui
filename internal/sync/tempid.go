package sync

import (
	"fmt"
	"strings"

	"github.com/demeyerthom/todoist-tui/internal/store"
)

// tmpPrefix is the prefix used for temporary IDs in optimistic updates.
const tmpPrefix = "tmp-"

// isTempID reports whether id is a temporary placeholder ID.
func isTempID(id string) bool {
	return strings.HasPrefix(id, tmpPrefix)
}

// resolveTempIDs walks the temp_id_mapping from a Sync API response and
// replaces tmp- prefixed placeholder IDs with the real server-assigned IDs.
// It updates both the owning entity's ID and any cross-entity references
// (Task.ProjectID, Task.SectionID, Task.ParentID, Section.ProjectID).
// After resolution, any orphaned entities still carrying a tmp- ID are deleted.
func resolveTempIDs(s *store.Store, mapping map[string]string) error {
	for tempID, realID := range mapping {
		if !isTempID(tempID) {
			continue
		}

		if err := resolveEntityID(s, tempID, realID); err != nil {
			return fmt.Errorf("sync: resolve temp id %s: %w", tempID, err)
		}

		if err := resolveTaskRefs(s, tempID, realID); err != nil {
			return fmt.Errorf("sync: resolve task refs for %s: %w", tempID, err)
		}

		if err := resolveSectionRefs(s, tempID, realID); err != nil {
			return fmt.Errorf("sync: resolve section refs for %s: %w", tempID, err)
		}
	}

	if err := removeOrphanedTempIDs(s); err != nil {
		return fmt.Errorf("sync: remove orphaned temp ids: %w", err)
	}

	return nil
}

// resolveEntityID finds the entity whose ID equals tempID across all entity
// types and updates it to realID. If no entity is found, it is skipped
// (it may have been deleted server-side).
func resolveEntityID(s *store.Store, tempID, realID string) error {
	// Try Task.
	if t, err := s.GetTask(tempID); err == nil {
		t.ID = realID
		if err := s.PutTask(t); err != nil {
			return fmt.Errorf("put task: %w", err)
		}
		// Delete the old temp-keyed entry since Put stores by ID.
		if err := s.DeleteTask(tempID); err != nil {
			// The new entry was already written under realID; the old key
			// may have been overwritten or will be cleaned up by orphan removal.
			_ = err
		}
		return nil
	}

	// Try Project.
	if p, err := s.GetProject(tempID); err == nil {
		p.ID = realID
		if err := s.PutProject(p); err != nil {
			return fmt.Errorf("put project: %w", err)
		}
		_ = s.DeleteProject(tempID)
		return nil
	}

	// Try Section.
	if sec, err := s.GetSection(tempID); err == nil {
		sec.ID = realID
		if err := s.PutSection(sec); err != nil {
			return fmt.Errorf("put section: %w", err)
		}
		_ = s.DeleteSection(tempID)
		return nil
	}

	// Try Label.
	if l, err := s.GetLabel(tempID); err == nil {
		l.ID = realID
		if err := s.PutLabel(l); err != nil {
			return fmt.Errorf("put label: %w", err)
		}
		_ = s.DeleteLabel(tempID)
		return nil
	}

	// Try Filter.
	if f, err := s.GetFilter(tempID); err == nil {
		f.ID = realID
		if err := s.PutFilter(f); err != nil {
			return fmt.Errorf("put filter: %w", err)
		}
		_ = s.DeleteFilter(tempID)
		return nil
	}

	// Entity not found in any bucket — skip (may have been deleted server-side).
	return nil
}

// resolveTaskRefs scans all tasks and updates any reference field that equals
// oldID to newID. Reference fields: ProjectID, SectionID, ParentID.
func resolveTaskRefs(s *store.Store, oldID, newID string) error {
	tasks, err := s.ListTasks()
	if err != nil {
		return fmt.Errorf("list tasks: %w", err)
	}

	for i := range tasks {
		t := &tasks[i]
		changed := false

		if t.ProjectID == oldID {
			t.ProjectID = newID
			changed = true
		}
		if t.SectionID == oldID {
			t.SectionID = newID
			changed = true
		}
		if t.ParentID == oldID {
			t.ParentID = newID
			changed = true
		}

		if changed {
			if err := s.PutTask(t); err != nil {
				return fmt.Errorf("put task %s: %w", t.ID, err)
			}
		}
	}

	return nil
}

// resolveSectionRefs scans all sections and updates ProjectID if it equals
// oldID to newID.
func resolveSectionRefs(s *store.Store, oldID, newID string) error {
	sections, err := s.ListSections()
	if err != nil {
		return fmt.Errorf("list sections: %w", err)
	}

	for i := range sections {
		sec := &sections[i]
		if sec.ProjectID == oldID {
			sec.ProjectID = newID
			if err := s.PutSection(sec); err != nil {
				return fmt.Errorf("put section %s: %w", sec.ID, err)
			}
		}
	}

	return nil
}

// removeOrphanedTempIDs scans all entity buckets and deletes any entity whose
// ID still has the tmp- prefix after resolution.
func removeOrphanedTempIDs(s *store.Store) error {
	// Tasks.
	tasks, err := s.ListTasks()
	if err != nil {
		return fmt.Errorf("list tasks: %w", err)
	}
	for _, t := range tasks {
		if isTempID(t.ID) {
			if err := s.DeleteTask(t.ID); err != nil {
				return fmt.Errorf("delete orphan task %s: %w", t.ID, err)
			}
		}
	}

	// Projects.
	projects, err := s.ListProjects()
	if err != nil {
		return fmt.Errorf("list projects: %w", err)
	}
	for _, p := range projects {
		if isTempID(p.ID) {
			if err := s.DeleteProject(p.ID); err != nil {
				return fmt.Errorf("delete orphan project %s: %w", p.ID, err)
			}
		}
	}

	// Sections.
	sections, err := s.ListSections()
	if err != nil {
		return fmt.Errorf("list sections: %w", err)
	}
	for _, sec := range sections {
		if isTempID(sec.ID) {
			if err := s.DeleteSection(sec.ID); err != nil {
				return fmt.Errorf("delete orphan section %s: %w", sec.ID, err)
			}
		}
	}

	// Labels.
	labels, err := s.ListLabels()
	if err != nil {
		return fmt.Errorf("list labels: %w", err)
	}
	for _, l := range labels {
		if isTempID(l.ID) {
			if err := s.DeleteLabel(l.ID); err != nil {
				return fmt.Errorf("delete orphan label %s: %w", l.ID, err)
			}
		}
	}

	// Filters.
	filters, err := s.ListFilters()
	if err != nil {
		return fmt.Errorf("list filters: %w", err)
	}
	for _, f := range filters {
		if isTempID(f.ID) {
			if err := s.DeleteFilter(f.ID); err != nil {
				return fmt.Errorf("delete orphan filter %s: %w", f.ID, err)
			}
		}
	}

	return nil
}
