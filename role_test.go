package acl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRole(t *testing.T) {

	assert.Equal(t, &Role{Id: "test"}, NewRole("test"))
}

func TestRole_AcquireFrom(t *testing.T) {
	a := assert.New(t)

	r1 := NewRole("r1").Grant("a", "b")
	r2 := NewRole("r2").Grant("c", "d")

	a.Equal(r2, r2.AcquireFrom(r1))
	a.Equal(r1, r1.AcquireFrom(r2))

	expected := []string{"a", "b", "c", "d"}
	a.Equal(expected, r1.granted)
	a.Equal(expected, r2.granted)
}

func TestRole_Revoke(t *testing.T) {
	a := assert.New(t)
	role := NewRole("test").Grant("a", "b", "c")

	a.Equal([]string{"a", "b", "c"}, role.granted)
}

func TestRole_Has(t *testing.T) {
	a := assert.New(t)
	role := NewRole("test").Grant("a", "b")

	a.True(role.Has("a"))
	a.True(role.Has("b"))
	a.False(role.Has("c"))
}

func TestRole_HasOneOf(t *testing.T) {
	a := assert.New(t)
	role := NewRole("test").Grant("a", "b")

	a.True(role.HasOneOf("a", "b"))
	a.True(role.HasOneOf("b", "c"))
	a.False(role.HasOneOf("c", "d"))
}

func TestRole_HasAllOf(t *testing.T) {
	a := assert.New(t)
	role := NewRole("test").Grant("a", "b")

	a.True(role.HasAllOf("a", "b"))
	a.False(role.HasAllOf("b", "c"))
	a.False(role.HasAllOf("c", "d"))
}

func TestRole_SetExaminer(t *testing.T) {
	a := assert.New(t)
	role := NewRole("test")

	a.Nil(role.examiner)
	role.SetExaminer(func(payload interface{}) bool { return false })
	a.NotNil(role.examiner)
}

func TestRole_examine(t *testing.T) {
	a := assert.New(t)
	role := NewRole("test")

	type fake struct{}
	a.False(role.examine(fake{}))

	role.SetExaminer(func(payload interface{}) bool { return true })
	a.True(role.examine(fake{}))
}
