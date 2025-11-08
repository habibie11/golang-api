package main

import (
	"bwastartup/auth"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
  api.POST("/avatars", authMiddleware(authService, userService),userHandler.UploadAvatar)

  router.Run()
}

// ambil nilai header Authorization: Bearer tokentokentoken
// dari header Authorization, ambil nilai tokennya saja
// validasi token
// ambil user_id
// ambil user dari db berdasarkan user_id lewat service
// set context isinya user

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
  return func (c *gin.Context) {
  authHeader := c.GetHeader("Authorization")

  if !strings.Contains(authHeader, "Bearer") {
    response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
    c.AbortWithStatusJSON(http.StatusUnauthorized, response)
    return
  }

  // dapatkan token
  tokenString := ""
  arrayToken := strings.Split(authHeader, " ")
  if len(arrayToken) == 2 {
    tokenString = arrayToken[1]
  }

  // validasi token
  token, err := authService.ValidateToken(tokenString)
  if err != nil {
    response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
    c.AbortWithStatusJSON(http.StatusUnauthorized, response)
    return
  }

  claim, ok := token.Claims.(jwt.MapClaims)

  if !ok || !token.Valid {
    response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
    c.AbortWithStatusJSON(http.StatusUnauthorized, response)
    return
  }

  userID := int(claim["user_id"].(float64))

  user, err := userService.GetUserByID(userID)
  if err != nil {
    response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
    c.AbortWithStatusJSON(http.StatusUnauthorized, response)
    return
  }

  // set user ke context
  c.Set("currentUser", user)
}
}