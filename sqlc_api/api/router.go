package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func SetupRoutes(router *gin.Engine) {
	router.GET("/user", homePage)
}

func homePage(c *gin.Context) {
	var user User = User{
		ID:    1,
		Name:  "John Doe",
		Email: "test@test.com",
	}
	c.JSON(http.StatusOK, user)
}
