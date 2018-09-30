package main

import (
	"github.com/gin-gonic/gin"
)

func initRoutes() {

	// gin.SetMode(gin.ReleaseMode)
	r = gin.New()
	// r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middlewareCSRF())
	r.Use(middlewareErrors())
	r.Use(middlewareRecover())

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	authMiddleware := jwtMiddleware()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.NoRoute(func(c *gin.Context) {
		render(c, gin.H{"error": "Not Found"}, "404.html")
	})

	r.GET("/", showHomePage)
	r.GET("/signup", showSignupPage)
	r.POST("/signup", registerNewUser)
	r.GET("/login", showLoginPage)
	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/logout", logout)

	auth := r.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	auth.GET("refresh_token", authMiddleware.RefreshHandler)

	auth.GET("/toodles", getAllToodles)
	auth.POST("/toodles", createAToodle)
	auth.GET("/toodles/:toodle_id", getAToodle)
	auth.PUT("/toodles/:toodle_id", updateAToodle)
	auth.DELETE("/toodles/:toodle_id", deleteAToodle)
	auth.PUT("/toodles/:toodle_id/complete", completeAToodle)

	//Routes specifically for form submitalls and not AJAX
	auth.POST("/toodles/:toodle_id", updateOrDeleteAToodle)
	auth.POST("/toodles/:toodle_id/complete", completeAToodle)

	r.Run()
}
