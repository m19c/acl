package acl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoleNew(t *testing.T) {
	assert.Equal(t, &Role{Id: "test"}, NewRole("test"))
}

func TestAcquireFrom(t *testing.T) {
	a := assert.New(t)

	r1 := NewRole("r1")
	r2 := NewRole("r2")

	a.Equal(r2, r2.AcquireFrom(r1))
	a.Equal(r1, r1.AcquireFrom(r2))
}

func TestRevoke(t *testing.T) {
	a := assert.New(t)
	role := NewRole("r1").Grant("a", "b", "c")

	a.True(role.Has("c"))
	role.Revoke("c")
	a.False(role.Has("c"))
}
