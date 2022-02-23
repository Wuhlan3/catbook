package model

import "gorm.io/gorm"

type Category struct{
	gorm.Model
	Breed string `json:"breed" gorm:"type:varchar(50)"`
	Description string `json:"content" gorm:"type:text; not null"`
}