package model

import "gorm.io/gorm"

type Comment struct{
	gorm.Model
	Uid uint `json:"uid" gorm:"not null"`
	Cid uint `json:"cid" gorm:"not null"`
	Content string `json:"content" gorm:"type:text; not null"`
}