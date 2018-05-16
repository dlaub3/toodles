package main

import (
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func initializeRoutes() {

	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key 123"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if userId == "admin" && password == "admin" {
				return userId, true
			}

			return userId, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			if userId == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "query:jwt",
		// TokenLookup: "query:token",
		//TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}

	auth := r.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// Get all todos
		auth.GET("/todos", getAllTodos)
		// Create a todo
		auth.POST("/todos", createATodo)
		// Update a todo
		auth.PUT("/todos/:todo_id", updateATodo)
		// Delete a todo
		auth.DELETE("/todos/:todo_id", deleteATodo)
		// Get a todo by ID
		auth.GET("/todos/:todo_id", getATodo)
		//Method specific to form submitalls
		auth.POST("/todos/:todo_id", updateOrDeleteTodo)

		auth.GET("refresh_token", authMiddleware.RefreshHandler)
	}

	// Handle the index route
	r.GET("/", showHomePage)
	// Handle the login route
	r.GET("/login", showLoginPage)
	r.POST("/login", authMiddleware.LoginHandler)
	// Handle the login route
	r.GET("/signup", showSignupPage)

}
