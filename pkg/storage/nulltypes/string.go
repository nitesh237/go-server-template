package nulltypes

import (
	"database/sql"
	"encoding/json"
)

// A wrapper struct over sql.NullString
type NullString struct {
	sql.NullString
}

// Method to get value
// Returns value string if valid is true else returns empty string
func (ns NullString) GetValue() string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// Factory method to initialize a NullString with valid true
func NewNullString(val string) NullString {
	if val == "" {
		return NullString{struct {
			String string
			Valid  bool
		}{String: val, Valid: false}}
	}
	return NullString{struct {
		String string
		Valid  bool
	}{String: val, Valid: true}}
}

// ref - https://gist.github.com/keidrun/d1b2791f840753e25070771b857af7ba
// ref - https://stackoverflow.com/questions/33072172/how-can-i-work-with-sql-null-values-and-json-in-a-good-way
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ns == nil {
		return json.Marshal(nil)
	}
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil && *s != "" {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

func EmptyNullString() NullString {
	return NullString{struct {
		String string
		Valid  bool
	}{String: "", Valid: true}}
}
