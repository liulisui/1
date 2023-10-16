package main

import (
	"CMS/config/database"
	"CMS/config/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	r := gin.Default()
	router.Init(r)

	err := r.Run()
	if err != nil {
		log.Fatal("Server start failed: ", err)
	}
}
