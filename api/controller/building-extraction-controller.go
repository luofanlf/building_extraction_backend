package controller

import (
	"building-extraction/api/dto"
	"building-extraction/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BuildingExtractionController struct {
	logger  *zap.Logger
	service *service.BuildingExtractionService
}

func NewBuildingExtractionController(service *service.BuildingExtractionService) *BuildingExtractionController {
	logger, _ := zap.NewProduction()

	// 调整堆栈级别为 Fatal
	logger = logger.WithOptions(zap.AddStacktrace(zapcore.FatalLevel))

	return &BuildingExtractionController{
		logger:  logger,
		service: service,
	}
}

func (c *BuildingExtractionController) HandleLogin(ctx *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		dto.FailWithMessage("bind login request failed", ctx)
		return
	}

	token, err := c.service.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: gin.H{
			"token":   token,
			"message": "login successful",
		},
	})
}

func (c *BuildingExtractionController) HandleRegister(ctx *gin.Context) {
	//bind register parameter
	var registerRequest dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		c.logger.Error("bind register request failed: ", zap.Error(err))
		dto.FailWithMessage("bind register request failed: ", ctx)
		return
	}

	//check password and confirm password
	if registerRequest.Password != registerRequest.ConfirmPassword {
		c.logger.Error("password is not equal to confirm password")
		dto.FailWithMessage("password is not equal to confirm password", ctx)
		return
	}

	//re
	err := c.service.Register(registerRequest.Username, registerRequest.Password)
	if err != nil {
		c.logger.Error("register failed: ", zap.Error(err))
		dto.FailWithMessage("register failed: "+err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: "login successful",
	})
}

func (c *BuildingExtractionController) HandleExtraction(ctx *gin.Context) {
	// 1. 从请求中获取文件
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		dto.FailWithMessage("no file uploaded or invalid form data", ctx)
		return
	}
	defer file.Close()

	fileURL, inputImage, outputImage, err := c.service.ExtractBuildings(file, header)
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	// 返回原始图片和掩码图片的URL
	ctx.JSON(http.StatusOK, dto.Response{
		Data: dto.ExtractionResponse{
			MaskUrl:     fileURL,
			InputImage:  inputImage,
			OutputImage: outputImage,
		},
	})
}

func (c *BuildingExtractionController) HandleGetProjects(ctx *gin.Context) {
	projects, err := c.service.GetAllProjects()
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	projectResponses := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = dto.ProjectToResponse(&project)
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: projectResponses,
	})
}

func (c *BuildingExtractionController) HandleSaveProject(ctx *gin.Context) {
	var saveProjectRequest dto.SaveProjectRequest
	if err := ctx.ShouldBindJSON(&saveProjectRequest); err != nil {
		dto.FailWithMessage("bind save project request failed", ctx)
		return
	}

	err := c.service.SaveProject(saveProjectRequest)
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: "save project successful",
	})
}

func (c *BuildingExtractionController) HandleGetProjectDetail(ctx *gin.Context) {
	projectId := ctx.Param("id")
	project, err := c.service.GetProjectById(projectId)
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: dto.ProjectToResponse(project),
	})
}
