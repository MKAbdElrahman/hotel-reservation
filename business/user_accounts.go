package business

import (
	"context"

	"github.com/mkabdelrahman/hotel-reservation/auth"
	"github.com/mkabdelrahman/hotel-reservation/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *Manager) GetUserToken(ctx context.Context, authParams AuthParams) (string, error) {

	user, err := m.UserStore.GetUserByEmail(ctx, authParams.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(authParams.Password)); err != nil {
		return "", err
	}

	ok, err := types.IsValidPassword(user.EncryptedPassword, authParams.Password)

	if err != nil || !ok {
		return "", err
	}

	token, err := auth.GenerateAuthToken(user.ID.Hex())

	if err != nil {
		return "", err
	}

	return token, nil
}

// func (m *Manager) ChangePassword(ctx context.Context, userID, newPassword string) error {
// 	return nil
// }

// func (m *Manager) SuspendUserAccount(ctx context.Context, userID string) error {
// 	return nil
// }

// func (m *Manager) ActivateUserAccount(ctx context.Context, userID string) error {

// 	return nil
// }
