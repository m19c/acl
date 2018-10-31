package acl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewResultSet(t *testing.T) {
	a := assert.New(t)
	r1 := NewRole("r1")
	r2 := NewRole("r2")

	rs := NewResultSet(r1, r2)
	a.Equal(map[string]*Role{
		"r1": r1,
		"r2": r2,
	}, rs.Matches)
}

func TestResultSet_HasRole(t *testing.T) {
	a := assert.New(t)
	r1 := NewRole("r1")
	rs := NewResultSet(r1)

	a.True(rs.HasRole("r1"))
	a.False(rs.HasRole("r2"))
}

func TestResultSet_GetRole(t *testing.T) {
	a := assert.New(t)
	r1 := NewRole("r1")
	rs := NewResultSet(r1)

	a.Equal(r1, rs.GetRole("r1"))
	a.Nil(rs.GetRole("r2"))
}

func TestResultSet_Has(t *testing.T) {
	a := assert.New(t)
	r1 := NewRole("r1").Grant("a")
	r2 := NewRole("r2").Grant("b")
	rs := NewResultSet(r1, r2)

	a.True(rs.Has("a"))
	a.True(rs.Has("b"))
	a.False(rs.Has("c"))
}

func TestResultSet_HasOneOf(t *testing.T) {
	a := assert.New(t)
	r1 := NewRole("r1").Grant("a")
	r2 := NewRole("r2").Grant("b")
	rs := NewResultSet(r1, r2)

	a.True(rs.HasOneOf("a", "b"))
	a.True(rs.HasOneOf("b", "c"))
	a.False(rs.HasOneOf("c", "d"))
}

func TestResultSet_HasAllOf(t *testing.T) {
	a := assert.New(t)
	r1 := NewRole("r1").Grant("a")
	r2 := NewRole("r2").Grant("b")
	rs := NewResultSet(r1, r2)

	a.True(rs.HasAllOf("a", "b"))
	a.False(rs.HasAllOf("b", "c"))
	a.False(rs.HasAllOf("c", "d"))
}