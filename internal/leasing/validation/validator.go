package validation

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type StepValidator interface {
	Validate(db *gorm.DB, data map[string]interface{}) map[string]string
}

// Helper to safely parse any float/int/string payload value into a uint
func parseUintVal(val interface{}) (uint, bool) {
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case float64:
		return uint(v), true
	case int:
		return uint(v), true
	case string:
		if v == "" {
			return 0, false
		}
		var u uint
		if _, err := fmt.Sscanf(v, "%d", &u); err == nil {
			return u, true
		}
	}
	return 0, false
}

// -----------------------------------------------------------------------------
// Step 1: Customer Validator
// -----------------------------------------------------------------------------
type CustomerStepValidator struct{}

func (v *CustomerStepValidator) Validate(db *gorm.DB, data map[string]interface{}) map[string]string {
	errors := make(map[string]string)
	cidVal, ok := data["customer_id"]
	if !ok || cidVal == nil || fmt.Sprintf("%v", cidVal) == "" {
		errors["customer_id"] = "Customer is required"
	} else if cid, okParsed := parseUintVal(cidVal); okParsed {
		var count int64
		db.Table("customers").Where("id = ? AND deleted_at IS NULL", cid).Count(&count)
		if count == 0 {
			errors["customer_id"] = "Selected customer does not exist"
		}
	} else {
		errors["customer_id"] = "Invalid customer ID format"
	}
	return errors
}

// -----------------------------------------------------------------------------
// Step 2: Introducers Validator
// -----------------------------------------------------------------------------
type IntroducerStepValidator struct{}

func (v *IntroducerStepValidator) Validate(db *gorm.DB, data map[string]interface{}) map[string]string {
	errors := make(map[string]string)
	introsVal, ok := data["introducers"]
	if ok && introsVal != nil {
		if intros, ok := introsVal.([]interface{}); ok {
			for idx, item := range intros {
				if introMap, ok := item.(map[string]interface{}); ok {
					introIDVal, present := introMap["introducer_id"]
					if !present || introIDVal == nil || fmt.Sprintf("%v", introIDVal) == "" {
						errors[fmt.Sprintf("introducers.%d.introducer_id", idx)] = "Introducer is required"
					} else if introID, okParsed := parseUintVal(introIDVal); okParsed {
						var count int64
						db.Table("introducers").Where("id = ? AND deleted_at IS NULL", introID).Count(&count)
						if count == 0 {
							errors[fmt.Sprintf("introducers.%d.introducer_id", idx)] = "Selected introducer does not exist"
						}
					} else {
						errors[fmt.Sprintf("introducers.%d.introducer_id", idx)] = "Invalid introducer ID format"
					}
				}
			}
		}
	}
	return errors
}

// -----------------------------------------------------------------------------
// Step 3: Vehicle Validator
// -----------------------------------------------------------------------------
type VehicleStepValidator struct{}

