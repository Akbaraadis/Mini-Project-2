package customers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ActorsUseCase interface {
	CreateCustomer(c *gin.Context)
}

type actorsUseCase struct {
	actorsRepo ActorsRepository
}

func NewActorsUseCase(actorsRepo ActorsRepository) ActorsUseCase {
	return &actorsUseCase{
		actorsRepo: actorsRepo,
	}
}

func (uc *actorsUseCase) CreateCustomer(c *gin.Context) {
	var actors Actor

	username := c.GetHeader("username")
	var flagVerified int
	if username == "superadmin" {
		flagVerified = 1
	} else {
		checker, err := uc.actorsRepo.FindByUsername(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if checker != nil && checker.RoleID == "2" {
			flagVerified = 2
		}
	}

	if flagVerified == 1 || flagVerified == 2 {
		if err := c.BindJSON(&actors); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if actors.RoleID == "3" {
			if err := uc.actorsRepo.Create(&actors); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't create except Customer Actor"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": actors})
}
