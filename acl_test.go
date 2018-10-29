package acl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const rightProfileEdit = "profile.edit"
const rightNewsList = "news.list"
const rightNewsCreate = "news.create"
const rightNewsEdit = "news.edit"
const rightNewsDelete = "news.delete"

func TestFullIntegration(t *testing.T) {
	a := assert.New(t)

	user := NewRole("user").Grant(rightProfileEdit)
	editor := NewRole("editor").Extend(user).Grant(rightNewsList, rightNewsCreate, rightNewsEdit)
	admin := NewRole("admin").Extend(editor).Grant(rightNewsDelete)

	// register the roles defined above
	manager, err := NewManager().Register(user, editor, admin)
	a.Nil(err)
	a.Equal(user, manager.Get(user.Id))
	a.Equal(editor, manager.Get(editor.Id))
	a.Equal(admin, manager.Get(admin.Id))

	// user: granted rights
	a.True(user.IsAllowed(rightProfileEdit))
	a.True(user.IsAllAllowed(rightProfileEdit))
	a.False(user.IsAllowed(rightNewsList))
	a.False(user.IsAllowed(rightNewsCreate))
	a.False(user.IsAllowed(rightNewsEdit))
	a.False(user.IsAllowed(rightNewsDelete))

	// editor: granted rights
	a.True(editor.IsAllAllowed(rightProfileEdit, rightNewsList, rightNewsCreate, rightNewsEdit))
	a.True(editor.IsAllowed(rightProfileEdit))
	a.True(editor.IsAllowed(rightNewsList))
	a.True(editor.IsAllowed(rightNewsCreate))
	a.True(editor.IsAllowed(rightNewsEdit))
	a.False(editor.IsAllowed(rightNewsDelete))

	// admin: granted rights
	a.True(admin.IsAllAllowed(rightProfileEdit, rightNewsList, rightNewsCreate, rightNewsEdit, rightNewsDelete))
	a.True(admin.IsAllowed(rightProfileEdit))
	a.True(admin.IsAllowed(rightNewsList))
	a.True(admin.IsAllowed(rightNewsCreate))
	a.True(admin.IsAllowed(rightNewsEdit))
	a.True(admin.IsAllowed(rightNewsDelete))
	a.False(admin.IsAllAllowed(rightProfileEdit, rightNewsList, rightNewsCreate, rightNewsEdit, rightNewsDelete, "test"))
}
