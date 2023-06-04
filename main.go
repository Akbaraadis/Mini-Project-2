package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math"
	"net/http"
	"strconv"
)

type Actors struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Role_id  string `json:"role_id"`
	Flag_act string `json:"flag_act"`
	Flag_ver string `json:"flag_ver"`
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
	var checker Actors

	username := c.GetHeader("username")
	var flag_verified int
	if username == "superadmin" {
		flag_verified = 1
	} else {
		if err := db.Where("username", username).First(&checker).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if checker.Role_id == "2" {
			flag_verified = 2
		}
	}

	if flag_verified == 1 || flag_verified == 2 {
		// Baca data JSON dari body permintaan
		if err := c.BindJSON(&actors); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if actors.Role_id == "3" {
			err := db.Select("username", "password", "role_id", "flag_act").Create(&actors).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't create except Customer Actor"})
			return
		}
	}

	// Tampilkan respons berhasil
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": actors})
}

func createAdmin(c *gin.Context) {
	var actors Actors

	var checker Actors

	username := c.GetHeader("username")
	var flag_verified int
	if username == "superadmin" {
		flag_verified = 1
	} else {
		if err := db.Where("username", username).First(&checker).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if checker.Role_id == "2" {
			flag_verified = 2
		}
	}

	if flag_verified == 2 {
		// Baca data JSON dari body permintaan
		if err := c.BindJSON(&actors); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if actors.Role_id != "2" {
			err := db.Select("username", "password", "role_id", "flag_act").Create(&actors).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
		return
	}

	// Tampilkan respons berhasil
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": actors})

}

func getWaitingApproved(c *gin.Context) {
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

func getCustomer(c *gin.Context) {
	var actors []Actors
	var checker Actors

	username := c.GetHeader("username")
	var flag_verified int
	if username == "superadmin" {
		flag_verified = 1
	} else {
		if err := db.Where("username", username).First(&checker).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if checker.Role_id == "2" {
			flag_verified = 2
		}
	}

	if flag_verified == 1 || flag_verified == 2 {
		// Dapatkan nilai halaman dan ukuran halaman dari query string untuk pagination
		pageStr := c.DefaultQuery("page", "1")
		sizeStr := c.DefaultQuery("size", "10")

		// Konversi nilai halaman dan ukuran halaman ke tipe data yang sesuai
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 1 {
			size = 10
		}

		// Dapatkan semua data admin dari database dengan kondisi WHERE role_id = 2
		var totalRecords int64
		db.Model(&Actors{}).Where("role_id = ?", "2").Count(&totalRecords)

		// Hitung offset berdasarkan halaman dan ukuran halaman
		offset := (page - 1) * size
		if offset < 0 {
			offset = 0
		}

		// Query data admin dengan pagination
		if err := db.Where("role_id = ?", "3").Offset(offset).Limit(size).Find(&actors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tampilkan data admin dan informasi pagination
		c.JSON(http.StatusOK, gin.H{
			"users": actors,
			"page":  page,
			"size":  size,
			"total": totalRecords,
			"pages": int(math.Ceil(float64(totalRecords) / float64(size))),
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
		return
	}

}

func getActorsById(c *gin.Context) {
	var actors Actors
	userID := c.Param("id")

	var checker Actors

	username := c.GetHeader("username")
	var flag_verified int
	if username == "superadmin" {
		flag_verified = 1
	} else {
		if err := db.Where("username", username).First(&checker).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if checker.Role_id == "2" {
			flag_verified = 2
		}
	}

	if flag_verified == 1 || flag_verified == 2 {
		// Dapatkan data user dari database berdasarkan ID
		if err := db.First(&actors, userID).Error; err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tampilkan data user
		c.JSON(http.StatusOK, gin.H{"user": actors})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
	}

}

func getAdmin(c *gin.Context) {
	var actors []Actors
	var checker Actors

	username := c.GetHeader("username")
	var flag_verified int
	if username == "superadmin" {
		flag_verified = 1
	} else {
		if err := db.Where("username", username).First(&checker).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if checker.Role_id == "2" {
			flag_verified = 2
		}
	}

	if flag_verified == 1 || flag_verified == 2 {
		// Dapatkan nilai halaman dan ukuran halaman dari query string untuk pagination
		pageStr := c.DefaultQuery("page", "1")
		sizeStr := c.DefaultQuery("size", "10")

		// Konversi nilai halaman dan ukuran halaman ke tipe data yang sesuai
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 1 {
			size = 10
		}

		// Dapatkan semua data admin dari database dengan kondisi WHERE role_id = 2
		var totalRecords int64
		db.Model(&Actors{}).Where("role_id = ?", "2").Count(&totalRecords)

		// Hitung offset berdasarkan halaman dan ukuran halaman
		offset := (page - 1) * size
		if offset < 0 {
			offset = 0
		}

		// Query data admin dengan pagination
		if err := db.Where("role_id = ?", "2").Offset(offset).Limit(size).Find(&actors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tampilkan data admin dan informasi pagination
		c.JSON(http.StatusOK, gin.H{
			"users": actors,
			"page":  page,
			"size":  size,
			"total": totalRecords,
			"pages": int(math.Ceil(float64(totalRecords) / float64(size))),
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
		return
	}
}

func updateAdmin(c *gin.Context) {
	var actors Actors
	userID := c.Param("id")

	username := c.GetHeader("username")

	if username == "superadmin" {
		// Dapatkan data user dari database berdasarkan ID
		if err := db.First(&actors, userID).Error; err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Baca data JSON dari body permintaan
		if err := c.ShouldBindJSON(&actors); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Simpan perubahan ke database
		if err := db.Select("username", "password", "flag_act", "flag_ver").Save(&actors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tampilkan respons berhasil
		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": actors})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
	}
}

func deleteCustomer(c *gin.Context) {
	var actor Actors
	var deleter Actors

	username := c.GetHeader("username")
	var flag_verified int
	if username == "superadmin" {
		flag_verified = 1
	} else {
		if err := db.Where("username", username).First(&deleter).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if deleter.Role_id == "2" {
			flag_verified = 2
		}
	}

	actorID := c.Param("id")
	// Dapatkan data user dari database berdasarkan ID
	if err := db.First(&actor, actorID).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Hapus data user dari database
	if flag_verified == 1 {
		if err := db.Delete(&actor).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if flag_verified == 2 && actor.Role_id == "3" {
		if err := db.Delete(&actor).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
		return
	}

	// Tampilkan respons berhasil
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/customers", createCustomer)
	r.POST("/admin", createAdmin)
	r.GET("/approved", getWaitingApproved)
	r.GET("/customers", getCustomer)
	r.GET("/customers/:id", getActorsById)
	r.GET("/admin", getAdmin)
	r.GET("/admin:id", getActorsById)
	r.PUT("/admin/:id", updateAdmin)
	r.DELETE("/customers/:id", deleteCustomer)

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
