package routes

import (
	"github.com/SiddharthaKR/golang-jwt-project/controllers"
	"github.com/SiddharthaKR/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func CompanyRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/companies", controllers.CreateCompany())
	incomingRoutes.GET("/companies", controllers.GetCompanies())
	incomingRoutes.GET("/companies/:company_id", controllers.GetCompany())
	incomingRoutes.PUT("/companies/:company_id", controllers.UpdateCompany())
	incomingRoutes.DELETE("/companies/:company_id", controllers.DeleteCompany())
}
