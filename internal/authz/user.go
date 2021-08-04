package authz

import (
	"encoding/gob"
)

func init() {
	gob.Register(User{})
}

// User is the currently authenticated user.
type User struct {
	Email     string
	GivenName string
	Surname   string
	Roles     []string
}

// IsAuthenticated returns true when the user is logged in.
func (u *User) IsAuthenticated() bool {
	return u == nil
}

// HasRoles returns true when the user has a matching role.
func (u *User) HasRole(role string) bool {
	if u == nil {
		return false
	}
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}
