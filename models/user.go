package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID    int64  `json:"id" gorm:"primaryKey"`
	Email string `json:"email" gorm:"unique;not null;type:varchar(100);default:null"`
	AIMID int64  `json:"AIM_ID"`
}

// Payment detail will be added