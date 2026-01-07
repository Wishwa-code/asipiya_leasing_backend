package controllers

import (
	"fmt"
	"garment-management-backend/internal/garmentOperation/operationModels"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
			EmployeeID:          input.EmployeeID,
			Date:                today,
			TotalWorkingMinutes: input.TotalWorkingMinutes,
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
		// var existingSum int64
		// tx.Model(&operationModels.DailyAmount{}).Where("style_id = ?", styleID).
		//     Select("COALESCE(SUM(output), 0)").Row().Scan(&existingSum)

		// if int(existingSum)+input.Outputs[i] > style.Target {
		//     tx.Rollback()
		//     remaining := style.Target - int(existingSum)
		//     c.JSON(http.StatusConflict, gin.H{
		//         "error": fmt.Sprintf("Style %s has remaining capacity of %d", style.StyleCode, remaining),
		//     })
		//     return
		// }

		// Efficiency Calculation
		efficiency := 0.0
		if input.WorkingMinutes[i] > 0 {
			efficiency = (float64(input.Outputs[i]) * style.SMV / input.WorkingMinutes[i]) * 100
		}

		amount := operationModels.DailyAmount{
			DailyReportID:     report.ID,
			EmployeeID:        input.EmployeeID,
			StyleID:           styleID,
			Output:            input.Outputs[i],
			SMV:               style.SMV,
			WorkingMinutes:    input.WorkingMinutes[i],
			ProcessEfficiency: efficiency,
		}

		if err := tx.Create(&amount).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
			return
		}
	}

	// 3. Update Report Averages
	// var allAmounts []operationModels.DailyAmount
	// tx.Where("daily_report_id = ?", report.ID).Find(&allAmounts)

	// totalEff := 0.0
	// for _, a := range allAmounts {
	//     totalEff += a.ProcessEfficiency
	// }

	// report.EntriesCount = len(allAmounts)
	// report.UserEfficiency = totalEff / float64(len(allAmounts))
	// tx.Save(&report)

	// tx.Commit()
	// c.JSON(http.StatusCreated, gin.H{"message": "Daily amounts added successfully! 🎉"})
	var allAmounts []operationModels.DailyAmount
	// Preload "Style" to access SMV for precise calculation
	if err := tx.Preload("Style").Where("daily_report_id = ?", report.ID).Find(&allAmounts).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recalculate totals"})
		return
	}

	var totalEarnedMinutes float64
	var totalWorkingMinutes float64

	for _, a := range allAmounts {
		totalEarnedMinutes += float64(a.Output) * a.Style.SMV
		totalWorkingMinutes = input.TotalWorkingMinutes
	}

	if totalWorkingMinutes > 0 {
		report.UserEfficiency = (totalEarnedMinutes / totalWorkingMinutes) * 100
	} else {
		report.UserEfficiency = 0
	}

	report.EntriesCount = len(allAmounts)
	tx.Save(&report)

	tx.Commit()
}

func (ctrl *DailyAmountController) Index(c *gin.Context) {
	var reports []operationModels.DailyReport
	// Simple pagination logic
	// ctrl.DB.Preload("Employee").Order("date desc").Limit(10).Find(&reports)
	ctrl.DB.Preload("DailyAmounts.Style").Preload("DailyAmounts.Employee").Order("date desc").Find(&reports)
	c.JSON(http.StatusOK, reports)
}

func (ctrl *DailyAmountController) Update(c *gin.Context) {
	var input operationModels.DailyAmountRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id") // The DailyReport ID
	tx := ctrl.DB.Begin()

	// 1. Fetch the existing report
	var report operationModels.DailyReport
	if err := tx.First(&report, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Daily report not found"})
		return
	}

	// 2. Sync Report Level Data
	report.TotalWorkingMinutes = input.TotalWorkingMinutes
	report.EmployeeID = input.EmployeeID

	// 3. REMOVE existing amounts to replace them with new ones 🧹
	if err := tx.Where("daily_report_id = ?", report.ID).Delete(&operationModels.DailyAmount{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset records"})
		return
	}

	// 4. Create new records (Same logic as Store) 🔄
	for i, styleID := range input.Styles {
		var style operationModels.Style
		if err := tx.First(&style, styleID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Style %d not found", styleID)})
			return
		}

		efficiency := 0.0
		if input.WorkingMinutes[i] > 0 {
			efficiency = (float64(input.Outputs[i]) * style.SMV / input.WorkingMinutes[i]) * 100
		}

		amount := operationModels.DailyAmount{
			DailyReportID:     report.ID,
			EmployeeID:        input.EmployeeID,
			StyleID:           styleID,
			Output:            input.Outputs[i],
			SMV:               style.SMV,
			WorkingMinutes:    input.WorkingMinutes[i],
			ProcessEfficiency: efficiency,
		}

		if err := tx.Create(&amount).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update entry"})
			return
		}
	}

	// 5. Final Recalculation (Shared logic) 📈
	var allAmounts []operationModels.DailyAmount
	tx.Preload("Style").Where("daily_report_id = ?", report.ID).Find(&allAmounts)

	var totalEarnedMinutes float64
	for _, a := range allAmounts {
		totalEarnedMinutes += float64(a.Output) * a.Style.SMV
	}

	if input.TotalWorkingMinutes > 0 {
		report.UserEfficiency = (totalEarnedMinutes / input.TotalWorkingMinutes) * 100
	} else {
		report.UserEfficiency = 0
	}

	report.EntriesCount = len(allAmounts)
	tx.Save(&report)

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Daily report updated successfully! ✨", "report": report})
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
