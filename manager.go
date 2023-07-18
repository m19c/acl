// Copyright 2018 Marc Binder. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package acl

import (
	"fmt"
)

// Manager contains all registered roles.
type Manager struct {
	registry map[string]*Role
}

// NewManager creates a new manager instance.
func NewManager() *Manager {
	return &Manager{
		registry: map[string]*Role{},
	}
}

// Register transfers the given roles into the Manager registry.
//
//	func main() {
//	    m := NewManager().Register(NewRole("a").Grant("b"), NewRole("c").Grant("d"))
//	}
//
// Note, that each role MUST contain an unique identifier.
func (manager *Manager) Register(roles ...*Role) (*Manager, error) {
	for index, role := range roles {
		if manager.Get(role.Id) != nil {
			return manager, fmt.Errorf("cannot register role %s (on position %d) because the id is already in use", role.Id, index)
		}

		manager.registry[role.Id] = role
	}

	return manager, nil
}

// Ensure returns or creates a role with the given id.
func (manager *Manager) Ensure(id string) *Role {
	if value := manager.Get(id); value != nil {
		return value
	}

	role := NewRole(id)
	manager.registry[id] = role

	return role
}

// Get returns the role with the given id.
func (manager *Manager) Get(id string) *Role {
	if value, exists := manager.registry[id]; exists {
		return value
	}

	return nil
}

// Examine starts the examining process to determine the available roles.
func (manager *Manager) Examine(payload interface{}) *ResultSet {
	var matched []*Role

	for _, role := range manager.registry {
		if role.examine(payload) {
			matched = append(matched, role)
		}
	}

	return NewResultSet(matched...)
}
