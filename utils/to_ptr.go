package utils

import (
	"time"
)

// ToPtr is a generic function that converts string and time.Time literals to pointers
func ToPtr[T string|time.Time](v T) *T {
    return &v
}