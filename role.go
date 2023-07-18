// Copyright 2018 Marc Binder. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package acl

import (
	"sort"
)

// ExaminerFunc a function to determine whether a role can be added to a `ResultSet`.
type ExaminerFunc = func(payload interface{}) bool

// NewRole returns a new role instance.
func NewRole(id string) *Role {
	return &Role{
		Id: id,
	}
}

// Role contains all necessary information about one set of granted rights.
//
// Each role requires an identifier. It is possible to define multiple roles with the
// same identifier as long as each manager contains an unique set of identifiers.
type Role struct {
	Id       string
	granted  []string
	examiner ExaminerFunc
}

// Grant adds the given right(s) to the role.
//
//		   func main() {
//	        r := NewRole("a")
//	        r.Grant("right.a", "right.b")
//	    }
//
// Note, that duplications will be ignored.
func (role *Role) Grant(rights ...string) *Role {
	var granted []string

	for _, right := range rights {
		if determineIndex(role.granted, right) == -1 {
			granted = append(granted, right)
		}
	}

	if len(granted) > 0 {
		role.granted = append(role.granted, granted...)

		// to use sort.SearchStrings, the slice must be sorted in ascending order
		sort.Strings(role.granted)
	}

	return role
}

// Revoke removes the given right(s) from the role.
//
//		   func main() {
//	        r := NewRole("a")
//	        r.Grant("right.a", "right.b")
//	        r.Revoke("right.a")
//	    }
func (role *Role) Revoke(rights ...string) *Role {
	for _, right := range rights {
		if index := determineIndex(role.granted, right); index >= 0 {
			role.granted = append(role.granted[:index], role.granted[index+1:]...)
		}
	}

	return role
}

// AcquireFrom grabs the rights from the given roles to add them.
//
//		   func main() {
//	        r1 := NewRole("r1").Grant("right.a")
//	        r2 := NewRole("r2").AcquireFrom(r1).Grant("right.b")
//	    }
func (role *Role) AcquireFrom(roles ...*Role) *Role {
	for _, ar := range roles {
		role.Grant(ar.granted...)
	}

	return role
}

// Has checks that the given right has been granted.
//
// To resolve whether a right is available or not, the function uses a binary search to
// determine the actual index of the given right(s). Therefore, the array of granted
// rights is always sorted alphabetically.
func (role *Role) Has(right string) bool {
	return determineIndex(role.granted, right) >= 0
}

// HasOneOf checks that at least one of the given rights has been granted.
//
//	func main() {
//	    r := NewRole("r")
//	    r.Grant("right.a")
//	    r.HasOneOf("right.a", "right.b")
//	}
func (role *Role) HasOneOf(rights ...string) bool {
	for _, right := range rights {
		if determineIndex(role.granted, right) >= 0 {
			return true
		}
	}

	return false
}

// HasAllOf verifies that all specified rights are present.
//
//	func main() {
//	    r := NewRole("r").Grant("a", "b", "c")
//	    r.HasAllOf("a", b", "c")
//	}
func (role *Role) HasAllOf(rights ...string) bool {
	registry := map[string]bool{}
	total := len(rights)

	resolve := func(right string) bool {
		if !registry[right] {
			total--
		}

		registry[right] = true

		return total == 0
	}

	for _, right := range rights {
		if role.Has(right) && resolve(right) {
			return true
		}
	}

	return false
}

// SetExaminer sets / overwrites the examiner.
//
// The examiner is used to determine whether a role can be added to a `ResultSet`.
//
//	type User struct {
//	    isAdmin bool
//	}
//
//	func main() {
//	    r := NewRole("admin").Grant("godmode").SetExaminer(func (payload interface{}) bool {
//	        user := payload.(User)
//	        return user.isAdmin
//	    })
//
//	    rs := NewManager().Register(r).Examine(User{isAdmin: true})
//	}
func (role *Role) SetExaminer(examiner ExaminerFunc) *Role {
	role.examiner = examiner
	return role
}

// examine calls the examiner function to determine if a role can be added to a `ResultSet`.
//
// Without examiner, the function will always return false.
func (role *Role) examine(payload interface{}) bool {
	if role.examiner == nil {
		return false
	}

	return role.examiner(payload)
}
