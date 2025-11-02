package main

import (
	"bwastartup/auth"
	"bwastartup/handler"
	"bwastartup/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
  dsn := "root:123456@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err != nil {
	log.Fatal(err.Error())
  }

  // input dari user
  // handler, mapping input dari user -> struct input
  // service : melakukan mapping dari struct input ke entity user
  // repository : simpan entity user ke database
  
  // dependency injection
  userRepository := user.NewRepository(db)

  // service punya dependency repository
  userService := user.NewService(userRepository)

  // userService.SaveAvatar(1, "images/avatar1.png")

  // inisial jwt ke var
  authService := auth.NewService()

  // handler punya dependency service
  userHandler := handler.NewUserHandler(userService, authService)

  router := gin.Default()
  api := router.Group("/api/v1")
  api.POST("/users", userHandler.RegisterUser)
  api.POST("/sessions", userHandler.Login)
  api.POST("/email_checkers", userHandler.CheckEmailAvailability)
  api.POST("/avatars", userHandler.UploadAvatar)

  router.Run()
}