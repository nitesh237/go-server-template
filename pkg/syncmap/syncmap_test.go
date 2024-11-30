package syncmap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap_Delete(t *testing.T) {
	t.Parallel()
	m := Map[string, int]{}
	m.Store("foo", 100)
	m.Delete("foo")
	a := require.New(t)
	a.Equal(m.Get("foo"), 0, "key not deleted")
}

func TestMap_Get(t *testing.T) {
	t.Parallel()
	m := Map[string, int]{}
	m.Store("foo", 100)
	a := require.New(t)
	a.Equal(m.Get("foo"), 100, "failed to get key")
}

func TestMap_LoadAndStore(t *testing.T) {
	t.Parallel()
	m := Map[string, int]{}
	m.Store("foo", 100)
	a := require.New(t)
	val, ok := m.Load("foo")
	a.Equal(ok, true, "key doesn't exist")
	a.Equal(val, 100, "value is different")
	_, ok = m.Load("bar")
	a.Equal(ok, false, "unknown key exists")
}

func TestMap_LoadAndDelete(t *testing.T) {
	t.Parallel()
	m := Map[string, int]{}
	m.Store("foo", 100)
	a := require.New(t)
	val, loaded := m.LoadAndDelete("foo")
	a.Equal(loaded, true, "vaule not loaded")
	a.Equal(val, 100, "unexpected value")
	_, ok := m.Load("foo")
	a.Equal(ok, false, "key not deleted")
	_, loaded = m.LoadAndDelete("bar")
	a.Equal(loaded, false, "invalid value loaded")
}

func TestMap_LoadOrStore(t *testing.T) {
	t.Parallel()
	m := Map[string, int]{}
	actual, loaded := m.LoadOrStore("foo", 100)
	a := require.New(t)
	a.Equal(loaded, false, "unexpected value loaded")
	a.Equal(actual, 100, "unexpected value")
	actual, loaded = m.LoadOrStore("foo", 200)
	a.Equal(loaded, true, "unexpected value loaded")
	a.Equal(actual, 100, "unexpected value")
}

func TestMap_Range(t *testing.T) {
	t.Parallel()
	m := Map[string, int]{}
	expectedMap := map[string]int{
		"foo":    100,
		"bar":    200,
		"foobar": 300,
		"barfoo": 400,
	}
	for k, v := range expectedMap {
		m.Store(k, v)
	}
	got := map[string]int{}
	m.Range(func(k string, v int) bool {
		got[k] = v
		return true
	})
	a := require.New(t)
	a.Equal(expectedMap, got, "mismatch in kv pairs")
}
