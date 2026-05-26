package models

import "gorm.io/gorm"

// LeasingVehicleDocumentImage maps to the former leasing_vehicle_has_document_images table.
type LeasingVehicleDocumentImage struct {
	gorm.Model
	LeasingVehicleID     *uint   `gorm:"column:leasing_vehicle_id" json:"leasing_vehicle_id"` // Nullable for early draft uploads
	LeasingApplicationID *uint   `gorm:"column:leasing_application_id" json:"leasing_application_id"` // Helps track early drafts
	ImageType            string  `gorm:"column:image_type;type:varchar(100)" json:"image_type"`
	ImageURL             string  `gorm:"column:image_url;type:text" json:"image_url"`

	// Relationships
	LeasingVehicle     *LeasingVehicle     `gorm:"foreignKey:LeasingVehicleID" json:"leasing_vehicle,omitempty"`
	LeasingApplication *LeasingApplication `gorm:"foreignKey:LeasingApplicationID" json:"leasing_application,omitempty"`
}

func (LeasingVehicleDocumentImage) TableName() string {
	return "leasing_vehicle_document_images"
}
