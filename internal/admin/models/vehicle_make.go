package models

import "gorm.io/gorm"

type VehicleMake struct {
	gorm.Model
	VehicleMake   string       `gorm:"column:vehicle_make;size:255;not null" json:"vehicle_make"`
	VehicleTypeID uint         `gorm:"column:vehicle_type_id;not null" json:"vehicle_type_id"`
	VehicleType   *VehicleType `gorm:"foreignKey:VehicleTypeID" json:"vehicle_type,omitempty"`
	Status        string       `gorm:"column:status;default:'Active'" json:"status"`
}

func (VehicleMake) TableName() string {
	return "vehicle_makes"
}
