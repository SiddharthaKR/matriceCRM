package routes

import (
	controller "github.com/SiddharthaKR/golang-jwt-project/controllers"
	"github.com/SiddharthaKR/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

///// will contain basic meets, call etc....

func CustomerRoutes(incomingRoutes *gin.Engine) {
	// Apply authentication middleware to all customer routes
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/all-customers", controller.GetAllCustomers())

	incomingRoutes.GET("/company/:company_id/customers", controller.GetCustomersByCompany())
	incomingRoutes.GET("/company/:company_id/customers/:customer_id", controller.GetCompanyCustomerByID())
	incomingRoutes.PUT("/company/:company_id/customers/:customer_id", controller.UpdateCompanyCustomerByID())
	incomingRoutes.DELETE("/company/:company_id/customers/:customer_id", controller.DeleteComapnyCustomerByID())

	incomingRoutes.GET("/customer/:user_id", controller.GetCustomer())
	// Define routes for customer operations
	// incomingRoutes.POST("/customers", controller.CreateCustomer())            // Create a new customer
	// incomingRoutes.PUT("/customers/:customer_id", controller.UpdateCustomer()) // Update an existing customer
	// incomingRoutes.DELETE("/customers/:customer_id", controller.DeleteCustomer()) // Delete a customer by ID
}
