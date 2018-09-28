package main

import (
	"github.com/gin-gonic/gin"
)

var r *gin.Engine
var config Config

func main() {
	initConfig()
	initDB()
	initRoutes()
}
