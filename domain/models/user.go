package models

import "time"

type User struct {
	ID             string
	Username       string
	Email          string
	FirstName      string
	LastName       string
	Password       string
	AvatarURL      string
	About          string
	Status         string
	IsActivated    bool
	Deleted        bool
	FollowerCount  int64
	FollowingCount int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Version        int
}
