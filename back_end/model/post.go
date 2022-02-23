package model

import "gorm.io/gorm"

type Post struct{
	gorm.Model
	Uid uint `json:"uid" gorm:"not null"`
	Cid uint `json:"cid" gorm:"not null"`
	ImgUrl string `json:"img_url" gorm:"type:varchar(50);not null"`
}