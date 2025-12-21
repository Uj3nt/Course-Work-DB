package handlers

import (
	"formula1/database"
	"formula1/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetActiveQuestions(c *gin.Context) {
	var questions []models.Question
	if err := database.DB.Preload("Options").Where("is_closed = ?", false).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении вопросов"})
		return
	}
	c.JSON(http.StatusOK, questions)
}

func SubmitPrediction(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input struct {
		QuestionID uint `json:"question_id" binding:"required"`
		OptionID   uint `json:"option_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	var question models.Question
	if err := database.DB.First(&question, input.QuestionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Вопрос не найден"})
		return
	}
	if question.IsClosed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Прием ответов на этот вопрос завершен"})
		return
	}

	prediction := models.UserPrediction{
		UserID:     userID,
		QuestionID: input.QuestionID,
		OptionID:   input.OptionID,
	}

	if err := database.DB.Save(&prediction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить прогноз"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ваш голос принят!"})
}

func CreateQuestion(c *gin.Context) {
	var input struct {
		Title    string   `json:"title" binding:"required"`
		Text     string   `json:"text" binding:"required"`
		RewardXP int      `json:"reward_xp"`
		Options  []string `json:"options" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Заполните все поля, минимум 2 варианта"})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		question := models.Question{
			Title:    input.Title,
			Text:     input.Text,
			RewardXP: input.RewardXP,
		}

		if err := tx.Create(&question).Error; err != nil {
			return err
		}

		for _, optText := range input.Options {
			option := models.Option{
				QuestionID: question.QuestionID,
				Text:       optText,
			}
			if err := tx.Create(&option).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании вопроса"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Опрос успешно создан"})
}

func ResolveQuestion(c *gin.Context) {
	questionID := c.Param("id")
	var input struct {
		CorrectOptionID uint `json:"correct_option_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Укажите ID правильного ответа"})
		return
	}

	var question models.Question
	if err := database.DB.First(&question, questionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Вопрос не найден"})
		return
	}

	database.DB.Model(&question).Updates(models.Question{
		IsClosed:        true,
		CorrectOptionID: &input.CorrectOptionID,
	})

	var winnersCount int64

	result := database.DB.Model(&models.UserPrediction{}).
		Where("question_id = ? AND option_id = ? AND is_rewarded = ?", questionID, input.CorrectOptionID, false).
		Find(&[]models.UserPrediction{})

	winnersCount = result.RowsAffected

	if winnersCount > 0 {
		database.DB.Model(&models.User{}).
			Where("user_id IN (SELECT user_id FROM user_predictions WHERE question_id = ? AND option_id = ?)", questionID, input.CorrectOptionID).
			Update("experience", gorm.Expr("experience + ?", question.RewardXP))

		database.DB.Model(&models.UserPrediction{}).
			Where("question_id = ? AND option_id = ?", questionID, input.CorrectOptionID).
			Update("is_rewarded", true)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Результаты подведены!",
		"winners_count": winnersCount,
		"xp_rewarded":   question.RewardXP,
	})
}

func GetUserProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	var rank int64
	database.DB.Model(&models.User{}).Where("experience > ?", user.Experience).Count(&rank)
	rank += 1

	c.JSON(http.StatusOK, gin.H{
		"username":   user.Username,
		"experience": user.Experience,
		"rank":       rank,
		"is_admin":   user.IsAdmin,
	})
}
