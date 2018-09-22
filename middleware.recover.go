package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func middlewareRecover() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			// c.Writer.Status() == 500
			if err := recover(); err != nil {
				c.Keys["genError"] = "😑 oh snap! Please try again."
				log.Printf("[Recovery] %s panic recovered.:\n%s\n", time.Now(), err)
				showErrorPage(c)
			}
		}()

		c.Next()
	}
}
