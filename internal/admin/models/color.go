package models

import "gorm.io/gorm"

type Color struct {
	gorm.Model
	ColorName string `gorm:"column:color_name;size:100;not null;unique" json:"color_name"`
	Status    string `gorm:"column:status;default:'Active'" json:"status"`
}

func (Color) TableName() string {
	return "colors"
}
