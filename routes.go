package main

import (
	"time"

	"github.com/dlaub3/gin-jwt"
	. "github.com/dlaub3/toodles/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

// Role for logged in users
var Role string

func initializeRoutes() {

	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		SendCookie:   true,
		SecureCookie: false,
		Realm:        "test zone",
		Key:          []byte("secret key 123"),
		Timeout:      time.Hour,
		MaxRefresh:   time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			user := User{}
			Mongo.C(CollectionToodlers).Find(bson.M{"email": userId}).One(&user)

			if user.Password == password {
				return "", true
			}

			return userId, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {

			Role = "user"
			if userId == "user" {
				Role = "user"
				return true
			} else if userId == "admin" {
				Role = "admin"
				return true
			}

			// @dev
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			render(c, gin.H{
				"title":   "403 Can't touch this.",
				"payload": "403 Can't touch this."}, "error.html")

		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// TokenLookup: "query:jwt",
		// TokenLookup: "query:token",
		TokenLookup: "cookie:token",
		// TokenLookup: "header:Authorization",

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
	r.POST("/signup", registerNewUser)

}
