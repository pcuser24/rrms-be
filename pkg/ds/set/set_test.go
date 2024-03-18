package set

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPrimitiveTypes(t *testing.T) {
	s := NewSet[string]()
	s.Add("abc123")
	s.Add("def456")
	require.True(t, s.Contains("abc123"))
	require.True(t, s.Contains("def456"))
	require.False(t, s.Contains("xyz789"))
	s.Remove("abc123")
	require.False(t, s.Contains("abc123"))
	require.True(t, s.Contains("def456"))
}

func TestUUID(t *testing.T) {
	s := NewSet[uuid.UUID]()
	u1 := uuid.MustParse("978fe220-663f-464c-9afd-fe05de7be44b")
	u2 := uuid.MustParse("978fe220-663f-464c-9afd-fe05de7be44b")
	s.Add(u1)
	require.True(t, s.Contains(u1))
	require.True(t, s.Contains(u2))
	u1 = uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f")
	require.False(t, s.Contains(u1))
}

func TestStructType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	s := NewSet[Person]()
	s.Add(Person{"Alice", 25})
	s.Add(Person{"Bob", 30})
	require.True(t, s.Contains(Person{"Alice", 25}))
	require.True(t, s.Contains(Person{"Bob", 30}))
	require.False(t, s.Contains(Person{"Alice", 35}))
	s.Remove(Person{"Alice", 25})
	require.False(t, s.Contains(Person{"Alice", 25}))
	require.True(t, s.Contains(Person{"Bob", 30}))
}

func TestPointerToStructType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	s := NewSet[*Person]()
	alice := Person{"Alice", 25}
	bob := Person{"Bob", 30}
	_alice := &alice
	_alice.Age = 35
	s.Add(&alice)
	s.Add(&bob)
	require.True(t, s.Contains(&alice))
	require.True(t, s.Contains(&bob))
	require.True(t, s.Contains(_alice))
	s.Remove(_alice)
	require.False(t, s.Contains(&alice))
	require.True(t, s.Contains(&bob))
}
