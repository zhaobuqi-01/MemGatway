package cache

import "time"

const (
	defaultExpiration = 24 * time.Hour
	cleanupInterval   = 48 * time.Hour
)
