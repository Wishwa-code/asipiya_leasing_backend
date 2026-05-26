package models

import (
	"gorm.io/gorm"
	"time"
)

// LeasingLoan represents the financial details of the lease
type LeasingLoan struct {
	gorm.Model
	LeasingVehicleID       *uint     `gorm:"column:leasing_vehicle_id" json:"leasing_vehicle_id"`
	LeasingApplicationID   *uint     `gorm:"column:leasing_application_id" json:"leasing_application_id"`
	LoanNo                 string    `gorm:"column:loan_no;type:varchar(100)" json:"loan_no"`
	LoanType               string    `gorm:"column:loan_type;type:varchar(50)" json:"loan_type"`
	CreatedDateTime        *time.Time `gorm:"column:created_date_time" json:"created_date_time"`
	ApprovalEndDateTime    *time.Time `gorm:"column:approval_end_date_time" json:"approval_end_date_time"`
	DisbursedDateTime      *time.Time `gorm:"column:disbursed_date_time" json:"disbursed_date_time"`
	CreatedBy              *uint     `gorm:"column:created_by" json:"created_by"`
	UpdatedBy              *uint     `gorm:"column:updated_by" json:"updated_by"`
	DisbursedUser          *uint     `gorm:"column:disbursed_user" json:"disbursed_user"`
	CustomerID             *uint     `gorm:"column:customer_id" json:"customer_id"`
	RouteID                *uint     `gorm:"column:route_id" json:"route_id"`
	CenterID               *uint     `gorm:"column:center_id" json:"center_id"`
	GroupsID               *uint     `gorm:"column:groups_id" json:"groups_id"`
	ProductID              *uint     `gorm:"column:product_id" json:"product_id"`
	ProductCode            string    `gorm:"column:product_code;type:varchar(50)" json:"product_code"`
	InterestMethod         string    `gorm:"column:interest_method;type:varchar(50)" json:"interest_method"`
	LoanPeriodType         string    `gorm:"column:loan_period_type;type:varchar(50)" json:"loan_period_type"`
	LoanPeriod             int       `gorm:"column:loan_period" json:"loan_period"`
	InterestPeriodType     string    `gorm:"column:interest_period_type;type:varchar(50)" json:"interest_period_type"`
	InterestPeriod         int       `gorm:"column:interest_period" json:"interest_period"`
	InterestRate           float64   `gorm:"column:interest_rate" json:"interest_rate"`
	CollectionPeriodType   string    `gorm:"column:collection_period_type;type:varchar(50)" json:"collection_period_type"`
	CollectionDuration     int       `gorm:"column:collection_duration" json:"collection_duration"`
	CollectionDateType     string    `gorm:"column:collection_date_type;type:varchar(50)" json:"collection_date_type"`
	FirstCollectionDate    string    `gorm:"column:first_collection_date;type:varchar(50)" json:"first_collection_date"`
	CollectionDay          int       `gorm:"column:collection_day" json:"collection_day"`
	ProductItemID          *uint     `gorm:"column:product_item_id" json:"product_item_id"`
	PenaltyActiveStatus    string    `gorm:"column:panelty_active_status;type:varchar(20)" json:"panelty_active_status"`
	PenaltyMethod          string    `gorm:"column:panelty_method;type:varchar(50)" json:"panelty_method"`
	PenaltyRate            float64   `gorm:"column:panelty_rate" json:"panelty_rate"`
	PenaltyApplyType       string    `gorm:"column:panelty_apply_type;type:varchar(50)" json:"panelty_apply_type"`
	PenaltyStartAfterDays  int       `gorm:"column:panelty_start_after_days" json:"panelty_start_after_days"`
	TotalAdditionalCharge  float64   `gorm:"column:total_additional_charge" json:"total_additional_charge"`
	RecoveryAccountStatus  string    `gorm:"column:recovery_account_status;type:varchar(20)" json:"recovery_account_status"`
	SavingAccountStatus    string    `gorm:"column:saving_acount_status;type:varchar(20)" json:"saving_acount_status"`
	SavingAmountType       string    `gorm:"column:saving_amount_type;type:varchar(50)" json:"saving_amount_type"`
	SavingCollectionType   string    `gorm:"column:saving_collection_type;type:varchar(50)" json:"saving_collection_type"`
	SavingAmount           float64   `gorm:"column:saving_amount" json:"saving_amount"`
	SavingAccountMonthlyInterest float64 `gorm:"column:saving_account_monthly_interest" json:"saving_account_monthly_interest"`
	SavingAmountOnEveryInstallments float64 `gorm:"column:saving_amount_on_every_installments" json:"saving_amount_on_every_installments"`
	LoanAmount             float64   `gorm:"column:loan_amount" json:"loan_amount"`
	InterestAmount         float64   `gorm:"column:interest_amount" json:"interest_amount"`
	TotalLoanAmount        float64   `gorm:"column:total_loan_amount" json:"total_loan_amount"`
	OtherChargesTotal      float64   `gorm:"column:other_charges_total" json:"other_charges_total"`
	OtherChargesOnDisburse float64   `gorm:"column:other_charges_on_disburse" json:"other_charges_on_disburse"`
	OtherChargesOnFirstInstallment float64 `gorm:"column:other_charges_on_first_installment" json:"other_charges_on_first_installment"`
	OtherChargesOnEveryInstallments float64 `gorm:"column:other_charges_on_every_installments" json:"other_charges_on_every_installments"`
	OtherChargesAdditionalAdded float64 `gorm:"column:other_charges_additional_added" json:"other_charges_additional_added"`
	DisburseAmount         float64   `gorm:"column:disburse_amount" json:"disburse_amount"`
	InstallmentAmount      float64   `gorm:"column:installment_amount" json:"installment_amount"`
	TotalPenaltyAmount     float64   `gorm:"column:total_panalty_amount" json:"total_panalty_amount"`
	TotalPaidAmount        float64   `gorm:"column:total_paid_amount" json:"total_paid_amount"`
	CapitalBalance         float64   `gorm:"column:capital_balance" json:"capital_balance"`
	InterestBalance        float64   `gorm:"column:interest_balance" json:"interest_balance"`
	TotalLoanBalance       float64   `gorm:"column:total_loan_balance" json:"total_loan_balance"`
	PenaltyBalance         float64   `gorm:"column:panalty_balance" json:"panalty_balance"`
	OtherChargesBalance    float64   `gorm:"column:other_charges_balance" json:"other_charges_balance"`
	AdditionalChargesBalance float64 `gorm:"column:additional_charges_balance" json:"additional_charges_balance"`
	CollectableSavingsBalance float64 `gorm:"column:collectable_savings_balance" json:"collectable_savings_balance"`
	TotalBalance           float64   `gorm:"column:total_balance" json:"total_balance"`
	LoanStatus             string    `gorm:"column:loan_status;type:varchar(50)" json:"loan_status"`
	LendingOfficerID       *uint     `gorm:"column:lending_officer_id" json:"lending_officer_id"`
	RecoveryOfficerID      *uint     `gorm:"column:recovery_officer_id" json:"recovery_officer_id"`
	BranchID               *uint     `gorm:"column:branch_id" json:"branch_id"`
	BankAccountID          *uint     `gorm:"column:bank_account_id" json:"bank_account_id"`
	DeletedBy              *uint     `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedReason          string    `gorm:"column:deleted_reson;type:text" json:"deleted_reson"`
	ArrearsAmount          float64   `gorm:"column:arrears_amount" json:"arrears_amount"`
	TodayDueAmount         float64   `gorm:"column:today_due_amount" json:"today_due_amount"`
	LastCalculatedAt       *time.Time `gorm:"column:last_calculated_at" json:"last_calculated_at"`
	MaturityDate           string    `gorm:"column:maturity_date;type:varchar(50)" json:"maturity_date"`
	InspectionDate         string    `gorm:"column:inspection_date;type:varchar(50)" json:"inspection_date"`

	// Relationships
	LeasingVehicle     *LeasingVehicle     `gorm:"foreignKey:LeasingVehicleID" json:"leasing_vehicle,omitempty"`
	LeasingApplication *LeasingApplication `gorm:"foreignKey:LeasingApplicationID" json:"leasing_application,omitempty"`
	Customer           *Customer           `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Product            *Product            `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (LeasingLoan) TableName() string {
	return "leasing_loans"
}
