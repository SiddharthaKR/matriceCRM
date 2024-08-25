package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/akhil/golang-jwt-project/controllers"
	"github.com/akhil/golang-jwt-project/middleware"
)

func CompanyRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/companies", controllers.CreateCompany())
	incomingRoutes.GET("/companies", controllers.GetCompanies())
	incomingRoutes.GET("/companies/:company_id", controllers.GetCompany())
	incomingRoutes.PUT("/companies/:company_id", controllers.UpdateCompany())
	incomingRoutes.DELETE("/companies/:company_id", controllers.DeleteCompany())
}
