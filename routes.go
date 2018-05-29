package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dlaub3/gin-jwt"
	"github.com/dlaub3/toodles/crypt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	mgo "gopkg.in/mgo.v2"
)

// Mongo databse ORM
var Mongo *mgo.Database

const (
	// CollectionToodles contains toodles
	CollectionToodles = "toodles"
	// CollectionToodlers contains toodlers info
	CollectionToodlers = "toodlers"
)

func initializeRoutes() {

	type connection struct {
		Server   string
		Database string
		Mongo    *mgo.Database
	}

	// DB handle
	var connect = connection{}
	var config = Config{}

	config.Read()
	connect.Server = config.Server
	connect.Database = config.Database
	session, err := mgo.Dial(connect.Server)
	if err != nil {
		log.Fatal(err)
	}
	Mongo = session.DB(connect.Database)

	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		SendCookie:   true,
		SecureCookie: false,
		SendRedirect: true,
		RedirectURI:  "/toodles",
		Realm:        "test zone",
		Key:          []byte(config.SecretKey),
		Timeout:      time.Hour,
		MaxRefresh:   time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			user := User{}
			Mongo.C(CollectionToodlers).Find(bson.M{"email": userId}).One(&user)
			hash := user.Password
			csrf(c)
			return userId, crypt.CheckPasswordHash(password, hash, 32)
		},
		Authorizator: func(userId string, c *gin.Context) bool {

			user := User{}
			Mongo.C(CollectionToodlers).Find(bson.M{"email": userId}).One(&user)
			csrfToken, err := c.Request.Cookie("csrf")
			if err != nil {
				csrfToken, err = csrf(c)
			}

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
		//Method specifically for form submitalls and not JSON
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

func csrf(c *gin.Context) (*http.Cookie, error) {

	cookie, err := c.Request.Cookie("csrf")

	if err != nil {
		expire := time.Now().UTC().Add(time.Hour)
		maxage := int(expire.Unix() - time.Now().Unix())

		csrf, err := crypt.GenerateRandomString(32)
		cookie := http.Cookie{
			Name:     "csrf",
			Value:    csrf,
			Path:     "/",
			Expires:  expire,
			MaxAge:   maxage,
			HttpOnly: false, // only access with the secure option
			Secure:   false, //@dev change when in prod mode
			// No support for SameSite yet https://golang.org/src/net/http/cookie.go
		}
		http.SetCookie(c.Writer, &cookie)
		return &cookie, err
	}
	return cookie, err
}
