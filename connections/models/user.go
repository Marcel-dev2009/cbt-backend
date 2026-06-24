package models

import (
	"time"
)
type User struct {
  ID string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
  Username string `gorm:"uniqueIndex;not null" json:"username"`
  Email string  `gorm:"uniqueIndex;not null" json:"email"` 
  Grade string `json:"grade"`
  School string `json:"school"`
  Profile string `json:"profile"`
  Password string `gorm:"not null" json:"-"`
  // Result int `json:"result"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

/* 
You need a separate result in the db 
*/
