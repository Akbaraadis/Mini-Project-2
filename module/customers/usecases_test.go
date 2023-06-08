package customers_test

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"miniproject2/module/customers"
	"testing"
	"time"
)

type MockActorsRepository struct {
	GetActorByUsernameFunc func(username string) (*customers.Actors, error)
	UpdateTokenKeyFunc     func(actor *customers.Actors, token string) error
}

func (m *MockActorsRepository) GetActorByUsername(username string) (*customers.Actors, error) {
	return m.GetActorByUsernameFunc(username)
}

func (m *MockActorsRepository) UpdateTokenKey(actor *customers.Actors, token string) error {
	return m.UpdateTokenKeyFunc(actor, token)
}

func TestActorsUseCase_LoginAuth(t *testing.T) {
	repo := &MockActorsRepository{}
	uc := customers.NewActorsUseCase(repo)

	// Situasi sukses
	repo.GetActorByUsernameFunc = func(username string) (*customers.Actors, error) {
		return &customers.Actors{
			Username: username,
			Password: "password123",
		}, nil
	}
	repo.UpdateTokenKeyFunc = func(actor *customers.Actors, token string) error {
		// Implementasi simulasi update token key
		return nil
	}

	signedToken, err := uc.LoginAuth("superadmin", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, signedToken)

	// Situasi gagal - actor tidak ditemukan
	repo.GetActorByUsernameFunc = func(username string) (*customers.Actors, error) {
		return nil, errors.New("actor not found")
	}

	signedToken, err = uc.LoginAuth("admin", "password123")

	assert.Error(t, err)
	assert.Empty(t, signedToken)

	// Situasi gagal - password salah
	repo.GetActorByUsernameFunc = func(username string) (*customers.Actors, error) {
		return &customers.Actors{
			Username: username,
			Password: "pass123",
		}, nil
	}

	signedToken, err = uc.LoginAuth("user1", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, signedToken)
}

func TestActorsUseCase_LoginAuth_TokenExpiration(t *testing.T) {
	repo := &MockActorsRepository{}
	uc := customers.NewActorsUseCase(repo)

	// Situasi sukses
	actor := &customers.Actors{
		Username: "user1",
		Password: "pass123",
	}
	repo.GetActorByUsernameFunc = func(username string) (*customers.Actors, error) {
		return actor, nil
	}
	repo.UpdateTokenKeyFunc = func(actor *customers.Actors, token string) error {
		// Implementasi simulasi update token key
		return nil
	}

	signedToken, err := uc.LoginAuth("user1", "pass123")

	assert.NoError(t, err)
	assert.NotEmpty(t, signedToken)

	// Verifikasi token
	token, _ := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)

	//assert.Equal(t, actor.ID, claims["sub"])
	assert.Equal(t, actor.Username, claims["name"])

	expirationTime := int64(claims["exp"].(float64))
	expectedExpirationTime := time.Now().Add(time.Hour * 1).Unix()
	assert.InDelta(t, expectedExpirationTime, expirationTime, 1)
}
