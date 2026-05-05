---
# todoist-tui-mwib
title: Temp ID mapping resolution
status: todo
type: task
created_at: 2026-05-03T15:09:55Z
updated_at: 2026-05-03T15:09:55Z
parent: todoist-tui-efa6
blocked_by:
    - todoist-tui-hgg5
---

## Description

Process the temp_id_mapping from Sync API responses to replace tmp- prefixed placeholder IDs with real server IDs.

## Requirements

- After each sync response, walk temp_id_mapping map
- For each mapping (temp_id → real_id), find and update entities in the store that reference the temp_id
- Replace temp_id in: Task.ID, Task.ProjectID, Task.SectionID, Task.ParentID, and any other reference fields
- Remove orphaned temp entries from store
- This is needed to support optimistic UI updates in M2

## Examples

temp_id_mapping: {'tmp-abc123': '987654321'}
→ Find task with ID 'tmp-abc123', update its ID to '987654321'
→ Also update any references to 'tmp-abc123' in other entities

## Acceptance Criteria

- temp_id_mapping entries replace tmp- IDs with real IDs in store
- References to tmp- IDs in other entities are updated
- No orphaned tmp- entries remain after resolution
