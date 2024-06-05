package models

import "time"

type Memo struct {
	ID          string
	Type        string
	Content     string
	ResourceURL string
	Description string
	Likes       int64
	Shares      int64
	Deleted     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	OwnerID     string
	Version     int
}

type Like struct {
	ID        string
	MemoID    string
	LikedBy   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}

type Share struct {
	ID        string
	MemoID    string
	SharedBy  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}
