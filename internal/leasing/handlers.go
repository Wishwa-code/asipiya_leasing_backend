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
	lookupCtrl := &controllers.LookupController{DB: database.DB}
	leaseCtrl := &controllers.LeaseController{DB: database.DB}
	supplierCtrl := &controllers.SupplierController{DB: database.DB}

	leasing := rg.Group("/leasing")
	{
		products := leasing.Group("/products")
		{
			products.GET("", productCtrl.Index)                // GET  /api/leasing/products
			products.POST("", productCtrl.Store)               // POST /api/leasing/products
			products.GET("/:id", productCtrl.Get)              // GET  /api/leasing/products/:id
			products.PUT("/:id", productCtrl.Update)           // PUT  /api/leasing/products/:id
			products.GET("/:id/items", productCtrl.GetItems)   // GET  /api/leasing/products/:id/items
			products.POST("/status", productCtrl.UpdateStatus) // POST /api/leasing/products/status
		}

		leasing.POST("/calculate", leaseCtrl.CalculateSummary) // POST /api/leasing/calculate
	}

	// Supplier routes
	suppliers := rg.Group("/suppliers")
	{
		suppliers.GET("", supplierCtrl.Index)          // GET /api/suppliers
		suppliers.POST("", supplierCtrl.Store)         // POST /api/suppliers
		suppliers.PUT("/:id", supplierCtrl.Update)     // PUT /api/suppliers/:id
		suppliers.DELETE("/:id", supplierCtrl.Destroy) // DELETE /api/suppliers/:id
	}

	// Lookup routes
	lookups := rg.Group("/lookup")
	{
		lookups.GET("/banks", lookupCtrl.GetBanks)
		lookups.GET("/insurance-companies", lookupCtrl.GetInsuranceCompanies)
		lookups.GET("/vehicle-types", lookupCtrl.GetVehicleTypes)
		lookups.GET("/marketing-executives", lookupCtrl.GetMarketingExecutives)
	}

	// Customer routes (under /api — not /api/leasing — to match the frontend expectations)
	customers := rg.Group("/customers")
	{
		customers.GET("", customerCtrl.Index)                         // GET /api/customers
		customers.GET("/generate-id", customerCtrl.GenerateID)        // GET  /api/customers/generate-id
		customers.GET("/search", customerCtrl.Search)                 // GET /api/customers/search
		customers.POST("", customerCtrl.Store)                        // POST /api/customers
		customers.GET("/:id", customerCtrl.Get)                       // GET /api/customers/:id
		customers.POST("/:id/status", customerCtrl.UpdateStatus)      // POST /api/customers/:id/status
		customers.POST("/:id/location", customerCtrl.UpdateLocation)  // POST /api/customers/:id/location
		customers.GET("/:id/bank-accounts", customerCtrl.GetBankAccounts) // GET /api/customers/:id/bank-accounts
		customers.POST("/:id/documents", customerCtrl.UploadDocument) // POST /api/customers/:id/documents
	}

	// Location helper routes
	locations := rg.Group("/locations")
	{
		locations.GET("/cities", customerCtrl.GetCities) // GET /api/locations/cities?province=X
	}
}
