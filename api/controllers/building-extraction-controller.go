package controllers

import (
	"building-extraction/api/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BuildingExtractionController struct {
	logger *zap.Logger
}

func NewBuildingExtractionController() *BuildingExtractionController {
	logger, _ := zap.NewProduction()

	// 调整堆栈级别为 Fatal
	logger = logger.WithOptions(zap.AddStacktrace(zapcore.FatalLevel))

	return &BuildingExtractionController{
		logger: logger,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *BuildingExtractionController) HandleLogin(ctx *gin.Context) {
	//bind loogin parameter
	var loginRequest LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		dto.FailWithMessage("bind login request failed: ", ctx)
		return
	}

	//check account and password
	if loginRequest.Username != "admin" && loginRequest.Password != "admin" {
		dto.FailWithMessage("incorrect account or password: ", ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: "login successful",
	})

}
