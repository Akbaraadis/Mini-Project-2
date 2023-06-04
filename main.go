package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Actors struct {
	gorm.Model
	Username     string `json:"username"`
	Password     string `json:"password"`
	Role_id      string `json:"role_id"`
	Role_creator string `json:"role_creator"`
	Flag_act     string `json:"flag_act"`
	Flag_ver     string `json:"flag_ver"`
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
	var FlagRole int
	switch {
	case actors.Role_id == "3":
		FlagRole = 0
	case actors.Role_id == "2":
		FlagRole = 1
	}

	var FlagError int

	if FlagRole == 0 {
		err := db.Select("username", "password", "role_id").Create(&actors).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			FlagError = 1
			return
		}
	} else {
		if actors.Role_creator == "2" {
			err := db.Select("username", "password", "role_id").Create(&actors).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				FlagError = 1
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
			FlagError = 1
		}
	}

	// Tampilkan respons berhasil
	if FlagError != 1 {
		c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": actors})
	}
}

func getCustomer(c *gin.Context) {
	var actors []Actors

	username := c.GetHeader("username")

	if username == "superadmin" {
		// Dapatkan semua data user dari database dengan kondisi WHERE flag_ver = nil
		if err := db.Where("flag_ver", nil).Find(&actors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tampilkan data user
		c.JSON(http.StatusOK, gin.H{"users": actors})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Permission Denied"})
	}

}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/customers", createCustomer)
	r.GET("/customers", getCustomer)
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
