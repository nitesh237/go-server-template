package nulltypes

import "database/sql"

type NullFloat64 struct {
	sql.NullFloat64
}

// GetValue Method to get value
// Returns value float64 if valid is true else returns 0
func (n NullFloat64) GetValue() float64 {
	if n.Valid {
		return n.Float64
	}
	return 0
}

// Factory method to initialize a NullFloat64 with valid true
func NewNullFloat64(val float64) NullFloat64 {
	if val == 0 {
		return NullFloat64{struct {
			Float64 float64
			Valid   bool
		}{Float64: val, Valid: false}}
	}
	return NullFloat64{struct {
		Float64 float64
		Valid   bool
	}{Float64: val, Valid: true}}
}
