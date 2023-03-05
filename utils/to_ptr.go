package utils

import (
	"time"
)

// Generic function converts literals to pointers
func ToPtr[T string|time.Time](v T) *T {
    return &v
}