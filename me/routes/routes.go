package routes

import (
	"github.com/Marcel-dev2009/cbt-backend/connections/handlers"
	"github.com/Marcel-dev2009/cbt-backend/me/middleware"
	"github.com/gin-gonic/gin"
)
func Setup (r *gin.Engine){
 api := r.Group("/api")
 {  // Makes all our grouped with /api
  auth := api.Group("/auth")
  {
  auth.POST("/register" , handlers.Register)
  auth.POST("/login" , handlers.Login)
  auth.POST("/logout" , handlers.Logout)
	
  }
  
  protected := api.Group("/")
  protected.Use(middleware.RequireAuth)
  {
   protected.GET("/me" , handlers.Me)
   protected.POST("/profile", handlers.SaveProfile)
   protected.POST("/save" , handlers.SaveResult);	          
   protected.GET("/result" , handlers.GetResults);
  }
 }	
}