package acl

import (
	"sort"
	"strings"
)

type ExaminerFunc = func(payload interface{}) bool

type Role struct {
	Id       string
	extends  []*Role
	rights   []string
	examiner ExaminerFunc
}

func NewRole(id string) *Role {
	return &Role{
		Id:     id,
	}
}

func (role *Role) Extend(roles ...*Role) *Role {
	role.extends = append(role.extends, roles...)
	return role
}

func (role *Role) determineInternalRightIndex(right string) int {
	index := sort.SearchStrings(role.rights, right)

	if index < len(role.rights) && role.rights[index] == right {
		return index
	}

	return -1
}

func (role *Role) Grant(rights ...string) *Role {
	var granted []string

	for _, right := range rights {
		right = strings.ToLower(right)

		if role.determineInternalRightIndex(right) == -1 {
			granted = append(granted, right)
		}
	}

	role.rights = append(role.rights, granted...)

	// to use sort.SearchStrings, the slice must be sorted in ascending order
	sort.Strings(role.rights)

	return role
}

func (role *Role) Revoke(rights ...string) *Role {
	for _, right := range rights {
		index := role.determineInternalRightIndex(right)

		if index != -1 {
			role.rights = append(role.rights[:index], role.rights[index+1:]...)
		}
	}

	return role
}

func (role *Role) IsAllowed(rights ...string) bool {
	for _, right := range rights {
		if role.determineInternalRightIndex(right) != -1 {
			return true
		}
	}

	for _, extended := range role.extends {
		if extended.IsAllowed(rights...) {
			return true
		}
	}

	return false
}

func (role *Role) IsAllAllowed(rights ...string) bool {
	registry := map[string]bool{}
	total := len(rights)

	resolve := func(right string) bool {
		if !registry[right] {
			total = total - 1
		}

		registry[right] = true

		if total == 0 {
			return true
		}

		return false
	}

	for _, right := range rights {
		if role.determineInternalRightIndex(right) != -1 && resolve(right) {
			return true
		}
	}

	for _, extended := range role.extends {
		for _, right := range rights {
			if extended.IsAllowed(right) && resolve(right) {
				return true
			}
		}
	}

	return false
}

func (role *Role) SetExaminer(examiner ExaminerFunc) *Role {
	role.examiner = examiner
	return role
}

func (role *Role) Examine(payload interface{}) bool {
	if role.examiner == nil {
		return false
	}

	return role.examiner(payload)
}
