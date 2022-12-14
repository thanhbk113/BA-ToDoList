package middleware

import (
	"log"
	"net/http"
	"todolist/helpers"

	"github.com/gin-gonic/gin"
)

// authMiddleware.go is a middleware that checks if the user is logged in or not
// authorizes validates token and authorizes user to access the routes

func Authentcation() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get token from header
		token := c.GetHeader("Authorization")
		log.Fatal("Token get from Header:", token)
		// check token
		if token == "" {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort() //abort request if token is not valid
			return
		}

		//check token is valid or not
		claims, err := helpers.ValidateToken(token)
		if err != "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Set("user_id", claims.User_id) //set user_id to context

		c.Next() //proceed in the middleware chain
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
