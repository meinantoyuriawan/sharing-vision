package models

import "time"

type Posts struct {
	Id           int       `gorm:"primaryKey" json:"id"`
	Title        string    `gorm:"varchar(200)" json:"title"`
	Content      string    `gorm:"varchar(300)" json:"content"`
	Category     string    `gorm:"varchar(100)" json:"category"`
	Created_date time.Time `gorm:"autoCreateTime:false"`
	Updated_date time.Time `gorm:"autoCreateTime:false"`
	Status       string    `gorm:"varchar(100)" json:"status"`
}
