package routes

import (
	"github.com/PranavMasekar/go-jwt-project/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	// These are unprotected routes bcz user dont have any token at this time
	incomingRoutes.POST("users/signup", controllers.SignUp())
	incomingRoutes.POST("users/login", controllers.Login())
}
