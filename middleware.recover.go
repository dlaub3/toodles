package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func middlewareRecover() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Set("genError", "😑 oh snap! Please try again.")
				c.Set("httpStatus", http.StatusInternalServerError)
				log.Printf("[Recovery] %s panic recovered.:\n%s\n", time.Now(), err)
				showErrorPage(c)
			}
		}()

		c.Next()
	}
}
