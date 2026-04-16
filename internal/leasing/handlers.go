package leasing

import (
	"garment-management-backend/internal/database"
	"garment-management-backend/internal/leasing/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the leasing specific routes
func RegisterRoutes(rg *gin.RouterGroup) {
	productCtrl := &controllers.ProductController{DB: database.DB}
	customerCtrl := &controllers.CustomerController{DB: database.DB}

	leasing := rg.Group("/leasing")
	{
		products := leasing.Group("/products")
		{
			products.GET("", productCtrl.Index)                // GET  /api/leasing/products
			products.POST("", productCtrl.Store)               // POST /api/leasing/products
			products.GET("/:id", productCtrl.Get)              // GET  /api/leasing/products/:id
			products.POST("/status", productCtrl.UpdateStatus) // POST /api/leasing/products/status
		}
	}

	// Customer routes (under /api — not /api/leasing — to match the frontend expectations)
	customers := rg.Group("/customers")
	{
		customers.GET("/generate-id", customerCtrl.GenerateID)   // GET  /api/customers/generate-id
		customers.POST("", customerCtrl.Store)                    // POST /api/customers
		customers.POST("/:id/documents", customerCtrl.UploadDocument) // POST /api/customers/:id/documents
	}

	// Location helper routes
	locations := rg.Group("/locations")
	{
		locations.GET("/cities", customerCtrl.GetCities) // GET /api/locations/cities?province=X
	}
}
