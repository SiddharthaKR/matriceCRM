package routes

import (
	controller "github.com/SiddharthaKR/golang-jwt-project/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
	incomingRoutes.POST("customers/signup", controller.CustomerSignup()) // Customer signup route
	incomingRoutes.POST("customers/login", controller.CustomerLogin())   // Customer login route
}
