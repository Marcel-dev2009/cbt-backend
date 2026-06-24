package sessions
import (
  "context"
  "encoding/json"
  "fmt"
  "time"

 "github.com/google/uuid"
 "github.com/redis/go-redis/v9"
)
type SessionData struct{
  USERID string `json:"user_id"`
  Email string `json:"email"`	
}
const sessionDuration = 48 * time.Hour
func CreateSession(client *redis.Client , data SessionData)(string , error) {
 sessionID := uuid.New().String()
 json, err := json.Marshal(data)
 if err != nil {
  return "", fmt.Errorf("failed to marshal session : %w" , err)	
 }
 ctx := context.Background()
 err = client.Set(ctx , sessionID , json , sessionDuration).Err()
 if err != nil {
  return "", fmt.Errorf("failed to store session %w",err)	
 }
 return sessionID , nil
}
func GetSession(client *redis.Client , sessionID string) (*SessionData , error) {
  ctx := context.Background()
  val , err := client.Get(ctx , sessionID).Result()
  if err == redis.Nil{
   return nil, fmt.Errorf("session not found")	
  }
  if err != nil {
   return nil, fmt.Errorf("failed to get session : %w" , err)
  } 	
  var data SessionData
  err = json.Unmarshal([]byte(val),&data)
  if err != nil {
  return nil, fmt.Errorf("failed to unmarshal session: %w" , err);	
  }
  return  &data , nil
}
func DeleteSession(client *redis.Client , sessionID string) error{
  ctx := context.Background()
  err := client.Del(ctx,sessionID).Err()
  if err != nil {
   return fmt.Errorf("failed to delete session : %w" , err)	
  }	
  return nil
}