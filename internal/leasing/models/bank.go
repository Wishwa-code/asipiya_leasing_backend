package models

import "gorm.io/gorm"

type Bank struct {
	gorm.Model
	Name      string `gorm:"column:name" json:"name"`
	Code      string `gorm:"column:code" json:"code"`
	ShortName string `gorm:"column:short_name" json:"short_name"`
	Status    string `gorm:"column:status;default:'Active'" json:"status"`
}

func (Bank) TableName() string {
	return "banks"
}
