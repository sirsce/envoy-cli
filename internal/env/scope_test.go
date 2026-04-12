package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeScope(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestScopeManager_AddGet(t *testing.T) {
	sm := NewScopeManager("")
	sm.Add("dev", makeScope("DB_HOST", "localhost", "PORT", "5432"))

	entries, err := sm.Get("dev")
	require.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "DB_HOST", entries[0].Key)
}

func TestScopeManager_GetMissing(t *testing.T) {
	sm := NewScopeManager("")
	_, err := sm.Get("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestScopeManager_FallbackToBase(t *testing.T) {
	sm := NewScopeManager("base")
	sm.Add("base", makeScope("DB_HOST", "localhost", "LOG_LEVEL", "info"))
	sm.Add("prod", makeScope("DB_HOST", "prod-db.example.com"))

	entries, err := sm.Get("prod")
	require.NoError(t, err)

	em := make(map[string]string)
	for _, e := range entries {
		em[e.Key] = e.Value
	}

	// prod overrides DB_HOST
	assert.Equal(t, "prod-db.example.com", em["DB_HOST"])
	// LOG_LEVEL falls back to base
	assert.Equal(t, "info", em["LOG_LEVEL"])
}

func TestScopeManager_NoFallbackWhenNoBase(t *testing.T) {
	sm := NewScopeManager("")
	sm.Add("dev", makeScope("APP_ENV", "development"))

	entries, err := sm.Get("dev")
	require.NoError(t, err)
	assert.Len(t, entries, 1)
}

func TestScopeManager_List(t *testing.T) {
	sm := NewScopeManager("base")
	sm.Add("base", makeScope("KEY", "val"))
	sm.Add("staging", makeScope("KEY", "staging-val"))
	sm.Add("prod", makeScope("KEY", "prod-val"))

	names := sm.List()
	assert.Len(t, names, 3)
	assert.ElementsMatch(t, []string{"base", "staging", "prod"}, names)
}

func TestScopeManager_BaseScope_NoFallback(t *testing.T) {
	sm := NewScopeManager("base")
	sm.Add("base", makeScope("X", "1", "Y", "2"))

	entries, err := sm.Get("base")
	require.NoError(t, err)
	assert.Len(t, entries, 2)
}
