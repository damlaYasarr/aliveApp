package models 
import "gorm.io/gorm"

type User struct{
	gorm.Model:
	ID  int64 `json:"id" gorm:"primarykey"`
	Email string `json:"email"` `gorm:"unique;not null;type:varchar(100);default:null"`
	AIM_ID int64 `json:"AIM_ID"` 
}

//payment detail will be added