package main

import (
	"time"

	jwt "github.com/dlaub3/gin-jwt"
	"github.com/dlaub3/toodles/crypt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func jwtMiddleware() *jwt.GinJWTMiddleware {

	// the JWT middleware
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
			err := mongo.C(collectionToodlers).Find(bson.M{"email": userId}).One(&user)
			csrf(c)

			errors := make(map[string]string)
			errors["Email"] = "Your username and password do not match."

			if err != nil {
				c.Set("httpStatus", 401)
				c.Set("error", errors)
				return userId, false
			}
			hash := user.Password
			if crypt.CheckPasswordHash(password, hash, 32) != true {
				c.Set("httpStatus", 401)
				c.Set("error", errors)
				return userId, false
			}
			return userId, true
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			user := User{}
			if err := mongo.C(collectionToodlers).Find(bson.M{"email": userId}).One(&user); err != nil {
				showErrorPage(c)
				return false
			}
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

	return authMiddleware
}
