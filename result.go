package acl

type ResultSet struct {
	Matches map[string]*Role
}

// NewResultSet creates a new `ResultSet` instance
func NewResultSet(roles ...*Role) *ResultSet {
	matches := map[string]*Role{}

	for _, role := range roles {
		matches[role.Id] = role
	}

	return &ResultSet{
		Matches: matches,
	}
}

// HasRole determines whether a role is present.
func (result *ResultSet) HasRole(id string) bool {
	if _, exists := result.Matches[id]; exists {
		return true
	}

	return false
}

// GetRole returns the role with the given identifier or nil.
func (result *ResultSet) GetRole(id string) *Role {
	if result.HasRole(id) {
		return result.Matches[id]
	}

	return nil
}

// Has checks that at least one role contains the given right.
func (result *ResultSet) Has(right string) bool {
	for _, role := range result.Matches {
		if role.Has(right) {
			return true
		}
	}

	return false
}

// HasOneOf checks that at least one role contains at least one of the given rights.
func (result *ResultSet) HasOneOf(rights ...string) bool {
	for _, right := range rights {
		for _, role := range result.Matches {
			if role.Has(right) {
				return true
			}
		}
	}

	return false
}

// HasAllOf verifies that all specified rights are present.
func (result *ResultSet) HasAllOf(rights ...string) bool {
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
		for _, role := range result.Matches {
			if role.Has(right) && resolve(right) {
				return true
			}
		}
	}

	return false
}
