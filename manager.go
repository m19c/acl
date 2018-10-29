package acl

import (
	"errors"
	"fmt"
)

type Manager struct {
	registry map[string]*Role
}

func NewManager() *Manager {
	return &Manager{
		registry: map[string]*Role{},
	}
}

func (manager *Manager) Register(roles ...*Role) (*Manager, error) {
	for index, role := range roles {
		if manager.Get(role.Id) != nil {
			return manager, errors.New(fmt.Sprintf("cannot register group %s (on position %d) because the id is already in use", role.Id, index))
		}

		manager.registry[role.Id] = role
	}

	return manager, nil
}

func (manager *Manager) Get(id string) *Role {
	if value, exists := manager.registry[id]; exists {
		return value
	}

	return nil
}

func (manager *Manager) Examine(payload interface{}) []*Role {
	var result []*Role

	for _, role := range manager.registry {
		if role.Examine(payload) {
			result = append(result, role)
		}
	}

	return result
}
