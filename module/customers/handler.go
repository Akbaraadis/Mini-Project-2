package customers

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ActorsHandler struct {
	useCase ActorsUseCase
}

func NewActorsHandler(useCase ActorsUseCase) *ActorsHandler {
	return &ActorsHandler{
		useCase: useCase,
	}
}

//func (h *ActorsHandler) CreateCustomer(c *gin.Context) {
//	var actors Actors
//	if err := c.BindJSON(&actors); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	err := h.useCase.CreateCustomer(&actors)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": actors})
//}
//
//func (h *ActorsHandler) CreateAdmin(c *gin.Context) {
//	var actors Actors
//	if err := c.BindJSON(&actors); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	err := h.useCase.CreateAdmin(&actors)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": actors})
//}

func (h *ActorsHandler) LoginAuth(c *gin.Context) {
	username, password, status := c.Request.BasicAuth()
	hash := sha256.New()
	hash.Write([]byte(password))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	password = hashString

	if !status {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tokenKey, err := h.useCase.LoginAuth(username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login has been successful", "token_key": tokenKey})
}
