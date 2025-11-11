package repository

import (
	"fmt"
	"messenger-project/internal/database"
	"messenger-project/internal/models"
)

func CreateUser(user *models.User) error {
	userMap := map[string]any{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"password":   user.Password,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	err := database.Create(database.DB, "users", userMap)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByUsername(username string) (*models.User, error) {
	usersData, err := database.Read(database.DB, "users", []string{"username"}, []any{username})
	if err != nil {
		return nil, err
	}
	if len(usersData) == 0 {
		return nil, fmt.Errorf("No such user: %s", username)
	}
	userData := usersData[0]
	if len(userData) == 0 {
		return nil, fmt.Errorf("No such user: %s", username)
	}
	if userData == nil {
		return nil, fmt.Errorf("No such user: %s", username)
	}
	user := &models.User{
		ID:       userData["id"].(string),
		Username: userData["username"].(string),
		Email:    userData["email"].(string),
		Password: userData["password"].(string),
	}
	return user, nil
}
