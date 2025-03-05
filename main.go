package main

import (
	"building-extraction/api/controller"
	"building-extraction/internal/model"
	"building-extraction/internal/service"
	"fmt"
	"log"

	"building-extraction/api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//connect to postgres
	dsn := "host=localhost user=myuser password=mypassword dbname=building_extraction port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	fmt.Println("数据库连接成功")

	// 自动迁移
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Printf("自动迁移警告: %v", err)
	} else {
		fmt.Println("自动迁移成功")
	}
	r := gin.Default()
	r.MaxMultipartMemory = 20 << 20

	//health check
	r.GET("/api/message", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})

	service := service.NewBuildingExtractionService(db)
	ctrl := controller.NewBuildingExtractionController(service)

	// API 路由组
	api := r.Group("/api")

	// 公开路由
	api.POST("/login", ctrl.HandleLogin)
	api.POST("/register", ctrl.HandleRegister)

	// 需要认证的路由
	authorized := api.Group("")
	authorized.Use(middleware.AuthMiddleware())
	{
		// authorized.POST("/extraction", ctrl.HandleExtraction)
		authorized.POST("/upload", ctrl.UploadHandler)
	}

	r.Run(":8080")
}
