package models

import (
	adminModels "garment-management-backend/internal/admin/models"
	"gorm.io/gorm"
)

// LeasingVehicle represents the vehicle linked to a Leasing Application
type LeasingVehicle struct {
	gorm.Model
	LeasingApplicationID *uint   `gorm:"column:leasing_application_id" json:"leasing_application_id"`
	VehicleTypeID        *uint   `gorm:"column:vehicle_type_id" json:"vehicle_type_id"`
	VehicleMakeID        *uint   `gorm:"column:vehicle_make_id" json:"vehicle_make_id"`
	VehicleModelID       *uint   `gorm:"column:vehicle_model_id" json:"vehicle_model_id"`
	VehicleStatus        string  `gorm:"column:vehicle_status;type:varchar(50)" json:"vehicle_status"`
	EngineCc             string  `gorm:"column:engine_cc;type:varchar(50)" json:"engine_cc"`
	ChasisNo             string  `gorm:"column:chasis_no;type:varchar(100)" json:"chasis_no"`
	ManufacturingYear    string  `gorm:"column:manufacturing_year;type:varchar(10)" json:"manufacturing_year"`
	ColorID              *uint   `gorm:"column:color_id" json:"color_id"`
	Usage                string  `gorm:"column:usage;type:varchar(255)" json:"usage"`
	CountryOfOrigin      string  `gorm:"column:country_of_origin;type:varchar(100)" json:"country_of_origin"`
	TypeOfBody           string  `gorm:"column:type_of_body;type:varchar(100)" json:"type_of_body"`
	Equipment            string  `gorm:"column:equipment;type:text" json:"equipment"`
	RegisteredYear       string  `gorm:"column:registered_year;type:varchar(10)" json:"registered_year"`
	RegisteredNo         string  `gorm:"column:registered_no;type:varchar(100)" json:"registered_no"`
	ValuationCompanyID   *uint   `gorm:"column:valuation_company_id" json:"valuation_company_id"`
	InsuranceCompanyID   *uint   `gorm:"column:insurance_company_id" json:"insurance_company_id"`
	InsuranceAmount      float64 `gorm:"column:insurance_amount" json:"insurance_amount"`
	InsurancePremium     float64 `gorm:"column:insurance_premium" json:"insurance_premium"`
	InsuranceStartDate   string  `gorm:"column:insurance_start_date;type:varchar(50)" json:"insurance_start_date"`
	InsuranceExpiryDate  string  `gorm:"column:insurance_expiry_date;type:varchar(50)" json:"insurance_expiry_date"`
	SupplierID           *uint   `gorm:"column:supplier_id" json:"supplier_id"`
	SupplierRno          string  `gorm:"column:supplier_rno;type:varchar(100)" json:"supplier_rno"`
	MarketValue          float64 `gorm:"column:market_value" json:"market_value"`
	ForcedSaleValue      float64 `gorm:"column:forced_sale_value" json:"forced_sale_value"`
	InvoiceValue         float64 `gorm:"column:invoice_value" json:"invoice_value"`

	// Relationships
	LeasingApplication *LeasingApplication           `gorm:"foreignKey:LeasingApplicationID" json:"leasing_application,omitempty"`
	VehicleType        *adminModels.VehicleType      `gorm:"foreignKey:VehicleTypeID" json:"vehicle_type,omitempty"`
	VehicleMake        *adminModels.VehicleMake      `gorm:"foreignKey:VehicleMakeID" json:"vehicle_make,omitempty"`
	VehicleModel       *adminModels.VehicleModel     `gorm:"foreignKey:VehicleModelID" json:"vehicle_model,omitempty"`
	ValuationCompany   *ValuationCompany             `gorm:"foreignKey:ValuationCompanyID" json:"valuation_company,omitempty"`
	InsuranceCompany   *InsuranceCompany             `gorm:"foreignKey:InsuranceCompanyID" json:"insurance_company,omitempty"`
	Supplier           *Supplier                     `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Images             []LeasingVehicleDocumentImage `gorm:"foreignKey:LeasingVehicleID" json:"images,omitempty"`
}

func (LeasingVehicle) TableName() string {
	return "leasing_vehicles"
}
