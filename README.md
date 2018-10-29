# ACL

[![Build Status](https://travis-ci.org/MrBoolean/acl.svg?branch=master)](https://travis-ci.org/MrBoolean/acl)

### TL;DR;

> ...

## Example

```go
func main() {
    user := acl.NewRole("user").Grant("profile.edit")
    editor := acl.NewRole("editor").Extend(user).Grant("news.list", "news.create", "news.edit")
    admin := acl.NewRole("admin").Extend(editor).Grant("news.delete")

    manager := NewManager().Register(user, editor, admin)
}
```
