package model

// SyncMeta stores key-value metadata for sync state (e.g. sync tokens).
type SyncMeta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Bucket returns the bbolt bucket name for sync metadata.
func (SyncMeta) Bucket() string { return "sync_meta" }