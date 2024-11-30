package nulltypes

import "database/sql"

type NullBool struct {
	sql.NullBool
}

// Method to get value
// Returns value bool if valid is true else returns false
func (n NullBool) GetValue() bool {
	if n.Valid {
		return n.Bool
	}
	return false
}

// Factory method to initialize a NullBool with valid true
func NewNullBool(val bool) NullBool {
	if !val {
		return NullBool{struct {
			Bool  bool
			Valid bool
		}{Bool: val, Valid: false}}
	}
	return NullBool{struct {
		Bool  bool
		Valid bool
	}{Bool: val, Valid: true}}
}
