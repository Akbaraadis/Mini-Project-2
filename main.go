package main

import (
	"github.com/gin-gonic/gin"
	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Actors struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Role_id  string `json:"role_id"`
}

var db *gorm.DB
var err error

func initDB() {
	dsn := "root:mysql@tcp(localhost:3306)/milestone1?parseTime=true"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

func createCustomer(c *gin.Context) {
	var actors Actors

	// Baca data JSON dari body permintaan
	if err := c.BindJSON(&actors); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan data ke database
	err := db.Create(&actors).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Tampilkan respons berhasil
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": actors})
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/customers", createCustomer)
	//r.GET("/users", getUsers)
	//r.GET("/users/:id", getUserById)
	//r.PUT("/users/:id", updateUser)
	//r.DELETE("/users/:id", deleteUser)

	return r
}

func main() {
	initDB()
	r := setupRouter()

	// Jalankan server di port 8080
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
