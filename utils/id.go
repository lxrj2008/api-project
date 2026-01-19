package utils

import "github.com/google/uuid"

// NewID returns a string UUIDv4.
func NewID() string {
	return uuid.NewString()
}
