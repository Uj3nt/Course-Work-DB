package handlers

import (
	"net/http"

	"formula1/database"
	"formula1/models"

	"github.com/gin-gonic/gin"
)

func GetRaces(c *gin.Context) {
	var races []models.Race
	year := c.Query("year")

	query := database.DB.Preload("Circuit").Preload("Race.Circuit")

	if year != "" {
		query = query.Where("year = ?", year)
	}

	query.Order("round asc").Find(&races)
	c.JSON(http.StatusOK, races)
}

func GetRaceDetail(c *gin.Context) {
	id := c.Param("id")
	var race models.Race
	var results []models.Result

	if err := database.DB.Preload("Circuit").First(&race, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Race not found"})
		return
	}

	database.DB.
		Preload("Driver").
		Preload("Constructor").
		Preload("Race").
		Preload("Race.Circuit").
		Where("race_id = ?", id).
		Order("position_order asc").
		Find(&results)

	c.JSON(http.StatusOK, gin.H{
		"race":    race,
		"results": results,
	})
}

func GetDriverStandings(c *gin.Context) {
	year := c.Query("year")
	var standings []models.DriverStanding

	subQuery := database.DB.Table("races").
		Select("MAX(race_id)").
		Where("year = ?", year)

	err := database.DB.Debug().
		Preload("Driver").       
		Preload("Race").        
		Preload("Race.Circuit"). 
		Where("race_id = (?)", subQuery).
		Order("position asc").
		Find(&standings).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, standings)
}

func GetDriverDetail(c *gin.Context) {
	id := c.Param("id")
	var driver models.Driver

	if err := database.DB.First(&driver, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Гонщик не найден"})
		return
	}

	c.JSON(http.StatusOK, driver)
}

func GetDriverStats(c *gin.Context) {
	driverID := c.Param("id")
	year := c.Query("year") 

	type Stats struct {
		Wins        int64   `json:"wins"`
		Poles       int64   `json:"poles"`
		TotalPoints float64 `json:"total_points"`
	}
	var stats Stats

	winsQuery := database.DB.Model(&models.Result{}).Where("driver_id = ? AND position_order = 1", driverID)
	if year != "" {
		winsQuery = winsQuery.Joins("JOIN races ON races.race_id = results.race_id").Where("races.year = ?", year)
	}
	winsQuery.Count(&stats.Wins)

	polesQuery := database.DB.Model(&models.Qualifying{}).Where("driver_id = ? AND position = 1", driverID)
	if year != "" {
		polesQuery = polesQuery.Joins("JOIN races ON races.race_id = qualifying.race_id").Where("races.year = ?", year)
	}
	polesQuery.Count(&stats.Poles)

	pointsQuery := database.DB.Model(&models.Result{}).Select("SUM(points)").Where("driver_id = ?", driverID)
	if year != "" {
		pointsQuery = pointsQuery.Joins("JOIN races ON races.race_id = results.race_id").Where("races.year = ?", year)
	}
	pointsQuery.Scan(&stats.TotalPoints)

	c.JSON(http.StatusOK, stats)
}

func GetDriverSeasons(c *gin.Context) {
	driverID := c.Param("id")
	var years []int

	database.DB.Table("results").
		Joins("JOIN races ON races.race_id = results.race_id").
		Where("results.driver_id = ?", driverID).
		Distinct("races.year").
		Order("races.year desc").
		Pluck("races.year", &years)

	c.JSON(http.StatusOK, years)
}
