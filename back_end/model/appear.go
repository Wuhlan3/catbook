package model

import "gorm.io/gorm"

type Appear struct{
	gorm.Model
	Cid uint `json:"cid" gorm:"not null"`
	LocationId uint `json:"location_id" gorm:"not null"`
	Specific string `json:"specific" gorm:"type:varchar(200)"`
}