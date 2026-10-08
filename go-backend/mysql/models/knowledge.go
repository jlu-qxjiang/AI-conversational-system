package models

import (
	"database/sql"
)

type Knowledge struct {
	ID        int64          `json:"id" db:"id"`
	Title     string         `json:"title" db:"title"`
	FileName  string         `json:"file_name" db:"file_name"`
	Path      string         `json:"path" db:"path"`
	Content   sql.NullString `json:"content" db:"content"`
	CreatedAt int64          `json:"created_at" db:"created_at"` // 时间戳 ✅
	UpdatedAt int64          `json:"updated_at" db:"updated_at"` // 时间戳 ✅
	Source    string         `json:"source" db:"source"`
}
