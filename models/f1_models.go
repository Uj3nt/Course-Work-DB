package models

import (
	"time"
)

type Circuit struct {
	CircuitID  uint   `gorm:"primaryKey"`
	CircuitRef string `gorm:"uniqueIndex"`
	Name       string
	Location   string
	Country    string
}

type Constructor struct {
	ConstructorID  uint   `gorm:"primaryKey"`
	ConstructorRef string `gorm:"uniqueIndex"`
	Name           string `gorm:"uniqueIndex"`
	Nationality    string
	URL            string
}

type Driver struct {
	DriverID    uint   `gorm:"primaryKey"`
	DriverRef   string `gorm:"uniqueIndex"`
	Number      *int
	Code        *string
	Forename    string
	Surname     string
	DOB         time.Time
	Nationality string
	URL         string
}

type Race struct {
	RaceID    uint `gorm:"primaryKey"`
	Year      int
	Round     int
	CircuitID uint
	Circuit   Circuit `gorm:"foreignKey:CircuitID;references:CircuitID"`
	Name      string
	Date      time.Time
	Time      *string
	URL       string
}

type Result struct {
	ResultID uint `gorm:"primaryKey;column:result_id"`
	RaceID        uint `gorm:"column:race_id"`
	DriverID      uint `gorm:"column:driver_id"`
	ConstructorID uint `gorm:"column:constructor_id"`
	Race        Race        `gorm:"foreignKey:RaceID;references:RaceID"`
	Driver      Driver      `gorm:"foreignKey:DriverID;references:DriverID"`
	Constructor Constructor `gorm:"foreignKey:ConstructorID;references:ConstructorID"`
	Number        *int    `gorm:"column:number"`
	Grid          int     `gorm:"column:grid"`
	PositionText  string  `gorm:"column:position_text"`
	PositionOrder int     `gorm:"column:position_order"`
	Points        float64 `gorm:"column:points"`
	Laps          int     `gorm:"column:laps"`
	Time          *string `gorm:"column:time"`
}

type Qualifying struct {
	QualifyId     uint `gorm:"primaryKey"`
	RaceId        uint
	Race          Race `gorm:"foreignKey:RaceId"`
	DriverId      uint
	Driver        Driver `gorm:"foreignKey:DriverId"`
	ConstructorId uint
	Constructor   Constructor `gorm:"foreignKey:ConstructorId"`
	Position      int
	Q1            *string
	Q2            *string
	Q3            *string
}

type DriverStanding struct {
	DriverStandingsID uint `gorm:"primaryKey"`
	RaceID            uint 
	DriverID          uint 
	Race   Race   `gorm:"foreignKey:RaceID;references:RaceID"`
	Driver Driver `gorm:"foreignKey:DriverID;references:DriverID"`
	Points       float64
	Position     int
	PositionText string
	Wins         int
}

type ConstructorStanding struct {
	ConstructorStandingsID uint `gorm:"primaryKey"`
	RaceID                 uint
	Race                   Race `gorm:"foreignKey:RaceID"`
	ConstructorID          uint
	Constructor            Constructor `gorm:"foreignKey:ConstructorID"`
	Points                 float64
	Position               int
	PositionText           string
	Wins                   int
}

type LapTime struct {
	RaceID       uint `gorm:"primaryKey"` 
	DriverID     uint `gorm:"primaryKey"` 
	Lap          int  `gorm:"primaryKey"` 
	Position     int
	Time         string
	Milliseconds int
}
