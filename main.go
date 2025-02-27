package main

import (
	"building-extraction/api/controllers"
	"building-extraction/internal/model"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//connect to postgres
	dsn := "host=localhost user=root password=root dbname=building_extraction port=5432 sslmode=disable"
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

	//health check
	r.GET("/api/message", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})

	ctrl := controllers.NewBuildingExtractionController()
	r.POST("/api/login", ctrl.HandleLogin)
	r.Run(":8080")
}
