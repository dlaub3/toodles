package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func middlewareCSRF() gin.HandlerFunc {

	return func(c *gin.Context) {

		csrfToken, err := c.Request.Cookie("csrf")
		if err != nil {
			csrfToken, err = csrf(c)
		}
		c.Set("csrftoken", csrfToken.Value)

		if c.Request.Method != "GET" {
			validCsrf := isCSRFTokenValid(c)
			if !validCsrf {
				c.Set("httpStatus", http.StatusBadRequest)
				c.Set("error", "refresh")
				showErrorPage(c)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
