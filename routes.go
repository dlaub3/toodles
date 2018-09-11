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
			err := Mongo.C(CollectionToodlers).Find(bson.M{"email": userId}).One(&user)
			csrf(c)
			if err != nil {
				c.Set("httpStatus", 401)
				c.Set("error", "Your username and password do not match.")
				return userId, false
			}
			hash := user.Password
			if crypt.CheckPasswordHash(password, hash, 32) != true {
				c.Set("httpStatus", 401)
				c.Set("error", "Your username and password do not match.")
				return userId, false
			}
			return userId, true
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			user := User{}
			Mongo.C(CollectionToodlers).Find(bson.M{"email": userId}).One(&user)
			csrfToken, err := c.Request.Cookie("csrf")
			if err != nil {
				csrfToken, err = csrf(c)
			}
			c.Set("csrftoken", csrfToken.Value)
			c.Set("uid", user.UID)
			c.Set("error", "")
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			if c.Request.URL.Path == "/login" {
				showLoginPage(c)
			} else {
				handleUnauthorized(c)
			}
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
		auth.GET("/toodles", getAllToodles)
		auth.POST("/toodles", createAToodle)
		auth.GET("/toodles/:toodle_id", getAToodle)
		auth.PUT("/toodles/:toodle_id", updateAToodle)
		auth.DELETE("/toodles/:toodle_id", deleteAToodle)
		auth.PUT("/toodles/:toodle_id/complete", completeToodle)

		//Routes specifically for form submitalls and not AJAX
		auth.POST("/toodles/:toodle_id", updateOrDeleteToodle)
		auth.POST("/toodles/:toodle_id/complete", completeToodle)

		auth.GET("refresh_token", authMiddleware.RefreshHandler)
	}

	r.GET("/", showHomePage)
	r.GET("/signup", showSignupPage)
	r.POST("/signup", registerNewUser)
	r.GET("/login", showLoginPage)
	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/logout", logout)

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
