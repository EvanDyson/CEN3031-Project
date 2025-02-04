// All functions that concern behaviors of logged in users
package controllers

import (
	"fmt"
	"net/http"

	"AttackOnCollege/back_end/src/database"
	"AttackOnCollege/back_end/src/models"

	"github.com/gin-gonic/gin"
)

func Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}

type ProfileRequest struct {
	Username      string `form:"username" gorm:"unique"`
	LastName      string `form:"lastName"`
	Email         string `form:"email" gorm:"unique"`
	Major         string `form:"major" binding:"required"`
	College       string `form:"college" binding:"required"`
	DOB           string `form:"dob"`
	FirstName     string `form:"firstName"`
	Age           int    `form:"age"`
	CurrentCourse string `form:"currCourse"`
}

type AchievementRes struct {
	Title            string `form:"title"`
	Description      string `form:"description"`
	ExperiencePoints int    `form:"XPgain"`
}

type AssignmentRes struct {
  ID      uint   `form:"assignmentID"`
  Title   string `form:"assignmentName"`
  DueDate string `form:"dueDate"`
  Course  string `form:"courseCode"`
}

// Return user information to be displayed on the user profile
func GetUser(context *gin.Context) {
	var user models.User
	// Get the authorization header from the request
	tokenString := context.GetHeader("Authorization")

	// Store information about the user with given token
	tokenString = tokenString[1 : len(tokenString)-1]
	record := database.UserDB.Where("token = ?", tokenString).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	var request ProfileRequest
	// Send back only information needed by the front end for now.
	mapUserToRequest(&user, &request)
  var response []interface{}
  response = append(response, request)
  assignments, size := GetAssignmentHelper(&user)
  response = append(response, size)
  response = append(response, assignments)

	context.JSON(http.StatusOK, response)
}

func GetAssignmentHelper(user *models.User) ([]AssignmentRes, int) {
  var assignments []AssignmentRes
  for i := range user.Assignments {
    var assignment models.Assignment
    record := database.AssignmentDB.Where("id = ?", user.Assignments[i]).First(&assignment)
    if record.Error != nil {
      fmt.Printf(record.Error.Error())
      break
    }

    var course models.Course
    check := database.CourseDB.Where("id = ?", user.CourseID).First(&course)
    if check.Error != nil {
      fmt.Printf(check.Error.Error())
      break
    }

    assignments = append(assignments, AssignmentRes{ID: assignment.ID,
      Title: assignment.Title,
      DueDate: assignment.DueDate,
      Course: course.CourseCode})

  }

  return assignments, len(assignments)
}

// Modifies the information in the user profile
func EditUser(context *gin.Context) {
	var user models.User
	tokenString := context.GetHeader("Authorization")

	tokenString = tokenString[1 : len(tokenString)-1]
	record := database.UserDB.Where("token = ?", tokenString).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	var request = struct {
		Username  string `form:"username" gorm:"unique"`
		FirstName string `form:"firstName"`
		LastName  string `form:"lastName"`
	}{}
	if err := context.Bind(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	user.Username = request.Username
	user.FirstName = request.FirstName
	user.LastName = request.LastName

	database.UserDB.Save(&user)
	user.Password = "Hidden"
	user.Token = ""
	context.JSON(http.StatusAccepted, user)
}

func ChangePassword(context *gin.Context) {
	var user models.User
	tokenString := context.GetHeader("Authorization")

	tokenString = tokenString[1 : len(tokenString)-1]
	record := database.UserDB.Where("token = ?", tokenString).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	var request = struct {
		Password string `form:"password"`
	}{}
	if err := context.Bind(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if err := user.HashPassword(request.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	database.UserDB.Save(&user)
	context.JSON(http.StatusAccepted, "Successfully changed password!")
}

func DeleteUser(context *gin.Context) {
	tokenString := context.GetHeader("Authorization")
	tokenString = tokenString[1 : len(tokenString)-1]
	record := database.UserDB.Where("token = ?", tokenString).Delete(&models.User{})
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, "Succesully deleted user!")
}

func LogOut(context *gin.Context) {
	var user models.User
	tokenString := context.GetHeader("Authorization")
	tokenString = tokenString[1 : len(tokenString)-1]
	record := database.UserDB.Where("token = ?", tokenString).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	user.Token = ""
	database.UserDB.Save(&user)
}

func GetAchievements(context *gin.Context) {
	var user models.User
	// Get the authorization header from the request
	tokenString := context.GetHeader("Authorization")

	// Store information about the user with given token
	tokenString = tokenString[1 : len(tokenString)-1]
	record := database.UserDB.Where("token = ?", tokenString).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	var temp []AchievementRes
	for _, id := range user.Achievements {
		var a models.Achievement
		r := database.AchievementDB.Where("id = ?", int(id)).First(&a)
		if r.Error != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
			context.Abort()
			return
		}
		var res AchievementRes
		res.Title = a.Title
		res.Description = a.Description
		res.ExperiencePoints = a.ExperiencePoints
		temp = append(temp, res)
	}
	var achievements []interface{}
  achievements = append(achievements, len(temp))
	achievements = append(achievements, temp)
	context.JSON(http.StatusOK, achievements)
}

func mapUserToRequest(user *models.User, request *ProfileRequest) {
	request.Username = user.Username
	request.Email = user.Email
	request.Major = user.Major
	request.College = user.College
	request.DOB = user.DOB
	request.FirstName = user.FirstName
	request.LastName = user.LastName
	request.Age = user.Age
	request.CurrentCourse = user.CurrentCourse
}
