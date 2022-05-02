package routes

import (
	"github.com/PranavMasekar/go-jwt-project/controllers"
	"github.com/PranavMasekar/go-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// These are protected routes so user should have token to access these routes
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
}
