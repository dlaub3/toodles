package main

import (
	"github.com/gin-gonic/gin"
)

func initializeRoutes() {

	authMiddleware := jwtMiddleware()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/", showHomePage)
	r.GET("/signup", showSignupPage)
	r.POST("/signup", registerNewUser)
	r.GET("/login", showLoginPage)
	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/logout", logout)

	auth := r.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/toodles", getAllToodles)
		auth.POST("/toodles", createAToodle)
		auth.GET("/toodles/:toodle_id", getAToodle)
		auth.PUT("/toodles/:toodle_id", updateAToodle)
		auth.DELETE("/toodles/:toodle_id", deleteAToodle)
		auth.PUT("/toodles/:toodle_id/complete", completeAToodle)

		//Routes specifically for form submitalls and not AJAX
		auth.POST("/toodles/:toodle_id", updateOrDeleteAToodle)
		auth.POST("/toodles/:toodle_id/complete", completeAToodle)

		auth.GET("refresh_token", authMiddleware.RefreshHandler)
	}

}
