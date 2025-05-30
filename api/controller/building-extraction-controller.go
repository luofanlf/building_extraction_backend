package controller

import (
	"building-extraction/api/dto"
	"building-extraction/internal/model"
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

	//2. 获得提取用户信息
	// 从 JWT token 中获取用户信息
	username, exists := ctx.Get("username")
	if !exists {
		dto.FailWithMessage("unauthorized", ctx)
		return
	}

	// 获取用户信息
	user, err := c.service.GetUserInfo(username.(string))
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	fileURL, inputImage, outputImage, err := c.service.ExtractBuildings(file, header, user.ID)
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
	// 从 JWT token 中获取用户信息
	username, exists := ctx.Get("username")
	if !exists {
		dto.FailWithMessage("unauthorized", ctx)
		return
	}

	// 获取用户信息
	user, err := c.service.GetUserInfo(username.(string))
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	projects, err := c.service.GetAllProjects(user.ID)
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

	// 从 JWT token 中获取用户信息
	username, exists := ctx.Get("username")
	if !exists {
		dto.FailWithMessage("unauthorized", ctx)
		return
	}

	// 获取用户信息
	user, err := c.service.GetUserInfo(username.(string))
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	err = c.service.SaveProject(saveProjectRequest, user.ID)
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

func (c *BuildingExtractionController) HandleDeleteProject(ctx *gin.Context) {
	projectId := ctx.Param("id")
	err := c.service.DeleteProject(projectId)
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: "delete project successful",
	})
}

func (c *BuildingExtractionController) HandleGetUserProfile(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		dto.FailWithMessage("unauthorized", ctx)
		return
	}

	// 获取用户信息
	user, err := c.service.GetUserInfo(username.(string))
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: dto.UserToResponse(user),
	})

}

func (c *BuildingExtractionController) HandleUpdatePassword(ctx *gin.Context) {
	var updatePasswordRequest dto.UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&updatePasswordRequest); err != nil {
		dto.FailWithMessage("bind update password request failed", ctx)
		c.logger.Error("bind update password request failed: ", zap.Error(err))
		return
	}

	username, exists := ctx.Get("username")
	if !exists {
		dto.FailWithMessage("unauthorized", ctx)
		c.logger.Error("unauthorized")
		return
	}

	// 获取用户信息
	user, err := c.service.GetUserInfo(username.(string))
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		c.logger.Error("get user info failed: ", zap.Error(err))
		return
	}

	err = c.service.UpdatePassword(user, updatePasswordRequest.CurrentPassword, updatePasswordRequest.NewPassword)
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		c.logger.Error("update password failed: ", zap.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: "update password successfully",
	})
}

func (c *BuildingExtractionController) HandleGetAdminStats(ctx *gin.Context) {
	// 检查用户是否是管理员
	username, exists := ctx.Get("username")
	if !exists {
		dto.FailWithMessage("unauthorized", ctx)
		return
	}

	user, err := c.service.GetUserInfo(username.(string))
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	if user.UserRole != 1 { // 假设 1 是管理员角色
		dto.FailWithMessage("permission denied", ctx)
		return
	}

	stats, err := c.service.GetAdminStats()
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: stats,
	})
}

func (c *BuildingExtractionController) HandleCreateRequest(ctx *gin.Context) {
	var createRequest dto.CreateRequestRequest
	if err := ctx.ShouldBindJSON(&createRequest); err != nil {
		dto.FailWithMessage("bind request failed", ctx)
		return
	}

	username, exists := ctx.Get("username")
	if !exists {
		dto.FailWithMessage("unauthorized", ctx)
		return
	}

	user, err := c.service.GetUserInfo(username.(string))
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	err = c.service.CreateExtractionRequest(user.ID, createRequest.RequestCount, createRequest.Reason)
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: "request created successfully",
	})
}

func (c *BuildingExtractionController) HandleGetPendingRequests(ctx *gin.Context) {
	// 检查是否是管理员
	username, exists := ctx.Get("username")
	if !exists {
		dto.FailWithMessage("unauthorized", ctx)
		return
	}

	user, err := c.service.GetUserInfo(username.(string))
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	if user.UserRole != 1 {
		dto.FailWithMessage("permission denied", ctx)
		return
	}

	requests, err := c.service.GetPendingRequests()
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: requests,
	})
}

func (c *BuildingExtractionController) HandleRequest(ctx *gin.Context) {
	var handleRequest dto.HandleRequestRequest
	if err := ctx.ShouldBindJSON(&handleRequest); err != nil {
		dto.FailWithMessage("bind request failed", ctx)
		return
	}

	// 检查是否是管理员
	username, exists := ctx.Get("username")
	if !exists {
		dto.FailWithMessage("unauthorized", ctx)
		return
	}

	user, err := c.service.GetUserInfo(username.(string))
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	if user.UserRole != 1 {
		dto.FailWithMessage("permission denied", ctx)
		return
	}

	err = c.service.HandleRequest(handleRequest.RequestID, model.RequestStatus(handleRequest.Status), handleRequest.Reply)
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: "request handled successfully",
	})
}

func (c *BuildingExtractionController) HandleGetAllUsers(ctx *gin.Context) {
	// 检查是否是管理员
	username, exists := ctx.Get("username")
	if !exists {
		dto.FailWithMessage("unauthorized", ctx)
		return
	}

	user, err := c.service.GetUserInfo(username.(string))
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	if user.UserRole != 1 {
		dto.FailWithMessage("permission denied", ctx)
		return
	}

	users, err := c.service.GetAllUsers()
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: users,
	})
}

func (c *BuildingExtractionController) HandleGetUserRequests(ctx *gin.Context) {
	// 获取当前登录用户信息
	username, exists := ctx.Get("username")
	if !exists {
		dto.FailWithMessage("unauthorized", ctx)
		return
	}

	user, err := c.service.GetUserInfo(username.(string))
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	requests, err := c.service.GetUserRequests(user.ID)
	if err != nil {
		dto.FailWithMessage(err.Error(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data: requests,
	})
}
