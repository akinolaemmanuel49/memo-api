package helpers

import "time"

var (
	IdleTimeout  = 1 * time.Minute
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 10 * time.Second
	MaxAge       = 12 * time.Hour
)

const (
	DefaultPage         string = "1"
	DefaultPageSize     string = "10"
	AccessTokenDuration        = 365 * 24 * time.Hour
)
