package store

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// syncMetaBucket is the bbolt bucket name for sync metadata.
const syncMetaBucket = "sync_meta"

// syncTokenKey is the key used to store the sync token within the sync_meta bucket.
const syncTokenKey = "sync_token"

// lastSyncTimeKey is the key used to store the last successful sync timestamp.
const lastSyncTimeKey = "last_sync_time"

// GetSyncToken returns the stored sync token, or an empty string if none has
// been stored yet (e.g. on first launch, which triggers a full sync).
func (s *Store) GetSyncToken() (string, error) {
	var token string
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(syncMetaBucket))
		if b == nil {
			// Bucket doesn't exist yet — no token stored.
			return nil
		}
		v := b.Get([]byte(syncTokenKey))
		if v != nil {
			token = string(v)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("store: get sync token: %w", err)
	}
	return token, nil
}

// SetSyncToken persists the sync token so that incremental syncs can resume
// from the last known state.
func (s *Store) SetSyncToken(token string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(syncMetaBucket))
		if err != nil {
			return fmt.Errorf("create bucket %q: %w", syncMetaBucket, err)
		}
		return b.Put([]byte(syncTokenKey), []byte(token))
	})
	if err != nil {
		return fmt.Errorf("store: set sync token: %w", err)
	}
	return nil
}

// GetLastSyncTime returns the timestamp of the last successful sync as an
// RFC3339 string, or an empty string if no sync has completed yet.
func (s *Store) GetLastSyncTime() (string, error) {
	var ts string
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(syncMetaBucket))
		if b == nil {
			return nil
		}
		v := b.Get([]byte(lastSyncTimeKey))
		if v != nil {
			ts = string(v)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("store: get last sync time: %w", err)
	}
	return ts, nil
}

// SetLastSyncTime records the timestamp of the last successful sync.
// It stores the time in RFC3339 format.
func (s *Store) SetLastSyncTime(t time.Time) error {
	ts := t.Format(time.RFC3339)
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(syncMetaBucket))
		if err != nil {
			return fmt.Errorf("create bucket %q: %w", syncMetaBucket, err)
		}
		return b.Put([]byte(lastSyncTimeKey), []byte(ts))
	})
	if err != nil {
		return fmt.Errorf("store: set last sync time: %w", err)
	}
	return nil
}