package main

import (
	"log"
	"os"

	"github.com/Marcel-dev2009/cbt-backend/config"
	"github.com/Marcel-dev2009/cbt-backend/me/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)
func main(){
  config.Load()
  r := gin.Default()
  r.Use(cors.New(config.CORSConfig()));
  r.GET("/" , func(c *gin.Context) {
    c.JSON(200 , gin.H{"status":"backend server is running smoothly"})
  })
  routes.Setup(r)
  port := os.Getenv("PORT")
  if port == ""{
  port = "8080"	
  }
  log.Println("Server is running on port" , port)
  r.Run(":" + port)
}