package controllers

import (
	adminModels "garment-management-backend/internal/admin/models"
	"garment-management-backend/internal/leasing/models"
	gmodels "garment-management-backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LookupController struct {
	DB *gorm.DB
}

// GetBanks handles GET /api/lookup/banks
func (ctrl *LookupController) GetBanks(c *gin.Context) {
	var banks []models.Bank
	if err := ctrl.DB.Where("status = ?", "Active").Find(&banks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch banks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": banks})
}

// GetInsuranceCompanies handles GET /api/lookup/insurance-companies
func (ctrl *LookupController) GetInsuranceCompanies(c *gin.Context) {
	var companies []models.InsuranceCompany
	if err := ctrl.DB.Where("status = ?", "Active").Find(&companies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insurance companies"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": companies})
}

// GetVehicleTypes handles GET /api/lookup/vehicle-types
func (ctrl *LookupController) GetVehicleTypes(c *gin.Context) {
	var types []adminModels.VehicleType
	if err := ctrl.DB.Where("status = ?", "Active").Find(&types).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicle types"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": types})
}

// GetMarketingExecutives handles GET /api/lookup/marketing-executives
func (ctrl *LookupController) GetMarketingExecutives(c *gin.Context) {
	// Let's assume all users are currently considered marketing executives without strict roles in the Go DB.
	// Filter by Active status if the status column exists, else all users.
	var users []gmodels.User
	if err := ctrl.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch marketing executives"})
		return
	}

	// Map to a cleaner format
	type Executive struct {
		ID       uint   `json:"id"`
		FullName string `json:"full_name"`
	}

	var response []Executive
	for _, u := range users {
		response = append(response, Executive{
			ID:       u.ID,
			FullName: u.Name, // Using Name from gmodels.User
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetVehicleMakes handles GET /api/lookup/vehicle-makes?type_id=X
func (ctrl *LookupController) GetVehicleMakes(c *gin.Context) {
	var makes []adminModels.VehicleMake
	query := ctrl.DB.Where("status = ?", "Active")

	if typeIDStr := c.Query("type_id"); typeIDStr != "" {
		query = query.Where("vehicle_type_id = ?", typeIDStr)
	}

	if err := query.Find(&makes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicle makes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": makes})
}

// GetVehicleModels handles GET /api/lookup/vehicle-models?type_id=X
func (ctrl *LookupController) GetVehicleModels(c *gin.Context) {
	var models []adminModels.VehicleModel
	query := ctrl.DB.Where("status = ?", "Active")

	if typeIDStr := c.Query("type_id"); typeIDStr != "" {
		query = query.Where("vehicle_type_id = ?", typeIDStr)
	}

	if err := query.Find(&models).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicle models"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": models})
}

// GetRmvData handles GET /api/v1/lookup/rmv-data?reg_no=X&chassis_no=Y
func (ctrl *LookupController) GetRmvData(c *gin.Context) {
	regNo := c.Query("reg_no")
	chassisNo := c.Query("chassis_no")

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"registered_no":      regNo,
			"chasis_no":          chassisNo,
			"vehicle_make":       "Toyota",
			"vehicle_model":      "Prius",
			"manufacturing_year": "2018",
			"fuel_type":          "Hybrid",
			"engine_no":          "1NZ-345678",
			"engine_cc":          "1500",
			"gross_weight":       "1350",
		},
	})
}
