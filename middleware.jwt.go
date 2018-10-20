package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/dlaub3/toodles/crypt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func jwtMiddleware() *jwt.GinJWTMiddleware {

	type login struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	identityKey := "uid"

	// the JWT middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		SigningAlgorithm: "HS256",
		SendCookie:       true,
		SecureCookie:     false, //non HTTPS dev environments
		CookieHTTPOnly:   true,
		CookieDomain:     "localhost:8080",
		CookieName:       "token",
		Realm:            "test zone",
		Key:              []byte(config.SecretKey),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		IdentityKey:      identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// process JWT JWT and map to claims
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UID,
				}
			}
			return jwt.MapClaims{"error": "no claims to map"}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			// extract UID from claims UID identityKey
			claims := jwt.ExtractClaims(c)
			return claims["uid"].(string)
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// handle authorization for data set in IdentityHandler
			if v, ok := data.(string); ok && v != "" {
				return true
			}
			showErrorPage(c)
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			fmt.Println("unauth")
			if c.Request.URL.Path == "/login" {
				showLoginPage(c)
			} else {
				handleUnauthorized(c)
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			fmt.Println("auth")
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			userID := loginVals.Username
			password := loginVals.Password

			user := User{}
			err := mongo.C(collectionToodlers).Find(bson.M{"email": userID}).One(&user)
			csrf(c)

			errors := make(map[string]string)
			errors["Email"] = "Your username and password do not match."

			if err != nil {
				c.Set("httpStatus", http.StatusUnauthorized)
				c.Set("error", errors)
				return nil, jwt.ErrFailedAuthentication
			}
			hash := user.Password
			if crypt.CheckPasswordHash(password, hash, 32) != true {
				c.Set("httpStatus", http.StatusUnauthorized)
				c.Set("error", errors)
				return nil, jwt.ErrFailedAuthentication
			}

			return &user, nil
		},
		LoginResponse: func(c *gin.Context, status int, msg string, time time.Time) {
			c.Redirect(http.StatusFound, "/toodles")
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
		// TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}
