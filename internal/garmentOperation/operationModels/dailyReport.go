package operationModels

import (
    // "errors"
    "time"
    "gorm.io/gorm"
)

type DailyReport struct {
    gorm.Model
    EmployeeID     uint           `json:"employee_id"`
    Date           time.Time      `json:"date"`
    UserEfficiency float64        `json:"user_efficiency"`
    EntriesCount   int            `json:"entries_count"`
    TotalWorkingMinutes float64   `json:"total_working_minutes"`
    DailyAmounts   []DailyAmount  `gorm:"foreignKey:DailyReportID"`
}

type DailyAmount struct {
    gorm.Model
    DailyReportID    uint    `json:"daily_report_id"`
    EmployeeID       uint    `json:"employee_id"`
	Employee         Employee `gorm:"foreignKey:EmployeeID" json:"employee"`
    StyleID          uint    `json:"style_id"`
	Style            Style    `gorm:"foreignKey:StyleID" json:"style"`
    Output           int     `json:"output"`
    SMV              float64 `json:"smv"`    
    WorkingMinutes   float64 `json:"working_minutes"`
    ProcessEfficiency float64 `json:"process_efficiency"`
}

// DailyAmountRequest handles the bulk array input from your Django logic
type DailyAmountRequest struct {
    EmployeeID     uint      `json:"employee_id" binding:"required"`
    TotalWorkingMinutes float64   `json:"total_working_minutes" binding:"required"`
    Styles         []uint    `json:"style" binding:"required"`
    Outputs        []int     `json:"output" binding:"required"`
    WorkingMinutes []float64 `json:"working_minutes" binding:"required"`
}

type Style struct {
    gorm.Model 
    StyleCode         string  `gorm:"uniqueIndex:idx_style_category_unique;not null" json:"style_code"`
    StyleSerialNo      string  `json:"style_serial_no"`
    StyleCategoryCode string  `gorm:"uniqueIndex:idx_style_category_unique;not null" json:"style_category_code"`
    StylePart          string  `json:"style_part"`
    StyleOperationName string  `json:"style_operation_name"`
    StyleCategoryName  string  `json:"style_category_name"`
    SMV                float64 `json:"smv"`
    Target             int     `json:"target"`
}

type StyleRequest struct {
    StyleSerialNo      string  `json:"style_serial_no"`
    StyleCode          string  `json:"style_code" binding:"required"`
    StyleCategoryCode  string  `json:"style_category_code" binding:"required"` 
    StylePart          string  `json:"style_part"`
    StyleOperationName string  `json:"style_operation_name"`
    StyleCategoryName  string  `json:"style_category_name"`
    SMV                float64 `json:"smv" binding:"required"`
    Target             int     `json:"target"`
}