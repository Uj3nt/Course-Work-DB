package routes

import (
	"formula1/handlers"   // предполагаем, что логика теперь в хендлерах
	"formula1/middleware" // для проверки токенов

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	f1 := api.Group("/f1")
	{
		f1.GET("/races", handlers.GetRaces)
		f1.GET("/races/:id", handlers.GetRaceDetail)
		f1.GET("/standings/drivers", handlers.GetDriverStandings)

		f1.GET("/drivers/:id", handlers.GetDriverDetail)
		f1.GET("/drivers/:id/stats", handlers.GetDriverStats)
	}

	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	predictions := r.Group("/api/predictions")
	predictions.Use(middleware.AuthRequired()) // Применяем ко всей группе
	{
		predictions.GET("/", handlers.GetActiveQuestions)
		predictions.POST("/vote", handlers.SubmitPrediction)
	}

	userGroup := api.Group("/user")
	userGroup.Use(middleware.AuthRequired())
	{
		userGroup.GET("/profile", handlers.GetUserProfile)
	}

	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthRequired())
	admin.Use(middleware.AdminRequired())
	{
		admin.POST("/questions", handlers.CreateQuestion)
		admin.PATCH("/questions/:id/resolve", handlers.ResolveQuestion)
	}
}
