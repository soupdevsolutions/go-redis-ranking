package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"strconv"
)

func main() {
	router := gin.Default()
	router.RedirectTrailingSlash = false
	router.GET("/", hello)
	router.GET("/ranks", getRankPage)
	router.GET("/rank", getUser)
	router.POST("/register", register)

	router.Run("127.0.0.1:8080")
}

func hello(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"application": "Ranking System"})
}

func register(c *gin.Context) {
	var newUser userRequest

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	user, err := addUser(newUser.Name)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusCreated, user)
}

func getUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "`id` is required"})
		return
	}

	userData, err := getUserData(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, userData)
}

func getRankPage(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "1"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "`start` must be an integer"})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "`limit` must be an integer"})
		return
	}

	result, err := getRedisRanks(offset, limit)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}
