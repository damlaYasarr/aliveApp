package models

import "gorm.io/gorm"

type Time struct {
	gorm.Model
	AIM_ID       int64    `json:"aim_id" gorm:"primaryKey"`
	COMPLETE_DAYS []string `json:"complete_days" gorm:"type:text[]"` // Dize dizisi olarak tanımlanmış alan
}