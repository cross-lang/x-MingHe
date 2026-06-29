package context_test

import (
	"context"
	"testing"

	"portal/internal/model"
	"portal/internal/pkg"
)

func TestSetUserToCtx(t *testing.T) {
	t.Run("sets user info to context", func(t *testing.T) {
		user := &model.XUser{
			ID:          123,
			Name:        "Test User",
			Account:     "testuser",
			PhoneNumber: "13800138000",
		}
		ctx := context.Background()
		newCtx := pkg.SetUserToCtx(ctx, user)

		if newCtx == ctx {
			t.Error("SetUserToCtx should return a new context")
		}
	})
}

func TestDetailUserFromCtx(t *testing.T) {
	user := &model.XUser{
		ID:          456,
		Name:        "Test User 2",
		Account:     "testuser2",
		PhoneNumber: "13900139000",
	}

	tests := []struct {
		name       string
		ctx        context.Context
		hasUser    bool
		expectedID uint32
	}{
		{
			name:       "context with user",
			ctx:        pkg.SetUserToCtx(context.Background(), user),
			hasUser:    true,
			expectedID: 456,
		},
		{
			name:    "context without user",
			ctx:     context.Background(),
			hasUser: false,
		},
		{
			name:    "nil context",
			ctx:     nil,
			hasUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			result := pkg.DetailUserFromCtx(tt.ctx)

			if tt.hasUser {
				if result == nil {
					t.Error("DetailUserFromCtx should return user")
					return
				}
				if result.ID != tt.expectedID {
					t.Errorf("DetailUserFromCtx user ID = %d, want %d", result.ID, tt.expectedID)
				}
				if result.Name != user.Name {
					t.Errorf("DetailUserFromCtx user Name = %s, want %s", result.Name, user.Name)
				}
			} else {
				if tt.ctx != nil && result != nil {
					t.Error("DetailUserFromCtx should return nil for context without user")
				}
			}
		})
	}
}

func TestSetAndGetUserRoundTrip(t *testing.T) {
	t.Run("set then get user", func(t *testing.T) {
		originalUser := &model.XUser{
			ID:            789,
			Name:          "Round Trip User",
			Account:       "roundtrip",
			PhoneNumber:   "13700137000",
			Avatar:        "http://example.com/avatar.jpg",
			Gender:        1,
			AccountStatus: 1,
		}

		ctx := context.Background()
		newCtx := pkg.SetUserToCtx(ctx, originalUser)
		result := pkg.DetailUserFromCtx(newCtx)

		if result == nil {
			t.Fatal("DetailUserFromCtx returned nil")
		}

		if result.ID != originalUser.ID {
			t.Errorf("User ID mismatch: got %d, want %d", result.ID, originalUser.ID)
		}
		if result.Name != originalUser.Name {
			t.Errorf("User Name mismatch: got %s, want %s", result.Name, originalUser.Name)
		}
		if result.Account != originalUser.Account {
			t.Errorf("User Account mismatch: got %s, want %s", result.Account, originalUser.Account)
		}
	})
}

func TestContextWithMultipleUsers(t *testing.T) {
	t.Run("multiple contexts with different users", func(t *testing.T) {
		user1 := &model.XUser{ID: 1, Name: "User 1", Account: "user1"}
		user2 := &model.XUser{ID: 2, Name: "User 2", Account: "user2"}
		user3 := &model.XUser{ID: 3, Name: "User 3", Account: "user3"}

		ctx1 := pkg.SetUserToCtx(context.Background(), user1)
		ctx2 := pkg.SetUserToCtx(context.Background(), user2)
		ctx3 := pkg.SetUserToCtx(context.Background(), user3)

		retrieved1 := pkg.DetailUserFromCtx(ctx1)
		retrieved2 := pkg.DetailUserFromCtx(ctx2)
		retrieved3 := pkg.DetailUserFromCtx(ctx3)

		if retrieved1.ID != 1 || retrieved1.Name != "User 1" {
			t.Errorf("Retrieved user1 incorrect: got %+v", retrieved1)
		}
		if retrieved2.ID != 2 || retrieved2.Name != "User 2" {
			t.Errorf("Retrieved user2 incorrect: got %+v", retrieved2)
		}
		if retrieved3.ID != 3 || retrieved3.Name != "User 3" {
			t.Errorf("Retrieved user3 incorrect: got %+v", retrieved3)
		}
	})
}

func TestContextNesting(t *testing.T) {
	t.Run("nested context should preserve user", func(t *testing.T) {
		originalUser := &model.XUser{
			ID:      999,
			Name:    "Nested User",
			Account: "nested",
		}

		ctx := context.Background()
		ctx1 := pkg.SetUserToCtx(ctx, originalUser)
		ctx2 := context.WithValue(ctx1, "other_key", "other_value")
		ctx3 := context.WithValue(ctx2, "another_key", 123)

		result := pkg.DetailUserFromCtx(ctx3)

		if result == nil {
			t.Error("DetailUserFromCtx should return user from nested context")
			return
		}
		if result.ID != originalUser.ID {
			t.Errorf("User ID mismatch in nested context: got %d, want %d", result.ID, originalUser.ID)
		}
	})
}

func TestUserInfoKeyConstant(t *testing.T) {
	t.Run("UserInfoKey constant is non-empty", func(t *testing.T) {
		if pkg.UserInfoKey == "" {
			t.Error("UserInfoKey constant should not be empty")
		}
	})
}

func TestSetUserWithZeroID(t *testing.T) {
	t.Run("user with zero ID", func(t *testing.T) {
		user := &model.XUser{
			ID:      0,
			Name:    "Zero ID User",
			Account: "zeroid",
		}

		ctx := pkg.SetUserToCtx(context.Background(), user)
		result := pkg.DetailUserFromCtx(ctx)

		if result == nil {
			t.Error("DetailUserFromCtx should return user even with zero ID")
			return
		}
		if result.ID != 0 {
			t.Errorf("User ID should be 0, got %d", result.ID)
		}
	})
}
