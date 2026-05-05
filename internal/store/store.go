package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	bolt "go.etcd.io/bbolt"
)

// ErrStoreNotOpen is returned when a store operation is attempted on a
// database that has not been opened or has been closed.
var ErrStoreNotOpen = errors.New("store: database not open")

// DefaultMode is the file mode used when creating the bbolt database file.
const DefaultMode os.FileMode = 0o600

// StoreConfig holds configuration for opening a bbolt store.
type StoreConfig struct {
	// Path is the filesystem path to the bbolt database file.
	// If empty, DBPath() is used to derive the default location.
	Path string
	// Mode is the file permission mode for the database file.
	// Defaults to 0o600 if zero.
	Mode os.FileMode
}

// Store wraps a bbolt.DB and provides open/close lifecycle management.
type Store struct {
	db *bolt.DB
}

// DBPath returns the default database file path following XDG conventions.
// It uses $XDG_DATA_HOME if set, otherwise falls back to ~/.local/share.
func DBPath() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "todoist-tui", "todoist.db")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		// Unlikely to fail in practice; use a reasonable fallback.
		home = filepath.Join("/tmp", "todoist-tui")
	}
	return filepath.Join(home, ".local", "share", "todoist-tui", "todoist.db")
}

// New opens (or creates) a bbolt database using the given configuration.
// It creates parent directories if they do not exist.
func New(cfg StoreConfig) (*Store, error) {
	path := cfg.Path
	if path == "" {
		path = DBPath()
	}

	mode := cfg.Mode
	if mode == 0 {
		mode = DefaultMode
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("store: create db directory: %w", err)
	}

	db, err := bolt.Open(path, mode, nil)
	if err != nil {
		return nil, fmt.Errorf("store: open db: %w", err)
	}

	s := &Store{db: db}
	if err := s.ensureBuckets(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("store: create buckets: %w", err)
	}

	return s, nil
}

// Close cleanly shuts down the underlying bbolt database.
// It returns ErrStoreNotOpen if the store was never opened or already closed.
func (s *Store) Close() error {
	if s.db == nil {
		return ErrStoreNotOpen
	}
	err := s.db.Close()
	s.db = nil
	return err
}

// buckets is the set of all bbolt bucket names used by the store.
var buckets = []string{
	"projects",
	"tasks",
	"sections",
	"labels",
	"filters",
	"sync_meta",
}

// ensureBuckets creates all required buckets if they do not already exist.
func (s *Store) ensureBuckets() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, name := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
				return fmt.Errorf("store: create bucket %q: %w", name, err)
			}
		}
		return nil
	})
}

// Put JSON-marshals value and stores it under key in the given bucket.
func (s *Store) Put(bucket, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("store: marshal value: %w", err)
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("store: bucket %q not found", bucket)
		}
		return b.Put([]byte(key), data)
	})
}

// Get retrieves the value stored under key in bucket and JSON-unmarshals it
// into out, which must be a pointer.
func (s *Store) Get(bucket, key string, out any) error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("store: bucket %q not found", bucket)
		}
		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("store: key %q not found in bucket %q", key, bucket)
		}
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("store: unmarshal value: %w", err)
		}
		return nil
	})
}

// Delete removes the value stored under key from the given bucket.
func (s *Store) Delete(bucket, key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("store: bucket %q not found", bucket)
		}
		return b.Delete([]byte(key))
	})
}

// List scans the entire bucket and JSON-unmarshals every item into out,
// which must be a non-nil pointer to a slice of the target type.
func (s *Store) List(bucket string, out any) error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("store: bucket %q not found", bucket)
		}
		slice := reflect.ValueOf(out).Elem()
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			item := reflect.New(slice.Type().Elem())
			if err := json.Unmarshal(v, item.Interface()); err != nil {
				return fmt.Errorf("store: unmarshal item %q: %w", string(k), err)
			}
			slice.Set(reflect.Append(slice, item.Elem()))
		}
		return nil
	})
}