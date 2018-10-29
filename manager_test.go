package acl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewManager(t *testing.T) {
	assert.Equal(t, &Manager{registry: map[string]*Role{}}, NewManager())
}

func TestRegister(t *testing.T) {
	manager := NewManager()
	role := NewRole("test")

	res, err := manager.Register(role)

	assert.Equal(t, manager, res)
	assert.NoError(t, err)
}

func TestRegisterSameRole(t *testing.T) {
	manager := NewManager()
	roleA := NewRole("test")
	roleB := NewRole("test")

	res, err := manager.Register(roleA, roleB)

	assert.Error(t, err)
	assert.Equal(t, manager, res)
}

func TestGet(t *testing.T) {
	manager := NewManager()

	assert.Nil(t, manager.Get("test"))

	role := NewRole("test")
	manager.Register(role)
	assert.NotNil(t, manager.Get("test"))
}

type userPayload struct {
	Id    string
	Roles []string
}

func TestExaminer(t *testing.T) {
	manager := NewManager()
	guest := NewRole("guest")
	user := NewRole("user").Grant("user.read").SetExaminer(func(payload interface{}) bool {
		user := payload.(userPayload)

		for _, role := range user.Roles {
			if role == "user" {
				return true
			}
		}

		return false
	})
	admin := NewRole("admin").Grant("user.delete").SetExaminer(func(payload interface{}) bool {
		user := payload.(userPayload)

		for _, role := range user.Roles {
			if role == "admin" {
				return true
			}
		}

		return false
	})
	manager.Register(guest, user, admin)

	granted := manager.Examine(userPayload{
		Id:    "test",
		Roles: []string{"user"},
	})
	assert.Len(t, granted, 1)
	assert.Equal(t, user, granted[0])

	granted = manager.Examine(userPayload{
		Id:    "test",
		Roles: []string{"user", "admin"},
	})
	assert.Len(t, granted, 2)

	var examinedUser *Role
	var examinedAdmin *Role
	for _, examined := range granted {
		if examined.Id == "user" {
			examinedUser = examined
		}

		if examined.Id == "admin" {
			examinedAdmin = examined
		}
	}
	assert.Equal(t, user, examinedUser)
	assert.Equal(t, admin, examinedAdmin)
}
