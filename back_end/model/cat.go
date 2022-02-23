package model

import (
	"gorm.io/gorm"
)

type Cat struct{
	gorm.Model
	CategoryId uint `json:"category_id" gorm:"not null"`
	Name string `json:"name" gorm:"type:varchar(20);not null;unique"`
	Color string `json:"color" gorm:"type:varchar(50);default:'未知'"`
}
