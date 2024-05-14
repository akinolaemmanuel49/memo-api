package models

import (
	"database/sql"
	"time"
)

type Follow struct {
	ID         string
	SubjectID  string
	FollowerID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Version    int
}

type Comment struct {
	ID          string
	OwnerID     string
	MemoID      string
	ParentID    sql.NullString
	CommentType string
	Content     string
	Caption     sql.NullString
	Transcript  sql.NullString
	Likes       int64
	Deleted     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}

type CommentParentChild struct {
	ID        string
	ParentID  string
	ChildID   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}
