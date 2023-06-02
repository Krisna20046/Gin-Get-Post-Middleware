package main

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var Posts = []Post{
	{ID: 1, Title: "Judul Postingan Pertama", Content: "Ini adalah postingan pertama di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Judul Postingan Kedua", Content: "Ini adalah postingan kedua di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = []User{
	{Username: "user1", Password: "pass1"},
	{Username: "user2", Password: "pass2"},
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Jika tidak ada header Authorization, kirim header WWW-Authenticate untuk meminta otentikasi
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Basic" {
			// Jika format header Authorization tidak sesuai dengan Basic Authentication, kembalikan status Unauthorized
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(authHeaderParts[1])
		if err != nil {
			// Jika terjadi error saat decoding header Authorization, kembalikan status Unauthorized
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		credentials := strings.Split(string(decoded), ":")
		if len(credentials) != 2 {
			// Jika format kredensial tidak sesuai, kembalikan status Unauthorized
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		username := credentials[0]
		password := credentials[1]

		validUser := false
		for _, user := range users {
			if user.Username == username && user.Password == password {
				// Jika username dan password cocok dengan salah satu pengguna yang valid, set validUser menjadi true
				validUser = true
				break
			}
		}

		if !validUser {
			// Jika tidak ada pengguna yang valid, kembalikan status Unauthorized
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()

	} 
}


func SetupRouter() *gin.Engine { //Fungsi SetupRouter() digunakan untuk mengkonfigurasi rute-rute pada mesin Gin.
	r := gin.Default()

	r.Use(authMiddleware()) // Middleware authMiddleware() ditambahkan dengan menggunakan Use(authMiddleware()). Middleware ini akan dijalankan untuk setiap permintaan ke rute-rute yang didefinisikan setelahnya.

	r.GET("/posts", func(c *gin.Context) {
		idStr := c.Query("id")
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka"})
				return
			}

			found := false
			for _, post := range Posts {
				if post.ID == id {
					found = true
					c.JSON(http.StatusOK, gin.H{"post": post})
					break
				}
			}

			if !found {
				c.JSON(http.StatusNotFound, gin.H{"error": "Postingan tidak ditemukan"})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"posts": Posts})
		}
		// TODO: answer here
	})

	r.POST("/posts", func(c *gin.Context) {
		var newPost Post
		if err := c.ShouldBindJSON(&newPost); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		newPost.ID = len(Posts) + 1
		newPost.CreatedAt = time.Now()
		newPost.UpdatedAt = time.Now()
		Posts = append(Posts, newPost)
		c.JSON(http.StatusCreated, gin.H{"message": "Postingan berhasil ditambahkan", "post": newPost})
		// TODO: answer here
	})

	return r
}

func main() {
	r := SetupRouter()

	r.Run(":8080")
}
