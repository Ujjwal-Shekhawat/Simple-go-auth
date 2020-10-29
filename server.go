package main

import (
	"io"
	"log"
	"os"

	controller "github.com/Ujjwal-Shekhawat/golang-gin-poc/controllers"
	"github.com/Ujjwal-Shekhawat/golang-gin-poc/middleware"
	"github.com/Ujjwal-Shekhawat/golang-gin-poc/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	imgService    service.ImgService       = service.New()
	imgController controller.ImgController = controller.New(imgService)
)

func setupLogger() {
	file, err := os.Create("server.log")
	if err != nil {
		log.Fatal("Error setting up logger")
	}

	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
}

func main() {
	error := godotenv.Load(".env")
	if error != nil {
		log.Fatal(error)
	}

	setupLogger()
	middleware.Init()

	server := gin.New()

	// Experemental group
	rGrp := server.Group("/auth")
	rGrp.Use(middleware.Logger())
	rGrp.GET("/login", middleware.Protect(), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "loggedin",
		})
	})

	rGrp.GET("/regester", middleware.Regester(), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "User regestered",
		})
	})

	server.Use(gin.Recovery(), middleware.Logger())

	server.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Hello from gin server"})
	})

	server.GET("/images", func(ctx *gin.Context) {
		ctx.JSON(200, imgController.FindAll())
	})

	server.POST("/images", middleware.Protect(), func(ctx *gin.Context) {
		// Later shift it to secure middleware
		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.JSON(200, gin.H{"message": imgController.Save(ctx)})
	})

	server.Run(":8080")
}
