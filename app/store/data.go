package store

import (
	"time"
)

type ValueType int

const (
	StringType ValueType = iota
	ListType
	StreamType
)

type Data struct {
	Typ       ValueType
	Value     any
	ExpiresAt time.Time
}

func (d Data) isExpired() bool {
	return !d.ExpiresAt.IsZero() && time.Now().After(d.ExpiresAt)
}
