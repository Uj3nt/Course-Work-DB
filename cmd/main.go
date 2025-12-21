package main

import (
	"formula1/database"
	"formula1/models"
	"formula1/routes"
	"formula1/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	utils.StartBackupScheduler()

	database.DB.AutoMigrate(
		&models.User{}, &models.Question{}, &models.UserPrediction{}, &models.Option{},
		&models.Circuit{}, &models.Constructor{}, &models.Driver{}, &models.Race{},
		&models.Qualifying{},
		&models.LapTime{},
	)
	database.DB.AutoMigrate(
		&models.ConstructorStanding{},
		&models.DriverStanding{},
		&models.Result{},
	)

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
	}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	routes.SetupRoutes(r)

	r.Run(":8080")
}
