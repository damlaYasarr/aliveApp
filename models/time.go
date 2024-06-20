package models

import ("gorm.io/gorm"
  "github.com/lib/pq")

type Time struct {
	gorm.Model
	AIM_ID       int64    `json:"aim_id" gorm:"primaryKey"`
	COMPLETE_DAYS pq.StringArray `gorm:"type:text[]"`
}