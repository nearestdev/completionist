package models

import "testing"

func TestUserRoleIsValid(t *testing.T) {
	tests := []struct {
		name string
		role UserRole
		want bool
	}{
		{name: "user", role: RoleUser, want: true},
		{name: "member", role: RoleMember, want: true},
		{name: "admin", role: RoleAdmin, want: true},
		{name: "unknown role", role: UserRole("superadmin"), want: false},
		{name: "empty", role: UserRole(""), want: false},
		{name: "wrong case", role: UserRole("Admin"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRoleIsMemberOrAbove(t *testing.T) {
	tests := []struct {
		name string
		role UserRole
		want bool
	}{
		{name: "member", role: RoleMember, want: true},
		{name: "admin", role: RoleAdmin, want: true},
		{name: "user", role: RoleUser, want: false},
		{name: "invalid", role: UserRole("superadmin"), want: false},
		{name: "empty", role: UserRole(""), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsMemberOrAbove(); got != tt.want {
				t.Errorf("IsMemberOrAbove() = %v, want %v", got, tt.want)
			}
		})
	}
}
