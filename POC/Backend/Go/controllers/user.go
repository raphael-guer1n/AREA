package controllers

import (
	"Go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var UserDB []models.User
var Counter int

func InitDataBase() {
	Counter = 1

	defaultUser := models.User{
		Id:       Counter,
		Name:     "Name",
		Email:    "Email",
		Password: "Password",
	}
	UserDB = append(UserDB, defaultUser)
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": UserDB})
}

func AddUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Counter++
	newUser := models.User{Id: Counter, Name: input.Name, Email: input.Email, Password: input.Password}
	UserDB = append(UserDB, newUser)
	c.JSON(http.StatusOK, gin.H{"data": "User added successfully"})
}

func DeleteUser(c *gin.Context) {
	var userId int
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var index int
	var doUserExist bool
	doUserExist = false
	for i, user := range UserDB {
		if user.Id == userId {
			doUserExist = true
			index = i
			break
		}
	}
	if !doUserExist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		return
	}
	UserDB = append(UserDB[:index], UserDB[index+1:]...)
	c.JSON(http.StatusOK, gin.H{"data": "User deleted successfully"})
}

func UpdateUser(c *gin.Context) {
	var userId int
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newUser.Id != userId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id does not match"})
		return
	}

	var index int
	var doUserExist bool
	for i, user := range UserDB {
		if user.Id == userId {
			index = i
			doUserExist = true
			break
		}
	}

	if !doUserExist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		return
	}
	UserDB[index] = newUser
	c.JSON(http.StatusOK, gin.H{"data": "User updated successfully"})
}
