package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	// "time"

	"github.com/Marcel-dev2009/cbt-backend/config"
	"github.com/Marcel-dev2009/cbt-backend/connections/models"
	"github.com/Marcel-dev2009/cbt-backend/me/sessions"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)
type RegisterInput struct {
 Username string `json:"username" binding:"required,min=3"`
 Email string `json:"email" binding:"required,email"`
 Password string `json:"password" binding:"required,min=6,max=15"`	
}
type Logininput struct {
 Email string  `json:"email" binding:"required,email"`
 Password string `json:"password" binding:"required"`
}
type Profileinput struct {
 Grade string `json:"grade" binding:"required"`
 School string `json:"school" binding:"required"`
 Profile string `json:"profile" binding:"required"`
}
type ResultInput struct{
  Subject string `json:"subject" binding:"required"`
  Score  int  `json:"score" binding:"required"`
  Total int `json:"total" binding:"required"`
}
func SaveProfile(c *gin.Context) {
 var input Profileinput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get userID from session set by middleware
	userID := c.MustGet("userID").(string)

	// Upload profile photo to Cloudinary
	photoURL := ""
	if input.Profile != "" {
		// Strip the base64 prefix e.g "data:image/png;base64,..."
		base64Data := input.Profile
		if idx := strings.Index(base64Data, ","); idx != -1 {
			base64Data = base64Data[idx+1:]
		}

		// Decode base64 string to bytes
		imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image data"})
			return
		}

		// Upload to Cloudinary
		ctx := c.Request.Context()
	uploadResult, err := config.CloudinaryClient.Upload.Upload(ctx, imageBytes, uploader.UploadParams{
    Folder:   "cbt-app/avatars",
    PublicID: fmt.Sprintf("user-%s", userID),
    	
});
if err != nil {
    fmt.Println("Cloudinary upload error:", err) // add this line
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload photo"})
    return
}
  fmt.Println("Upload result:", uploadResult)
fmt.Println("Secure URL:", uploadResult.SecureURL)
fmt.Println("Error:", uploadResult.Error)
		photoURL = uploadResult.SecureURL
	}

	// Update user record in DB
	updates := map[string]interface{}{
		"grade":  input.Grade,
		"school": input.School,
	}
	if photoURL != "" {
		updates["profile"] = photoURL
	}

	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "profile saved successfully",
		"photo_url": photoURL,
	})
};
func SaveResult(c *gin.Context){
 var input ResultInput
 if err := c.ShouldBindJSON(&input); err != nil{
  c.JSON(http.StatusBadRequest , gin.H{"error":err.Error()})
  return;
 } 
 userID := c.MustGet("userID").(string)
 result := models.Result{
  UserID: userID,
  Subject: input.Subject,
  Score: input.Score,
  Total: input.Total,
  TakenAt: time.Now(),
 }
 if err := config.DB.Create(&result).Error; err != nil{
  fmt.Println("DB error", err)
  c.JSON(http.StatusUnauthorized , gin.H{"error":"failed to save user"})
  return
 }
 c.JSON(http.StatusOK , gin.H{
  "message" : "Your attempt has been submitted , thank you!",
  "result" : gin.H{
   "subject" : result.Subject,
   "score" : result.Score,
   "total" : result.Total,
  },
 })
}

