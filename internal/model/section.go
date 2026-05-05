package model

// Section represents a section within a project.
type Section struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ProjectID    string `json:"project_id"`
	SectionOrder int    `json:"section_order"`
	IsDeleted    bool   `json:"is_deleted"`
	IsArchived   bool   `json:"is_archived"`
}

// Bucket returns the bbolt bucket name for sections.
func (Section) Bucket() string { return "sections" }