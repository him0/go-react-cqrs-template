package infrastructure

import (
	"testing"

	"github.com/example/go-vite-spec-kit-sample/internal/domain"
)

func TestInMemoryUserRepository_Save(t *testing.T) {
	repo := NewInMemoryUserRepository()
	user, _ := domain.NewUser("John Doe", "john@example.com")

	err := repo.Save(user)
	if err != nil {
		t.Errorf("Save() unexpected error: %v", err)
	}

	savedUser, err := repo.FindByID(user.ID)
	if err != nil {
		t.Errorf("FindByID() unexpected error: %v", err)
	}

	if savedUser.ID != user.ID {
		t.Errorf("FindByID() ID = %v, want %v", savedUser.ID, user.ID)
	}
}

func TestInMemoryUserRepository_FindByID(t *testing.T) {
	repo := NewInMemoryUserRepository()

	t.Run("user exists", func(t *testing.T) {
		user, _ := domain.NewUser("John Doe", "john@example.com")
		repo.Save(user)

		found, err := repo.FindByID(user.ID)
		if err != nil {
			t.Errorf("FindByID() unexpected error: %v", err)
		}

		if found.ID != user.ID {
			t.Errorf("FindByID() ID = %v, want %v", found.ID, user.ID)
		}
	})

	t.Run("user does not exist", func(t *testing.T) {
		_, err := repo.FindByID("non-existent-id")
		if err == nil {
			t.Error("FindByID() expected error, got nil")
		}
	})
}

func TestInMemoryUserRepository_FindAll(t *testing.T) {
	repo := NewInMemoryUserRepository()

	// Create test users
	for i := 0; i < 5; i++ {
		user, _ := domain.NewUser("User", "user@example.com")
		repo.Save(user)
	}

	users, err := repo.FindAll(10, 0)
	if err != nil {
		t.Errorf("FindAll() unexpected error: %v", err)
	}

	if len(users) != 5 {
		t.Errorf("FindAll() returned %d users, want 5", len(users))
	}
}

func TestInMemoryUserRepository_Count(t *testing.T) {
	repo := NewInMemoryUserRepository()

	// Initially empty
	count, err := repo.Count()
	if err != nil {
		t.Errorf("Count() unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("Count() = %d, want 0", count)
	}

	// Add users
	for i := 0; i < 3; i++ {
		user, _ := domain.NewUser("User", "user@example.com")
		repo.Save(user)
	}

	count, err = repo.Count()
	if err != nil {
		t.Errorf("Count() unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("Count() = %d, want 3", count)
	}
}

func TestInMemoryUserRepository_Delete(t *testing.T) {
	repo := NewInMemoryUserRepository()

	t.Run("delete existing user", func(t *testing.T) {
		user, _ := domain.NewUser("John Doe", "john@example.com")
		repo.Save(user)

		err := repo.Delete(user.ID)
		if err != nil {
			t.Errorf("Delete() unexpected error: %v", err)
		}

		_, err = repo.FindByID(user.ID)
		if err == nil {
			t.Error("FindByID() should return error after delete")
		}
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		err := repo.Delete("non-existent-id")
		if err == nil {
			t.Error("Delete() expected error, got nil")
		}
	})
}
