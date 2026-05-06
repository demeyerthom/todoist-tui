---
# todoist-tui-bge2
title: Implement resolveTempIDs in internal/sync/tempid.go
status: completed
type: task
priority: normal
created_at: 2026-05-06T10:01:28Z
updated_at: 2026-05-06T10:02:14Z
---

Create tempid.go with resolveTempIDs function that processes temp_id_mapping from Sync API responses, replacing tmp- prefixed placeholder IDs with real server IDs, updating cross-entity references, and cleaning up orphans.

## Summary of Changes

- Created `internal/sync/tempid.go` with the `resolveTempIDs` function
- `resolveTempIDs(s *store.Store, mapping map[string]string) error` — unexported, walks temp_id_mapping and replaces tmp- IDs with real IDs
- `resolveEntityID` — finds the entity owning a temp ID across all 5 entity types (Task, Project, Section, Label, Filter), updates its ID, and deletes the old temp-keyed entry
- `resolveTaskRefs` — scans all Tasks for references to oldID in ProjectID, SectionID, ParentID and updates them
- `resolveSectionRefs` — scans all Sections for references to oldID in ProjectID and updates them
- `removeOrphanedTempIDs` — scans all entity buckets and deletes any entity still carrying a tmp- prefix ID
- `isTempID` helper using `strings.HasPrefix(id, "tmp-")`
- All errors from store operations are returned; entities not found in any bucket are skipped
- `go build ./...`, `go vet ./...`, and `go test ./...` all pass
