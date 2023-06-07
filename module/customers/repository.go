package customers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
)

type ActorsRepository interface {
	//CreateActor(actors *Actors) error
	GetActorByUsername(username string) (*Actors, error)
	//GetActorByToken(username string) (*Actors, error)
	UpdateTokenKey(actor *Actors, tokenKey string) error
}

type actorsRepository struct {
	db *gorm.DB
}

func NewActorsRepository(db *gorm.DB) ActorsRepository {
	return &actorsRepository{
		db: db,
	}
}

func (repo *actorsRepository) CreateActor(actors *Actors) error {
	hash := sha256.New()
	hash.Write([]byte(actors.Password))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	actors.Password = hashString

	if actors.RoleID == "3" {
		return repo.db.Select("username", "password", "role_id", "flag_act").Create(actors).Error
	}

	return errors.New("Can't create actor except Customer Actor")
}

func (repo *actorsRepository) GetActorByUsername(username string) (*Actors, error) {
	var actor Actors
	err := repo.db.Where("username = ?", username).First(&actor).Error
	if err != nil {
		return nil, err
	}
	return &actor, nil
}

//func (repo *actorsRepository) GetActorByToken(token string) (*Actors, error) {
//	var actor Actors
//	err := repo.db.Where("token_key = ?", token).First(&actor).Error
//	if err != nil {
//		return nil, err
//	}
//	return &actor, nil
//}

func (repo *actorsRepository) UpdateTokenKey(actor *Actors, tokenKey string) error {
	actor.TokenKey = tokenKey
	return repo.db.Model(actor).Update("token_key", tokenKey).Error
}
