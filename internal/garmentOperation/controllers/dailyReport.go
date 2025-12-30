package controllers

import (
    "fmt"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
	"garment-management-backend/internal/garmentOperation/operationModels"

)

type DailyAmountController struct {
    DB *gorm.DB
}

func (ctrl *DailyAmountController) Store(c *gin.Context) {
    var input operationModels.DailyAmountRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    tx := ctrl.DB.Begin() // Start Transaction for data integrity 🛡️

    // 1. Get or Create Daily Report
    var report operationModels.DailyReport
    today := time.Now().Truncate(24 * time.Hour)
    
    if err := tx.Where("employee_id = ? AND date = ?", input.EmployeeID, today).
        FirstOrCreate(&report, operationModels.DailyReport{
            EmployeeID: input.EmployeeID,
            Date:       today,
        }).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle daily report"})
        return
    }

    // 2. Loop through styles and create amounts
    for i, styleID := range input.Styles {
        var style operationModels.Style // Assuming Style model exists
        if err := tx.First(&style, styleID).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Style %d not found", styleID)})
            return
        }

        // 🚀 Cumulative Validation
        var existingSum int64
        tx.Model(&operationModels.DailyAmount{}).Where("style_id = ?", styleID).
            Select("COALESCE(SUM(output), 0)").Row().Scan(&existingSum)

        if int(existingSum)+input.Outputs[i] > style.Target {
            tx.Rollback()
            remaining := style.Target - int(existingSum)
            c.JSON(http.StatusConflict, gin.H{
                "error": fmt.Sprintf("Style %s has remaining capacity of %d", style.StyleCode, remaining),
            })
            return
        }

        // Efficiency Calculation
        efficiency := 0.0
        if input.WorkingMinutes[i] > 0 {
            efficiency = (float64(input.Outputs[i]) * style.SMV / input.WorkingMinutes[i]) * 100
        }

        amount := operationModels.DailyAmount{
            DailyReportID:    report.ID,
            EmployeeID:       input.EmployeeID,
            StyleID:          styleID,
            Output:           input.Outputs[i],
            WorkingMinutes:   input.WorkingMinutes[i],
            ProcessEfficiency: efficiency,
        }

        if err := tx.Create(&amount).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
            return
        }
    }

    // 3. Update Report Averages
    var allAmounts []operationModels.DailyAmount
    tx.Where("daily_report_id = ?", report.ID).Find(&allAmounts)

    totalEff := 0.0
    for _, a := range allAmounts {
        totalEff += a.ProcessEfficiency
    }

    report.EntriesCount = len(allAmounts)
    report.UserEfficiency = totalEff / float64(len(allAmounts))
    tx.Save(&report)

    tx.Commit()
    c.JSON(http.StatusCreated, gin.H{"message": "Daily amounts added successfully! 🎉"})
}

func (ctrl *DailyAmountController) Index(c *gin.Context) {
    var reports []operationModels.DailyReport
    // Simple pagination logic
    // ctrl.DB.Preload("Employee").Order("date desc").Limit(10).Find(&reports)
    ctrl.DB.Preload("DailyAmounts.Style").Preload("DailyAmounts.Employee").Order("date desc").Limit(10).Find(&reports)
    c.JSON(http.StatusOK, reports)
}

// Highlight: New method to verify style output against target
func (ctrl *DailyAmountController) VerifyStyleTarget(c *gin.Context) {
    styleID := c.Param("id")
    var style operationModels.Style
    var totalOutput int64

    // 1. Fetch the Style to get the target
    if err := ctrl.DB.First(&style, styleID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Style not found"})
        return
    }

    // 2. Sum the Output from DailyAmount for this StyleID
    ctrl.DB.Model(&operationModels.DailyAmount{}).
        Where("style_id = ?", styleID).
        Select("COALESCE(SUM(output), 0)").
        Scan(&totalOutput)

    // 3. Compare and Return Status
    isOverTarget := totalOutput >= int64(style.Target)

    c.JSON(http.StatusOK, gin.H{
        "style_code":   style.StyleCode,
        "target":       style.Target,
        "total_output": totalOutput,
        "is_completed": isOverTarget,
        "remaining":    int64(style.Target) - totalOutput,
    })
}