func Register (c *gin.Context){
 var input RegisterInput
 if err := c.ShouldBindJSON(&input); err != nil {
	c.JSON(http.StatusBadRequest,gin.H{"error" : err.Error()})
	 return 
 }	
var existing models.User
if err := config.DB.Where("email = ? " , input.Email).First(&existing).Error; err == nil{
 c.JSON(http.StatusConflict, gin.H{"error":"email already exists"})
 return	
}
hashedPassword , err := bcrypt.GenerateFromPassword([]byte(input.Password) ,bcrypt.DefaultCost )
if err != nil {
 c.JSON(http.StatusInternalServerError , gin.H{"error" : "failed to hash password"})
 return	
}
 user := models.User{
  Username: input.Username,
  Email: input.Email,
  Password: string(hashedPassword),	
 }
 if err := config.DB.Create(&user).Error; err != nil{
  c.JSON(http.StatusInternalServerError , gin.H{"error" : "error creating user"})	
  return
 }
 // create session
 sessionData := sessions.SessionData{
  USERID: user.ID,	
  Email: user.Email,
 }
 sessionID , err := sessions.CreateSession(config.RedisClient , sessionData)
 if err != nil {
  c.JSON(http.StatusInternalServerError , gin.H{"error" : "failed to create session"})
  return
 }
 //setting my cookie
 c.SetCookie("session_id" , sessionID , 172800, "/" , "", false,true )
 c.JSON(http.StatusCreated ,gin.H{
   "message" : "registered sucessfully",
   "user" : gin.H{
    "id" : user.ID,
    "email" : user.Email,	
   },	
 })
}
func Login(c *gin.Context){
  var input Logininput
  if err := c.ShouldBindJSON(&input); err != nil {
   c.JSON(http.StatusBadRequest , gin.H{"error":err.Error()})
   return	
  }
  // find user by email
  var user  models.User;
  if err := config.DB.Where("email = ?" , input.Email).First(&user).Error; err != nil{
   c.JSON(http.StatusUnauthorized , gin.H{"error" :"No account found for this user"});
   return	
  }
  // compare password
  if err := bcrypt.CompareHashAndPassword([]byte(user.Password) , []byte(input.Password)); err != nil{
    c.JSON(http.StatusUnauthorized , gin.H{"error" : "Invalid credentials"})
    return	
  }	
  // create session
  sessionData := sessions.SessionData{
  USERID : user.ID,
  Email : user.Email,
  }
  sessionID , err := sessions.CreateSession(config.RedisClient , sessionData)
  if err != nil{
  c.JSON(http.StatusInternalServerError , gin.H{"error" : "failed to create session"})
  return
  }
  //set cookie
  c.SetCookie("session_id" , sessionID , 172800, "/" , "" , false ,true)
  c.JSON(http.StatusOK , gin.H{
   "message" : "logged in successfully",
   "user" : gin.H{
   "id" : user.ID,
   "email" : user.Email,	
   },	
  })
}

func Logout (c *gin.Context){
  sessionID, err := c.Cookie("session_id")
  if err != nil{
   c.JSON(http.StatusUnauthorized , gin.H{"error" : "no session found"})
   return	
  }	
  sessions.DeleteSession(config.RedisClient, sessionID)
c.SetCookie("session_id", "" ,-1 ,"/" , "" , false ,true)
c.JSON(http.StatusOK , gin.H{"message":"logged out successfully"})
}
func Me(c *gin.Context){
  sessionID , err := c.Cookie("session_id")
  if err != nil{
	c.JSON(http.StatusUnauthorized , gin.H{"error" : "not authenticated"})
	return
  }	
  sessionData, err := sessions.GetSession(config.RedisClient, sessionID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
		return
	}
  c.JSON(http.StatusOK ,gin.H{
   "user" : gin.H{
    "id" : sessionData.USERID,
    "email": sessionData.Email,	
   },	
  })
}
func GetResults(c *gin.Context){
  UserID := c.MustGet("userID").(string)
  var result []models.Result
  if err := config.DB.Preload("User").Where("user_id = ?" , UserID).Order("taken_at desc").Find(&result).Error; err != nil{
    c.JSON(http.StatusInternalServerError , gin.H{"error":"failed to get results"})
    return
  }  
  var response []gin.H
  for _,r := range result{
    response = append(response, gin.H{
     "username" : r.User.Username,
     "email" : r.User.Email,
      "subject" : r.Subject,
      "score" :r.Score,
      "total" : r.Total,
      "taken_at": r.TakenAt,
    })
  }
  c.JSON(http.StatusOK , gin.H{"result":response});
}