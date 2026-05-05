package model

// DueDate represents a due date attached to a task.
type DueDate struct {
	Date        string `json:"date"`
	IsRecurring bool   `json:"is_recurring"`
	Datetime    string `json:"datetime,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
}

// Task represents a Todoist task item.
type Task struct {
	ID           string   `json:"id"`
	Content      string   `json:"content"`
	Description  string   `json:"description"`
	ProjectID    string   `json:"project_id"`
	SectionID    string   `json:"section_id"`
	ParentID    string   `json:"parent_id,omitempty"`
	Labels       []string `json:"labels"`
	Priority     int      `json:"priority"`
	Due          *DueDate `json:"due,omitempty"`
	Completed    bool     `json:"completed"`
	Checked      bool     `json:"checked"`
	IsDeleted    bool     `json:"is_deleted"`
	AddedAt      string   `json:"added_at,omitempty"`
	CompletedAt  string   `json:"completed_at,omitempty"`
	UserID       string   `json:"user_id,omitempty"`
}

// Bucket returns the bbolt bucket name for tasks.
func (Task) Bucket() string { return "tasks" }