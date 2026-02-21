package domain

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		email    string
		wantErr  bool
	}{
		{
			name:     "valid user",
			userName: "John Doe",
			email:    "john@example.com",
			wantErr:  false,
		},
		{
			name:     "empty name",
			userName: "",
			email:    "john@example.com",
			wantErr:  true,
		},
		{
			name:     "empty email",
			userName: "John Doe",
			email:    "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.userName, tt.email)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error: %v", err)
				return
			}

			if user.Name != tt.userName {
				t.Errorf("NewUser() name = %v, want %v", user.Name, tt.userName)
			}

			if user.Email != tt.email {
				t.Errorf("NewUser() email = %v, want %v", user.Email, tt.email)
			}

			if user.ID == "" {
				t.Error("NewUser() ID should not be empty")
			}

			if user.CreatedAt.IsZero() {
				t.Error("NewUser() CreatedAt should not be zero")
			}

			if user.UpdatedAt.IsZero() {
				t.Error("NewUser() UpdatedAt should not be zero")
			}
		})
	}
}

func TestUser_Update(t *testing.T) {
	user, err := NewUser("John Doe", "john@example.com")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	originalUpdatedAt := user.UpdatedAt

	tests := []struct {
		name     string
		newName  string
		newEmail string
	}{
		{
			name:     "update name and email",
			newName:  "Jane Doe",
			newEmail: "jane@example.com",
		},
		{
			name:     "update name only",
			newName:  "John Smith",
			newEmail: "",
		},
		{
			name:     "update email only",
			newName:  "",
			newEmail: "smith@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.Update(tt.newName, tt.newEmail)
			if err != nil {
				t.Errorf("Update() unexpected error: %v", err)
			}

			if tt.newName != "" && user.Name != tt.newName {
				t.Errorf("Update() name = %v, want %v", user.Name, tt.newName)
			}

			if tt.newEmail != "" && user.Email != tt.newEmail {
				t.Errorf("Update() email = %v, want %v", user.Email, tt.newEmail)
			}

			if !user.UpdatedAt.After(originalUpdatedAt) {
				t.Error("Update() UpdatedAt should be updated")
			}
		})
	}
}