func (v *VehicleStepValidator) Validate(db *gorm.DB, data map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	// Helper to get string val or fallback
	getStringVal := func(keys ...string) (string, string) {
		for _, k := range keys {
			if val, ok := data[k]; ok && val != nil {
				return fmt.Sprintf("%v", val), k
			}
		}
		return "", keys[0]
	}

	// 1. Asset Type (vehicle_type_id)
	typeVal, typeKey := getStringVal("vehicle_type_id")
	if typeVal == "" {
		errors[typeKey] = "Asset type is required"
	} else if tid, ok := parseUintVal(typeVal); ok {
		var count int64
		db.Table("vehicle_types").Where("id = ? AND status = 'Active'", tid).Count(&count)
		if count == 0 {
			errors[typeKey] = "Selected asset type does not exist or is inactive"
		}
	} else {
		errors[typeKey] = "Invalid asset type format"
	}

	// 2. Make (vehicle_make_id)
	makeVal, makeKey := getStringVal("vehicle_make_id")
	if makeVal == "" {
		errors[makeKey] = "Vehicle make is required"
	} else if mid, ok := parseUintVal(makeVal); ok {
		var count int64
		db.Table("vehicle_makes").Where("id = ? AND status = 'Active'", mid).Count(&count)
		if count == 0 {
			errors[makeKey] = "Selected vehicle make does not exist or is inactive"
		}
	} else {
		errors[makeKey] = "Invalid make format"
	}

	// 3. Model (vehicle_model_id)
	modelVal, modelKey := getStringVal("vehicle_model_id")
	if modelVal == "" {
		errors[modelKey] = "Vehicle model is required"
	} else if mid, ok := parseUintVal(modelVal); ok {
		var count int64
		db.Table("vehicle_models").Where("id = ? AND status = 'Active'", mid).Count(&count)
		if count == 0 {
			errors[modelKey] = "Selected vehicle model does not exist or is inactive"
		}
	} else {
		errors[modelKey] = "Invalid model format"
	}

	// 4. Asset Status (vehicle_status)
	statusVal, statusKey := getStringVal("vehicle_status")
	if statusVal == "" {
		errors[statusKey] = "Asset status is required"
	}

	// 5. Engine CC (engine_cc)
	engineVal, engineKey := getStringVal("engine_cc")
	if engineVal == "" {
		errors[engineKey] = "Engine CC is required"
	}

	// 6. Chassis No (chassis_no / chasis_no)
	chassisVal, chassisKey := getStringVal("chassis_no", "chasis_no")
	if chassisVal == "" {
		errors[chassisKey] = "Chassis number is required"
	}

	// 7. Manufacturing Year (manu_year / manufacturing_year)
	manuVal, manuKey := getStringVal("manu_year", "manufacturing_year")
	if manuVal == "" {
		errors[manuKey] = "Manufacturing year is required"
	}

	// 8. Color (color_id)
	colorVal, colorKey := getStringVal("color_id")
	if colorVal == "" {
		errors[colorKey] = "Color is required"
	} else if cid, ok := parseUintVal(colorVal); ok {
		var count int64
		db.Table("colors").Where("id = ? AND status = 'Active'", cid).Count(&count)
		if count == 0 {
			errors[colorKey] = "Selected color does not exist or is inactive"
		}
	} else {
		errors[colorKey] = "Invalid color format"
	}

	// 9. Usage (usage_type / usage)
	usageVal, usageKey := getStringVal("usage_type", "usage")
	if usageVal == "" {
		errors[usageKey] = "Usage is required"
	}

	// 10. Country of Origin (manu_country / country_of_origin)
	countryVal, countryKey := getStringVal("manu_country", "country_of_origin")
	if countryVal == "" {
		errors[countryKey] = "Country of origin is required"
	}

	// 11. Reg Year (reg_year / registered_year)
	regYearVal, regYearKey := getStringVal("reg_year", "registered_year")
	if regYearVal == "" {
		errors[regYearKey] = "Registered year is required"
	}

	// 12. Reg No (reg_no / registered_no / registration_no)
	regNoVal, regNoKey := getStringVal("reg_no", "registered_no", "registration_no")
	if regNoVal == "" {
		errors[regNoKey] = "Registered number is required"
	}

	// 13. Supplier (supplier_id)
	supplierVal, supplierKey := getStringVal("supplier_id")
	if supplierVal == "" {
		errors[supplierKey] = "Supplier is required"
	} else if sid, ok := parseUintVal(supplierVal); ok {
		var count int64
		db.Table("suppliers").Where("id = ? AND deleted_at IS NULL", sid).Count(&count)
		if count == 0 {
			errors[supplierKey] = "Selected supplier does not exist"
		}
	} else {
		errors[supplierKey] = "Invalid supplier format"
	}

	validateNumeric := func(fieldName string, label string) {
		val, present := data[fieldName]
		if !present || val == nil || fmt.Sprintf("%v", val) == "" {
			errors[fieldName] = label + " is required"
			return
		}
		var num float64
		if _, err := fmt.Sscanf(fmt.Sprintf("%v", val), "%f", &num); err != nil {
			errors[fieldName] = label + " must be a valid number"
		} else if num < 0 {
			errors[fieldName] = label + " must be greater than or equal to 0"
		}
	}

	validateNumeric("market_value", "Market value")
	validateNumeric("forced_value", "Forced sale value")
	validateNumeric("invoice_value", "Invoice value")

	return errors
}

