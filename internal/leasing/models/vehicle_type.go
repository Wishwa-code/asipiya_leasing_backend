package models

import "gorm.io/gorm"

type VehicleType struct {
	gorm.Model
	VehicleTypeName string `gorm:"column:vehicle_type_name" json:"vehicle_type_name"`
	Description     string `gorm:"column:description" json:"description"`
	Status          string `gorm:"column:status;default:'Active'" json:"status"`
}

func (VehicleType) TableName() string {
	return "vehicle_types"
}
