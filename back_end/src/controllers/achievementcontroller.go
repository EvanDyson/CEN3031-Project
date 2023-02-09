package controllers

import (
	"net/http"

	"CEN3031-Project/back_end/src/database"
	"CEN3031-Project/back_end/src/models"
	
	"github.com/gin-gonic/gin"
)

type AchievementRequest struct {
	Title 				string `json:"title"`
	Description 		string `json:"description"`
	ExperiencePoints 	int `json:"expPts"`
}

func GetAchievement(user *models.User, title string) {
	var achievement models.Achievement
	record := database.AchievementDB.Where("title = ?", title).First(&achievement)

	if record.Error != nil {
		//context.JSON(http.StatusNotImplemented)
		return
	}

	user.Achievements = append(user.Achievements, (int64)(achievement.ID))
}

func AddAchievement(context *gin.Context) {
	var achievement models.Achievement
	var request AchievementRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	createAchievement(&achievement, &request)

	record := database.AchievementDB.Create(&achievement)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{"title": achievement.Title, "descrip": achievement.Description, "expPts": achievement.ExperiencePoints})

}

func createAchievement(achievement *models.Achievement, request *AchievementRequest) {
	achievement.Title = request.Title
	achievement.Description = request.Description
	achievement.ExperiencePoints = request.ExperiencePoints 
}

func GetAllAchievements(context *gin.Context) {
	var achievements []models.Achievement
	database.AchievementDB.Find(&achievements) 
	context.IndentedJSON(http.StatusAccepted, achievements)
}

/*
func DeleteAchievement(context *gin.Context) {
	var title = struct {
		Title string `json:"title"`
	}{}

	if err := context.ShouldBindJSON(&title); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.

}
*/