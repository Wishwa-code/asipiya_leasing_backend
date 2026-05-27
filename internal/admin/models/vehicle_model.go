package models

import "gorm.io/gorm"

type VehicleModel struct {
	gorm.Model
	VehicleModelName string       `gorm:"column:vehicle_model_name;size:255;not null" json:"vehicle_model_name"`
	VehicleTypeID    uint         `gorm:"column:vehicle_type_id;not null" json:"vehicle_type_id"`
	VehicleType      *VehicleType `gorm:"foreignKey:VehicleTypeID" json:"vehicle_type,omitempty"`
	Status           string       `gorm:"column:status;default:'Active'" json:"status"`
}

func (VehicleModel) TableName() string {
	return "vehicle_models"
}