// -----------------------------------------------------------------------------
// Step 4: Insurance Validator
// -----------------------------------------------------------------------------
type InsuranceStepValidator struct{}

func (v *InsuranceStepValidator) Validate(db *gorm.DB, data map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	company, ok := data["insurance_company"]
	if !ok || company == nil || fmt.Sprintf("%v", company) == "" {
		errors["insurance_company"] = "Insurance company is required"
	} else {
		var count int64
		db.Table("insurance_companies").Where("company_name = ? AND status = 'Active'", company).Count(&count)
		if count == 0 {
			errors["insurance_company"] = "Selected insurance company is invalid or inactive"
		}
	}

	validateNumeric := func(fieldName string, label string) {
		val, present := data[fieldName]
		if !present || val == nil || fmt.Sprintf("%v", val) == "" {
			errors[fieldName] = label + " is required"
			return
		}
		var num float64
		if _, err := fmt.Sscanf(fmt.Sprintf("%v", val), "%f", &num); err != nil {
			errors[fieldName] = label + " must be a valid number"
		} else if num < 0 {
			errors[fieldName] = label + " must be greater than or equal to 0"
		}
	}

	validateNumeric("insurance_amount", "Sum insured")
	validateNumeric("insurance_premium", "Insurance premium")

	startDateStr, okStart := data["insurance_start_date"].(string)
	endDateStr, okEnd := data["insurance_expiry_date"].(string)

	var startValid, endValid bool
	var startT, endT time.Time

	if !okStart || startDateStr == "" {
		errors["insurance_start_date"] = "Insurance start date is required"
	} else {
		var err error
		startT, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			errors["insurance_start_date"] = "Insurance start date must be a valid date (YYYY-MM-DD)"
		} else {
			startValid = true
		}
	}

	if !okEnd || endDateStr == "" {
		errors["insurance_expiry_date"] = "Insurance expiry date is required"
	} else {
		var err error
		endT, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			errors["insurance_expiry_date"] = "Insurance expiry date must be a valid date (YYYY-MM-DD)"
		} else {
			endValid = true
		}
	}

	if startValid && endValid {
		if !endT.After(startT) {
			errors["insurance_expiry_date"] = "Insurance expiry date must be after start date"
		}
	}

	return errors
}

// -----------------------------------------------------------------------------
// Step 5: Product (Lease Details) Validator
// -----------------------------------------------------------------------------
type ProductStepValidator struct{}

func (v *ProductStepValidator) Validate(db *gorm.DB, data map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	prodIDVal, ok := data["product_id"]
	if !ok || prodIDVal == nil || fmt.Sprintf("%v", prodIDVal) == "" {
		errors["product_id"] = "Product selection is required"
	} else if prodID, okParsed := parseUintVal(prodIDVal); okParsed {
		var count int64
		db.Table("products").Where("id = ? AND status = 'Active'", prodID).Count(&count)
		if count == 0 {
			errors["product_id"] = "Selected product does not exist or is inactive"
		}
	}

	period, okPeriod := data["period"]
	if !okPeriod || period == nil || fmt.Sprintf("%v", period) == "" {
		errors["period"] = "Period is required"
	} else {
		var val int
		if _, err := fmt.Sscanf(fmt.Sprintf("%v", period), "%d", &val); err != nil {
			errors["period"] = "Period must be a valid integer"
		} else if val < 1 {
			errors["period"] = "Period must be at least 1 month"
		}
	}

	validateNumeric := func(fieldName string, label string) {
		val, present := data[fieldName]
		if !present || val == nil || fmt.Sprintf("%v", val) == "" {
			errors[fieldName] = label + " is required"
			return
		}
		var num float64
		if _, err := fmt.Sscanf(fmt.Sprintf("%v", val), "%f", &num); err != nil {
			errors[fieldName] = label + " must be a valid number"
		} else if num < 0 {
			errors[fieldName] = label + " must be greater than or equal to 0"
		}
	}

	validateNumeric("interest_rate", "Interest rate")
	validateNumeric("loan_amount", "Loan amount")

	tccDate, okTcc := data["tcc_collection_date"].(string)
	if !okTcc || tccDate == "" {
		errors["tcc_collection_date"] = "First collection date is required"
	} else {
		if _, err := time.Parse("2006-01-02", tccDate); err != nil {
			errors["tcc_collection_date"] = "First collection date must be a valid date (YYYY-MM-DD)"
		}
	}

	return errors
}

