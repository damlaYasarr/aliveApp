package models 
import "gorm.io/gorm"

type Aim struct{
	gorm.Model:
	ID  int64 `json:"id" gorm:"primarykey"`
	Name string `json:"name"` `gorm:"type:varchar(50);default:null"`
	Startday string `json:"startday"` `gorm:"type:varchar(50);default:null"`
	Endday string `json:"endday"` `gorm:"type:varchar(100);default:null"`
	Notification_hour string `json:"notification_hour"` `gorm:"type:varchar(10);default:null"`

}