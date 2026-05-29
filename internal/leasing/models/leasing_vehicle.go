package models

import (
	"fmt"
	"time"

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
	RegistrationNo       string  `gorm:"column:registration_no;type:varchar(100)" json:"registration_no"`
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

	// New RMV / CR Master Data
	DateOfFirstRegistration string  `gorm:"column:date_of_first_registration;type:varchar(50)" json:"cr_first_reg_date"`
	AbsoluteOwner           string  `gorm:"column:absolute_owner;type:varchar(255)" json:"cr_absolute_owner"`
	RegisteredOwner         string  `gorm:"column:registered_owner;type:varchar(255)" json:"cr_registered_owner"`
	PreviousOwnersCount     int     `gorm:"column:previous_owners_count;default:0" json:"cr_previous_owners_count"`
	VariantCode             string  `gorm:"column:variant_code;type:varchar(100)" json:"cr_variant"`
	Transmission            string  `gorm:"column:transmission;type:varchar(50)" json:"cr_transmission"`
	SeatingCapacity         int     `gorm:"column:seating_capacity" json:"cr_seating_capacity"`
	UnladenWeight           float64 `gorm:"column:unladen_weight" json:"cr_unladen_weight"`
	MileageAtRegistration   float64 `gorm:"column:mileage_at_registration" json:"cr_mileage"`

	// CR Certificate Metadata
	CrSerialNo    string `gorm:"column:cr_serial_no;type:varchar(100)" json:"cr_serial_no"`
	CrIssueDate   string `gorm:"column:cr_issue_date;type:varchar(50)" json:"cr_issue_date"`
	CrIssueOffice string `gorm:"column:cr_issue_office;type:varchar(100)" json:"cr_issue_office"`
	CrType        string `gorm:"column:cr_type;type:varchar(100)" json:"cr_type"` // e.g. Original / Duplicate
	KeysReceived  int    `gorm:"column:keys_received;default:2" json:"cr_keys_received"`
	DuplicateKey  bool   `gorm:"column:duplicate_key;default:false" json:"duplicate_key"`

	// Verification Audit Fields
	RmvVerified       bool   `gorm:"column:rmv_verified;default:false" json:"rmv_verified"`
	RmvVerifiedAt     string `gorm:"column:rmv_verified_at;type:varchar(100)" json:"rmv_verified_at"`
	DataMatched       bool   `gorm:"column:data_matched;default:true" json:"data_matched"`
	DiscrepancyReason string `gorm:"column:discrepancy_reason;type:text" json:"discrepancy_reason"`

	// Relationships
	LeasingApplication *LeasingApplication           `gorm:"foreignKey:LeasingApplicationID" json:"leasing_application,omitempty"`
	VehicleType        *adminModels.VehicleType      `gorm:"foreignKey:VehicleTypeID" json:"vehicle_type,omitempty"`
	VehicleMake        *adminModels.VehicleMake      `gorm:"foreignKey:VehicleMakeID" json:"vehicle_make,omitempty"`
	VehicleModel       *adminModels.VehicleModel     `gorm:"foreignKey:VehicleModelID" json:"vehicle_model,omitempty"`
	Color              *adminModels.Color            `gorm:"foreignKey:ColorID" json:"color,omitempty"`
	ValuationCompany   *ValuationCompany             `gorm:"foreignKey:ValuationCompanyID" json:"valuation_company,omitempty"`
	InsuranceCompany   *InsuranceCompany             `gorm:"foreignKey:InsuranceCompanyID" json:"insurance_company,omitempty"`
	Supplier           *Supplier                     `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Images             []LeasingVehicleDocumentImage `gorm:"foreignKey:LeasingVehicleID" json:"images,omitempty"`
}

func (LeasingVehicle) TableName() string {
	return "leasing_vehicles"
}

// BeforeUpdate is a GORM hook called prior to updating a LeasingVehicle record
func (v *LeasingVehicle) BeforeUpdate(tx *gorm.DB) (err error) {
	// Retrieve the existing record from the database to compare values
	var old LeasingVehicle
	if err := tx.Unscoped().First(&old, v.ID).Error; err != nil {
		return nil // skip if record does not exist
	}

	// We only track modifications post-RMV check (i.e. if RMV was already verified in the database)
	if !old.RmvVerified {
		return nil
	}

	// Helper function to log field changes
	logChange := func(fieldName, oldValue, newValue string) {
		if oldValue == newValue {
			return
		}

		modifiedBy := "system"
		if tx.Statement.Context != nil {
			if uVal := tx.Statement.Context.Value("username"); uVal != nil {
				if uStr, ok := uVal.(string); ok {
					modifiedBy = uStr
				}
			}
		}

		audit := LeasingVehicleAudit{
			LeasingVehicleID: v.ID,
			FieldName:        fieldName,
			OldValue:         oldValue,
			NewValue:         newValue,
			ModifiedBy:       modifiedBy,
			Timestamp:        time.Now(),
		}
		tx.Create(&audit)
	}

	// Compare CR Master Data fields
	logChange("DateOfFirstRegistration", old.DateOfFirstRegistration, v.DateOfFirstRegistration)
	logChange("AbsoluteOwner", old.AbsoluteOwner, v.AbsoluteOwner)
	logChange("RegisteredOwner", old.RegisteredOwner, v.RegisteredOwner)
	logChange("PreviousOwnersCount", fmt.Sprintf("%d", old.PreviousOwnersCount), fmt.Sprintf("%d", v.PreviousOwnersCount))
	logChange("VariantCode", old.VariantCode, v.VariantCode)
	logChange("Transmission", old.Transmission, v.Transmission)
	logChange("SeatingCapacity", fmt.Sprintf("%d", old.SeatingCapacity), fmt.Sprintf("%d", v.SeatingCapacity))
	logChange("UnladenWeight", fmt.Sprintf("%.2f", old.UnladenWeight), fmt.Sprintf("%.2f", v.UnladenWeight))
	logChange("MileageAtRegistration", fmt.Sprintf("%.2f", old.MileageAtRegistration), fmt.Sprintf("%.2f", v.MileageAtRegistration))
	logChange("CrSerialNo", old.CrSerialNo, v.CrSerialNo)
	logChange("CrIssueDate", old.CrIssueDate, v.CrIssueDate)
	logChange("CrIssueOffice", old.CrIssueOffice, v.CrIssueOffice)
	logChange("CrType", old.CrType, v.CrType)

	return nil
}