// -----------------------------------------------------------------------------
// Step 6: Guarantors Validator
// -----------------------------------------------------------------------------
type GuarantorStepValidator struct{}

func (v *GuarantorStepValidator) Validate(db *gorm.DB, data map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	guarVal, ok := data["guarantors"]
	if ok && guarVal != nil {
		if guars, ok := guarVal.([]interface{}); ok {
			for idx, item := range guars {
				if gMap, ok := item.(map[string]interface{}); ok {
					custIDVal, present := gMap["customer_id"]
					if !present || custIDVal == nil || fmt.Sprintf("%v", custIDVal) == "" {
						errors[fmt.Sprintf("guarantors.%d.customer_id", idx)] = "Guarantor customer is required"
					} else if custID, okParsed := parseUintVal(custIDVal); okParsed {
						var count int64
						db.Table("customers").Where("id = ? AND deleted_at IS NULL", custID).Count(&count)
						if count == 0 {
							errors[fmt.Sprintf("guarantors.%d.customer_id", idx)] = "Selected guarantor customer does not exist"
						}
					}

					gType, present := gMap["type"]
					if !present || gType == nil || fmt.Sprintf("%v", gType) == "" {
						errors[fmt.Sprintf("guarantors.%d.type", idx)] = "Guarantor type is required"
					}
				}
			}
		}
	}

	return errors
}

// -----------------------------------------------------------------------------
// Step 7: PDC Security Validator
// -----------------------------------------------------------------------------
type PdcStepValidator struct{}

func (v *PdcStepValidator) Validate(db *gorm.DB, data map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	securityType, ok := data["pdc_security_type"].(string)
	if !ok || securityType == "" {
		errors["pdc_security_type"] = "Identification type is required"
		return errors
	}

	if securityType != "Cheque" && securityType != "CR Book" && securityType != "Deed" {
		errors["pdc_security_type"] = "Identification must be Cheque, CR Book, or Deed"
		return errors
	}

	refDetails, okRef := data["pdc_reference_details"].(string)
	if !okRef || refDetails == "" {
		errors["pdc_reference_details"] = "Reference details are required"
	}

	if securityType == "Cheque" {
		bankIDVal, okBank := data["pdc_bank_id"]
		if !okBank || bankIDVal == nil || fmt.Sprintf("%v", bankIDVal) == "" {
			errors["pdc_bank_id"] = "Bank is required"
		} else if bankID, okParsed := parseUintVal(bankIDVal); okParsed {
			var count int64
			db.Table("banks").Where("id = ? AND status = 'Active'", bankID).Count(&count)
			if count == 0 {
				errors["pdc_bank_id"] = "Selected bank does not exist or is inactive"
			}
		}

		chqDate, okChqDate := data["pdc_cheque_date"].(string)
		if !okChqDate || chqDate == "" {
			errors["pdc_cheque_date"] = "Cheque date is required"
		} else {
			if _, err := time.Parse("2006-01-02", chqDate); err != nil {
				errors["pdc_cheque_date"] = "Cheque date must be a valid date (YYYY-MM-DD)"
			}
		}

		chqNo, okChqNo := data["pdc_cheque_no"].(string)
		if !okChqNo || chqNo == "" {
			errors["pdc_cheque_no"] = "Cheque number is required"
		}

		ownership, okOwn := data["pdc_ownership"].(string)
		if !okOwn || ownership == "" {
			errors["pdc_ownership"] = "Ownership is required"
		}
	} else if securityType == "CR Book" {
		bookDate, okBookDate := data["pdc_book_date"].(string)
		if !okBookDate || bookDate == "" {
			errors["pdc_book_date"] = "Book date is required"
		} else {
			if _, err := time.Parse("2006-01-02", bookDate); err != nil {
				errors["pdc_book_date"] = "Book date must be a valid date (YYYY-MM-DD)"
			}
		}
	}

	return errors
}

