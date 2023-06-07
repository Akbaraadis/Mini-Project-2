package customers

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type ActorsUseCase interface {
	//CreateCustomer(actors *Actors) error
	//CreateAdmin(actors *Actors) error
	LoginAuth(username, password string) (string, error)
}

type actorsUseCase struct {
	repo ActorsRepository
}

func NewActorsUseCase(repo ActorsRepository) ActorsUseCase {
	return &actorsUseCase{
		repo: repo,
	}
}

//func (uc *actorsUseCase) CreateCustomer(actors *Actors) error {
//	return uc.repo.CreateActor(actors)
//}
//
//func (uc *actorsUseCase) CreateAdmin(actors *Actors) error {
//	return uc.repo.CreateActor(actors)
//}
//
//func (uc *actorsUseCase) CheckRole(token string) (*Actors, error) {
//	return uc.repo.GetActorByToken(token)
//}

func (uc *actorsUseCase) LoginAuth(username, password string) (string, error) {
	actor, err := uc.repo.GetActorByUsername(username)
	if err != nil {
		return "", err
	}

	if actor.Password != password {
		return "", errors.New("Wrong Password")
	}

	claims := jwt.MapClaims{
		"sub":  actor.ID,
		"name": actor.Username,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		return "", err
	}

	err = uc.repo.UpdateTokenKey(actor, signedToken)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
