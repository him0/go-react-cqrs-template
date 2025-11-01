package application

import (
	"testing"

	"github.com/example/go-vite-spec-kit-sample/internal/infrastructure"
)

func TestUserService_CreateUser(t *testing.T) {
	repo := infrastructure.NewInMemoryUserRepository()
	service := NewUserService(repo)

	user, err := service.CreateUser("John Doe", "john@example.com")
	if err != nil {
		t.Errorf("CreateUser() unexpected error: %v", err)
	}

	if user.Name != "John Doe" {
		t.Errorf("CreateUser() name = %v, want John Doe", user.Name)
	}

	if user.Email != "john@example.com" {
		t.Errorf("CreateUser() email = %v, want john@example.com", user.Email)
	}

	// Verify user was saved
	savedUser, err := repo.FindByID(user.ID)
	if err != nil {
		t.Errorf("FindByID() unexpected error: %v", err)
	}

	if savedUser.ID != user.ID {
		t.Errorf("FindByID() ID = %v, want %v", savedUser.ID, user.ID)
	}
}

func TestUserService_GetUser(t *testing.T) {
	repo := infrastructure.NewInMemoryUserRepository()
	service := NewUserService(repo)

	createdUser, _ := service.CreateUser("John Doe", "john@example.com")

	t.Run("user exists", func(t *testing.T) {
		user, err := service.GetUser(createdUser.ID)
		if err != nil {
			t.Errorf("GetUser() unexpected error: %v", err)
		}

		if user.ID != createdUser.ID {
			t.Errorf("GetUser() ID = %v, want %v", user.ID, createdUser.ID)
		}
	})

	t.Run("user does not exist", func(t *testing.T) {
		_, err := service.GetUser("non-existent-id")
		if err == nil {
			t.Error("GetUser() expected error, got nil")
		}
	})
}

func TestUserService_ListUsers(t *testing.T) {
	repo := infrastructure.NewInMemoryUserRepository()
	service := NewUserService(repo)

	// Create test users
	for i := 0; i < 5; i++ {
		service.CreateUser("User", "user@example.com")
	}

	users, total, err := service.ListUsers(10, 0)
	if err != nil {
		t.Errorf("ListUsers() unexpected error: %v", err)
	}

	if len(users) != 5 {
		t.Errorf("ListUsers() returned %d users, want 5", len(users))
	}

	if total != 5 {
		t.Errorf("ListUsers() total = %d, want 5", total)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	repo := infrastructure.NewInMemoryUserRepository()
	service := NewUserService(repo)

	createdUser, _ := service.CreateUser("John Doe", "john@example.com")

	t.Run("update existing user", func(t *testing.T) {
		updatedUser, err := service.UpdateUser(createdUser.ID, "Jane Doe", "jane@example.com")
		if err != nil {
			t.Errorf("UpdateUser() unexpected error: %v", err)
		}

		if updatedUser.Name != "Jane Doe" {
			t.Errorf("UpdateUser() name = %v, want Jane Doe", updatedUser.Name)
		}

		if updatedUser.Email != "jane@example.com" {
			t.Errorf("UpdateUser() email = %v, want jane@example.com", updatedUser.Email)
		}
	})

	t.Run("update non-existent user", func(t *testing.T) {
		_, err := service.UpdateUser("non-existent-id", "Jane Doe", "jane@example.com")
		if err == nil {
			t.Error("UpdateUser() expected error, got nil")
		}
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	repo := infrastructure.NewInMemoryUserRepository()
	service := NewUserService(repo)

	createdUser, _ := service.CreateUser("John Doe", "john@example.com")

	t.Run("delete existing user", func(t *testing.T) {
		err := service.DeleteUser(createdUser.ID)
		if err != nil {
			t.Errorf("DeleteUser() unexpected error: %v", err)
		}

		// Verify user was deleted
		_, err = repo.FindByID(createdUser.ID)
		if err == nil {
			t.Error("FindByID() should return error after delete")
		}
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		err := service.DeleteUser("non-existent-id")
		if err == nil {
			t.Error("DeleteUser() expected error, got nil")
		}
	})
}
