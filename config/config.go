package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Marcel-dev2009/cbt-backend/connections/models"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
  "github.com/cloudinary/cloudinary-go/v2"
	"gorm.io/gorm"
)
var DB *gorm.DB
var RedisClient *redis.Client // Defining type of our newly created variables without initializing them with a value
var CloudinaryClient *cloudinary.Cloudinary;
func Load (){
 _ = godotenv.Load()
 connectDB()
 connectRedis()
 connectCloudinary()	
}

func connectDB(){
 dsn := os.Getenv("DATABASE_URL");
 db , err := gorm.Open(postgres.Open(dsn) , &gorm.Config{})
 if err != nil{
  log.Fatal("Database Connection Failed" , err)	
 }
 DB = db
 err = db.AutoMigrate(&models.User{} , &models.Result{} , &models.SaveResult{});
 if err != nil {
  log.Fatal("Migration Failed" , err)
 }
 fmt.Println("✅ Models migrated")
 fmt.Println("✅ Database connected")	
}

func connectRedis (){
 opt , err := redis.ParseURL(os.Getenv("REDIS_URL")) // Just parses your url
 if err != nil {
  log.Fatal("Failed to parse Redis URL" , err)	
 }
  client := redis.NewClient(opt)
  // Test the connection
  ctx := context.Background()
  _,err = client.Ping(ctx).Result()
  if err != nil {
    log.Fatal("Failed to Connect to Redis" , err)	
  }
  RedisClient = client
  fmt.Println("✅ Redis connected")
}

func CORSConfig() cors.Config{
  config := cors.DefaultConfig()
  config.AllowOrigins = []string{"https://prep-mate-xyz.vercel.app/","http://localhost:3000"}
  config.AllowMethods = []string{"GET" , "POST" ,"PUT" , "PATCH" , "DELETE" , "OPTIONS"}
  config.AllowHeaders = []string{"Origin","Content-Type","Authorization"}
  config.AllowCredentials = true
  return  config
}
func connectCloudinary (){
  cld,err := cloudinary.NewFromParams(
    os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
  )
  if err != nil{
    log.Fatal("failed to connect to cloudinary", err)
  }
  CloudinaryClient = cld
  fmt.Println("Cloudinary connected")
}