---
# todoist-tui-xbbs
title: Bucket definitions and CRUD helpers
status: todo
type: task
created_at: 2026-05-03T15:09:20Z
updated_at: 2026-05-03T15:09:20Z
parent: todoist-tui-g83v
blocked_by:
    - todoist-tui-p2ij
---

## Description

Create buckets on init and implement generic Put, Get, Delete, and List helpers that JSON-serialize/deserialize domain models.

## Requirements

- Buckets: 'projects', 'tasks', 'sections', 'labels', 'filters', 'sync_meta'
- Store.ensureBuckets() creates all buckets in a WriteTxn called during New()
- Put(bucket, key string, value any) error — JSON-marshal value, store under key
- Get(bucket, key string, out any) error — JSON-unmarshal into out pointer
- Delete(bucket, key string) error
- List(bucket string, out any) error — scan entire bucket, append all items into out slice
- All operations use bbolt transactions (View/Update)
- Typed convenience methods or generic helpers for each model type

## Examples

```go
func (s *Store) Put(bucket, key string, value any) error {
    data, err := json.Marshal(value)
    if err != nil {
        return fmt.Errorf("marshal value: %w", err)
    }
    return s.db.Update(func(tx *bbolt.Tx) error {
        b := tx.Bucket([]byte(bucket))
        return b.Put([]byte(key), data)
    })
}
```

## Acceptance Criteria

- All 6 buckets created on Open
- Put/Get round-trip preserves data
- List returns all items in a bucket
- Delete removes items
- 'go build ./...' and 'go vet ./...' pass
