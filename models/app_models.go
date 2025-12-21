package models

import (
	"time"
)

type User struct {
	UserID     uint   `gorm:"primaryKey;column:user_id"`
	Username   string `gorm:"unique;not null"`
	Password   string `gorm:"not null"` // Здесь будем хранить хэш
	Experience int    `gorm:"column:experience;default:0"`
	CreatedAt  time.Time
	IsAdmin    bool `gorm:"column:is_admin;default:false"`
}

type Question struct {
	QuestionID      uint   `gorm:"primaryKey;column:question_id"`
	Title           string `gorm:"column:title;not null"`
	Text            string `gorm:"column:text;not null"`
	RewardXP        int    `gorm:"column:reward_xp;default:0"`
	IsClosed        bool   `gorm:"column:is_closed;default:false"`
	CorrectOptionID *uint  `gorm:"column:correct_option_id"`

	Options   []Option `gorm:"foreignKey:QuestionID"`
	CreatedAt time.Time
}

type Option struct {
	OptionID   uint   `gorm:"primaryKey;column:option_id"`
	QuestionID uint   `gorm:"column:question_id"`
	Text       string `gorm:"column:text;not null"`
}

type UserPrediction struct {
	UserID     uint `gorm:"primaryKey"`
	QuestionID uint `gorm:"primaryKey"`
	OptionID   uint `gorm:"column:option_id"`
	IsRewarded bool `gorm:"column:is_rewarded;default:false"`
	CreatedAt  time.Time
}
