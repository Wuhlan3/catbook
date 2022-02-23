package model

import "gorm.io/gorm"

type Location struct{
	gorm.Model
	Place string `json:"place" gorm:"type:varchar(50);not null;unique"`
}
