package main

import (
	"Go/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	controllers.InitDataBase()
	r := gin.Default()

	r.GET("/users", controllers.GetUsers)
	r.POST("/users", controllers.AddUser)
	r.DELETE("/users/:id", controllers.DeleteUser)
	r.PUT("/users/:id", controllers.UpdateUser)
	err := r.Run()
	if err != nil {
		return
	}
}
