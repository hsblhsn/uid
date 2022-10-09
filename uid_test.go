package uid_test

import (
	"crypto/rand"
	"sync"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hsblhsn/uid"
)

func BenchmarkNewWithEachEntropy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		entropy := ulid.Monotonic(rand.Reader, 0)
		ulid.MustNew(ulid.Now(), entropy)
	}
}

func BenchmarkNewWithDefaultEntropy(b *testing.B) {
	entropy := ulid.Monotonic(rand.Reader, 0)
	for i := 0; i < b.N; i++ {
		ulid.MustNew(ulid.Now(), entropy)
	}
}

func TestUnique(t *testing.T) {
	t.Parallel()
	const generate = 50000
	wg := sync.WaitGroup{}
	store := make(map[string]struct{})
	mu := sync.Mutex{}
	for i := 0; i < generate; i++ {
		wg.Add(1)
		go func() {
			mu.Lock()
			defer mu.Unlock()
			defer wg.Done()
			id := uid.MustNew("ULID0")
			idStr := id.String()
			if _, ok := store[idStr]; ok {
				t.Errorf("duplicate id: %s", id)
			}
			store[idStr] = struct{}{}
		}()
	}
	wg.Wait()
	assert.Equal(t, generate, len(store))
}

func TestIsValid(t *testing.T) {
	t.Parallel()

	id := uid.MustNew("TEST1")
	isValid := id.IsValid()
	require.True(t, isValid)

	id = uid.ID("INVALID_ID")
	isValid = id.IsValid()
	require.False(t, isValid)
}