package model

// Filter represents a Todoist filter.
type Filter struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Query     string `json:"query"`
	IsDeleted bool   `json:"is_deleted"`
}

// Bucket returns the bbolt bucket name for filters.
func (Filter) Bucket() string { return "filters" }