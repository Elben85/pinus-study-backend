package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func SignUp(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var User struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := c.ShouldBindJSON(&User)
		if err != nil {
			panic(err)
		}

		is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(User.Username)
		
		if !is_alphanumeric {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause": "username must be alphanumeric",
			})
			return
		}

		err2 := database.SignUp(db, User.Email, User.Username, User.Password)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
			//panic(err2)
		}

		//err := database.EditThreadById(db, threadid)
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func LogIn(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var User struct {
			NameOrEmail string `json:"username"`
			Password string `json:"password"`
		}
		err := c.ShouldBindJSON(&User)
		if err != nil {
			panic(err)
		}

		success, token, err2 := database.LogIn(db, User.NameOrEmail, User.Password)
		
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			panic(err2)
		}

		var status string
		if success {
			status = "success"
		} else {
			status = "failure"
			token = ""
		}
		//err := database.EditThreadById(db, threadid)
		c.JSON(http.StatusOK, gin.H{
			"status": status,
			"token":  token,
		})
	}
}
