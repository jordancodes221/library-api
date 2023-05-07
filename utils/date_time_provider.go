package utils

import (
	"time"
)

type DateTimeProvider interface {
	GetCurrentTime()  *time.Time
}

type ProductionDateTimeProvider struct {}

func (p *ProductionDateTimeProvider) GetCurrentTime() (*time.Time) {
	return ToPtr(time.Now())
}

type TestingDateTimeProvider struct {
	arbitraryTime time.Time
}

func (t *TestingDateTimeProvider) GetCurrentTime() (*time.Time) {
	return ToPtr(t.arbitraryTime)
}