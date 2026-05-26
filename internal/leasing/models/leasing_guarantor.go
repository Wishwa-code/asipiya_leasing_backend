package models

import "gorm.io/gorm"

// LeasingGuarantor represents a guarantor for the leasing application
type LeasingGuarantor struct {
	gorm.Model
	LeasingApplicationID *uint  `gorm:"column:leasing_application_id" json:"leasing_application_id"`
	CustomerID           *uint  `gorm:"column:customer_id" json:"customer_id"`
	GuarantorIndex       int    `gorm:"column:guarantor_index" json:"guarantor_index"`
	Type                 string `gorm:"column:type;type:varchar(50)" json:"type"`

	// Relationships
	LeasingApplication *LeasingApplication `gorm:"foreignKey:LeasingApplicationID" json:"leasing_application,omitempty"`
	Customer           *Customer           `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (LeasingGuarantor) TableName() string {
	return "leasing_guarantors"
}
