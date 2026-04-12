package models

type UserRole string

const (
	RoleUser   UserRole = "user"
	RoleMember UserRole = "member"
	RoleAdmin  UserRole = "admin"
)

func (r UserRole) IsValid() bool {
	switch r {
	case RoleUser, RoleMember, RoleAdmin:
		return true
	default:
		return false
	}
}

func (r UserRole) IsMemberOrAbove() bool {
	return r == RoleMember || r == RoleAdmin
}
