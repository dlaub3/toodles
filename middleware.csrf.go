package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func middlewareCSRF() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()

		if c.Request.Method != "GET" {
			validRequest := isCSRFTokenValid(c)
			if !validRequest {
				c.AbortWithError(http.StatusInternalServerError, errorInternalError).SetType(gin.ErrorTypePublic)
			}
		}
	}
}
