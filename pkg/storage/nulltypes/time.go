package nulltypes

import (
	"database/sql"
	"time"
)

// A wrapper struct over sql.NullTime
type NullTime struct {
	sql.NullTime
}

// Returns time if valid is true else returns time
func (ns NullTime) GetValue() time.Time {
	if ns.Valid {
		return ns.Time
	}
	return time.Time{}
}

// Factory method to initialize a NullTime with valid true
func NewNullTime(val time.Time) NullTime {
	empty := time.Time{}
	if val == empty {
		return NullTime{struct {
			Time  time.Time
			Valid bool
		}{Time: val, Valid: false}}
	}
	return NullTime{struct {
		Time  time.Time
		Valid bool
	}{Time: val, Valid: true}}
}
