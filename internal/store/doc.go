// Package store provides the bbolt storage layer for projects, tasks, labels, and sync state.
package store

import (
	_ "go.etcd.io/bbolt" // Embedded KV store
)
