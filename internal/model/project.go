package model

// Project represents a Todoist project.
type Project struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	IsFavorite bool   `json:"is_favorite"`
	IsInbox    bool   `json:"is_inbox"`
	ParentID   string `json:"parent_id,omitempty"`
	ViewStyle  string `json:"view_style"`
	IsDeleted  bool   `json:"is_deleted"`
	IsArchived bool   `json:"is_archived"`
	UserID     string `json:"user_id,omitempty"`
}

// Bucket returns the bbolt bucket name for projects.
func (Project) Bucket() string { return "projects" }