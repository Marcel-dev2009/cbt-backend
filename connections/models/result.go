package models

import (
 "time"
)

type Result struct {
    ID      string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
    UserID  string    `gorm:"not null" json:"user_id"`
    Subject string    `gorm:"not null" json:"subject"`
    Score   int       `gorm:"not null" json:"score"`
    Total   int       `gorm:"not null" json:"total"`
    TakenAt time.Time `json:"taken_at"`

    User User `gorm:"foreignKey:UserID" json:"user"`
}	