// -----------------------------------------------------------------------------
// Step 8: Cheque Define Validator
// -----------------------------------------------------------------------------
type ChequeStepValidator struct{}

func (v *ChequeStepValidator) Validate(db *gorm.DB, data map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	chequesVal, ok := data["cheques"]
	if !ok || chequesVal == nil {
		errors["cheques"] = "Payout definitions are required"
		return errors
	}

	cheques, ok := chequesVal.([]interface{})
	if !ok || len(cheques) == 0 {
		errors["cheques"] = "At least one payout definition is required"
		return errors
	}

	for idx, item := range cheques {
		if cMap, ok := item.(map[string]interface{}); ok {
			validateRequired := func(fieldName string, label string) {
				val, present := cMap[fieldName]
				if !present || val == nil || fmt.Sprintf("%v", val) == "" {
					errors[fmt.Sprintf("cheques.%d.%s", idx, fieldName)] = label + " is required"
				}
			}

			validateRequired("payee_name", "Payee name")
			validateRequired("nic_br_no", "NIC/BR no.")
			validateRequired("instructions", "Instructions")
			validateRequired("bank_name", "Bank name")
			validateRequired("branch_name", "Branch name")
			validateRequired("account_number", "Account number")

			valAmount, present := cMap["payment_amount"]
			if !present || valAmount == nil || fmt.Sprintf("%v", valAmount) == "" {
				errors[fmt.Sprintf("cheques.%d.payment_amount", idx)] = "Payment amount is required"
			} else {
				var num float64
				if _, err := fmt.Sscanf(fmt.Sprintf("%v", valAmount), "%f", &num); err != nil {
					errors[fmt.Sprintf("cheques.%d.payment_amount", idx)] = "Payment amount must be a valid number"
				} else if num <= 0 {
					errors[fmt.Sprintf("cheques.%d.payment_amount", idx)] = "Payment amount must be greater than 0"
				}
			}
		}
	}

	return errors
}

// -----------------------------------------------------------------------------
// Step 9: CR & Docs Validator
// -----------------------------------------------------------------------------
type DocsStepValidator struct{}

func (v *DocsStepValidator) Validate(db *gorm.DB, data map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	validateRequired := func(fieldName string, label string) {
		val, present := data[fieldName]
		if !present || val == nil || fmt.Sprintf("%v", val) == "" {
			errors[fieldName] = label + " is required"
		}
	}

	validateRequired("cr_serial_no", "CR Serial No")
	validateRequired("url_cr_front", "CR Front Image")
	validateRequired("url_cr_back", "CR Back Image")
	validateRequired("url_invoice", "Invoice Document")

	return errors
}

// -----------------------------------------------------------------------------
// Strategy Pattern Registry
// -----------------------------------------------------------------------------
var Registry = map[string]StepValidator{
	"step-customer":      &CustomerStepValidator{},
	"step-introducers":   &IntroducerStepValidator{},
	"step-vehicle":       &VehicleStepValidator{},
	"step-insurance":     &InsuranceStepValidator{},
	"step-lease-details": &ProductStepValidator{},
	"step-guarantors":    &GuarantorStepValidator{},
	"step-pdc-security":  &PdcStepValidator{},
	"step-cheque-define": &ChequeStepValidator{},
	"step-documents":     &DocsStepValidator{},
}

func ValidateStep(db *gorm.DB, stepName string, data map[string]interface{}) map[string]string {
	val, ok := Registry[stepName]
	if !ok {
		return nil
	}
	return val.Validate(db, data)
}
