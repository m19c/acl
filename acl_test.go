package acl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const rightProfileEdit = "profile.edit"
const rightNewsList = "news.list"
const rightNewsCreate = "news.create"
const rightNewsEdit = "news.edit"
const rightNewsDelete = "news.delete"

func createDummyExaminer(returnValue bool) ExaminerFunc {
	return func(payload interface{}) bool {
		return returnValue
	}
}

func TestEverything(t *testing.T) {
	a := assert.New(t)

	user := NewRole("user").Grant(rightProfileEdit)
	editor := NewRole("editor").AcquireFrom(user).Grant(rightNewsList, rightNewsCreate, rightNewsEdit)
	admin := NewRole("admin").AcquireFrom(editor).Grant(rightNewsDelete)

	// register the roles defined above
	manager, err := NewManager().Register(user, editor, admin)
	a.Nil(err)
	a.Equal(user, manager.Get(user.Id))
	a.Equal(editor, manager.Get(editor.Id))
	a.Equal(admin, manager.Get(admin.Id))

	// user: granted rights
	a.True(user.Has(rightProfileEdit))
	a.True(user.HasAllOf(rightProfileEdit))
	a.False(user.Has(rightNewsList))
	a.False(user.Has(rightNewsCreate))
	a.False(user.Has(rightNewsEdit))
	a.False(user.Has(rightNewsDelete))

	// editor: granted rights
	a.True(editor.HasAllOf(rightProfileEdit, rightNewsList, rightNewsCreate, rightNewsEdit))
	a.True(editor.Has(rightProfileEdit))
	a.True(editor.Has(rightNewsList))
	a.True(editor.Has(rightNewsCreate))
	a.True(editor.Has(rightNewsEdit))
	a.False(editor.Has(rightNewsDelete))

	// admin: granted rights
	a.True(admin.HasAllOf(rightProfileEdit, rightNewsList, rightNewsCreate, rightNewsEdit, rightNewsDelete))
	a.True(admin.Has(rightProfileEdit))
	a.True(admin.Has(rightNewsList))
	a.True(admin.Has(rightNewsCreate))
	a.True(admin.Has(rightNewsEdit))
	a.True(admin.Has(rightNewsDelete))
	a.False(admin.HasAllOf(rightProfileEdit, rightNewsList, rightNewsCreate, rightNewsEdit, rightNewsDelete, "test"))

	// finally, test the examiner and the result set
	user.SetExaminer(createDummyExaminer(true))
	admin.SetExaminer(createDummyExaminer(true))
	rs := manager.Examine(1337)

	a.False(rs.HasRole("guest"))
	a.True(rs.HasRole("user"))
	a.True(rs.HasRole("admin"))
}

func TestWorkWithRevokedRights(t *testing.T) {
	guest := NewRole("guest").Grant("register")
	user := NewRole("user").AcquireFrom(guest).Grant("view.profile").Revoke("register")
	admin := NewRole("admin").AcquireFrom(user).Grant("user.delete")

	a := assert.New(t)
	a.True(guest.Has("register"))
	a.True(user.Has("view.profile"))
	a.False(user.Has("register"))
	a.True(admin.Has("user.delete"))
	a.True(admin.Has("view.profile"))
	a.False(admin.Has("register"))
}
