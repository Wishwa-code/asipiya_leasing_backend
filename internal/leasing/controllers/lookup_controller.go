package controllers

import (
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
	var types []models.VehicleType
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
