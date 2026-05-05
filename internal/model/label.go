package model

// Label represents a Todoist label.
type Label struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	IsFavorite bool   `json:"is_favorite"`
	IsDeleted  bool   `json:"is_deleted"`
	ItemOrder  int    `json:"item_order"`
}

// Bucket returns the bbolt bucket name for labels.
func (Label) Bucket() string { return "labels" }