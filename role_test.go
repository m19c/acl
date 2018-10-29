package acl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoleNew(t *testing.T) {
	assert.Equal(t, &Role{Id: "test"}, NewRole("test"))
}

func TestExtend(t *testing.T) {
	a := assert.New(t)

	r1 := NewRole("r1")
	r2 := NewRole("r2")

	a.Equal(r2, r2.Extend(r1))
	a.Equal(r1, r1.Extend(r2))
}

func TestRevoke(t *testing.T) {
	a := assert.New(t)
	role := NewRole("r1").Grant("a", "b", "c")

	a.True(role.IsAllowed("c"))
	role.Revoke("c")
	a.False(role.IsAllowed("c"))
}
