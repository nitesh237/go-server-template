package nulltypes

import "database/sql"

type NullInt64 struct {
	sql.NullInt64
}

// Method to get value
// Returns value int64 if valid is true else returns 0
func (n NullInt64) GetValue() int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}

// Factory method to initialize a NullInt64 with valid true
func NewNullInt64(val int64) NullInt64 {
	if val == 0 {
		return NullInt64{struct {
			Int64 int64
			Valid bool
		}{Int64: val, Valid: false}}
	}
	return NullInt64{struct {
		Int64 int64
		Valid bool
	}{Int64: val, Valid: true}}
}

type NullInt16 struct {
	sql.NullInt16
}

// Method to get value
// Returns value int16 if valid is true else returns 0
func (n NullInt16) GetValue() int16 {
	if n.Valid {
		return n.Int16
	}
	return 0
}

// Factory method to initialize a NullInt16 with valid true
func NewNullInt16(val int16) NullInt16 {
	if val == 0 {
		return NullInt16{struct {
			Int16 int16
			Valid bool
		}{Int16: val, Valid: false}}
	}
	return NullInt16{struct {
		Int16 int16
		Valid bool
	}{Int16: val, Valid: true}}
}
