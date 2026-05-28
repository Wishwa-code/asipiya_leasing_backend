package models

import (
	"time"

	"gorm.io/gorm"
)

type LeasingVehicleAudit struct {
	gorm.Model
	LeasingVehicleID uint      `gorm:"column:leasing_vehicle_id" json:"leasing_vehicle_id"`
	FieldName        string    `gorm:"column:field_name;type:varchar(100)" json:"field_name"`
	OldValue         string    `gorm:"column:old_value;type:text" json:"old_value"`
	NewValue         string    `gorm:"column:new_value;type:text" json:"new_value"`
	ModifiedBy       string    `gorm:"column:modified_by;type:varchar(100)" json:"modified_by"`
	Timestamp        time.Time `gorm:"column:timestamp" json:"timestamp"`
}

func (LeasingVehicleAudit) TableName() string {
	return "leasing_vehicle_audits"
}
