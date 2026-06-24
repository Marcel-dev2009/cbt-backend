package middleware

import (
	"net/http"

	"github.com/Marcel-dev2009/cbt-backend/config"
	"github.com/Marcel-dev2009/cbt-backend/me/sessions"
	"github.com/gin-gonic/gin"
)
func RequireAuth(c *gin.Context){
  //Get session ID from cookie
  sessionID , err := c.Cookie("session_id")
  if err  != nil{
   c.JSON(http.StatusUnauthorized,gin.H{"error" : "not authenticated"})
   c.Abort()
   return	
  }	
  sessionData , err := sessions.GetSession(config.RedisClient , sessionID)
  if err != nil{
   c.JSON(http.StatusUnauthorized , gin.H{"error":"invalid or expired session"})
   c.Abort()
   return	
  }
   c.Set("userID" , sessionData.USERID)
   c.Set("email" , sessionData.Email)
   c.Next()
}