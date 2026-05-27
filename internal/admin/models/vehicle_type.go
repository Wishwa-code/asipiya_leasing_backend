package models

import "gorm.io/gorm"

type VehicleType struct {
	gorm.Model
	VehicleTypeName string `gorm:"column:vehicle_type_name;size:255;not null" json:"vehicle_type_name"`
	Description     string `gorm:"column:description;type:text" json:"description"`
	Status          string `gorm:"column:status;default:'Active'" json:"status"`
}

func (VehicleType) TableName() string {
	return "vehicle_types"
}
