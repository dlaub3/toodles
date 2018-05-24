package main

import (
	"net/http"
	"time"

	"github.com/dlaub3/gin-jwt"
	"github.com/dlaub3/toodles/crypt"
	. "github.com/dlaub3/toodles/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func initializeRoutes() {

	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		SendCookie:   true,
		SecureCookie: false,
		Realm:        "test zone",
		Key:          []byte("secret key 12345678910"),
		Timeout:      time.Hour,
		MaxRefresh:   time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			user := User{}
			Mongo.C(CollectionToodlers).Find(bson.M{"email": userId}).One(&user)
			hash := user.Password
			return userId, crypt.CheckPasswordHash(password, hash, 32)
		},
		Authorizator: func(userId string, c *gin.Context) bool {

			user := User{}
			Mongo.C(CollectionToodlers).Find(bson.M{"email": userId}).One(&user)
			csrf(c)
			csrfToken, _ := c.Request.Cookie("csrf")
			c.Keys["csrftoken"] = csrfToken.Value
			c.Keys["uid"] = user.UID

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
		// Get all toodles
		auth.GET("/toodles", getAllToodles)
		// Create a todo
		auth.POST("/toodles", createAToodle)
		// Update a todo
		auth.PUT("/toodles/:toodle_id", updateAToodle)
		// Delete a todo
		auth.DELETE("/toodles/:toodle_id", deleteAToodle)
		// Get a todo by ID
		auth.GET("/toodles/:toodle_id", getAToodle)
		//Method specific to form submitalls
		auth.POST("/toodles/:toodle_id", updateOrDeleteToodle)

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

func csrf(c *gin.Context) {

	_, err := c.Request.Cookie("csrf")

	if err != nil {
		expire := time.Now().UTC().Add(time.Hour)
		maxage := int(expire.Unix() - time.Now().Unix())

		csrf, _ := crypt.GenerateRandomString(32)
		cookie := http.Cookie{
			Name:     "csrf",
			Value:    csrf,
			Path:     "/",
			Expires:  expire,
			MaxAge:   maxage,
			HttpOnly: true,
			Secure:   false, //@dev change when in prod mode
			// No support for SameSite yet https://golang.org/src/net/http/cookie.go
		}
		http.SetCookie(c.Writer, &cookie)
	}
}
