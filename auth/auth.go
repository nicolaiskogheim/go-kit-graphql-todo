package auth

// Identifier Identifies the authenticated entity
type Identifier string

func (id Identifier) ToString() string {
	return string(id)
}
