package service

import (
	"building-extraction/api/dto"
	"building-extraction/internal/dao"
	"building-extraction/internal/model"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"building-extraction/pkg/auth"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type BuildingExtractionService struct {
	db     *gorm.DB
	logger *zap.Logger
	dao    *dao.BuildingExtractionDao
}

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	emailRegex    = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	passwordRegex = regexp.MustCompile(`^[A-Za-z\d]{8,}$`)
)

func NewBuildingExtractionService(db *gorm.DB) *BuildingExtractionService {
	logger, _ := zap.NewProduction()

	// 调整堆栈级别为 Fatal
	logger = logger.WithOptions(zap.AddStacktrace(zapcore.FatalLevel))

	return &BuildingExtractionService{
		logger: logger,
		db:     db,
		dao:    dao.NewBuildingExtractionDao(),
	}
}

func (s *BuildingExtractionService) validateUsernameOrEmail(input string) error {
	if usernameRegex.MatchString(input) || emailRegex.MatchString(input) {
		return nil
	}
	return fmt.Errorf("输入格式必须是用户名(3-20位字母数字下划线)或有效的邮箱地址")
}

func (s *BuildingExtractionService) Login(username string, password string) (string, error) {
	user, err := s.dao.GetUserByUsername(s.db, username)
	if err != nil {
		s.logger.Error("get user by username failed", zap.Error(err))
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("username %s not found", username)
	}
	if user.Password != password {
		return "", fmt.Errorf("incorrect password")
	}

	// 生成 token
	token, err := auth.GenerateToken(username)
	s.logger.Info("generate token", zap.String("token", token))
	if err != nil {
		s.logger.Error("generate token failed", zap.Error(err))
		return "", fmt.Errorf("generate token failed")
	}

	return token, nil
}

func (s *BuildingExtractionService) Register(username string, password string) error {
	registerUser := model.User{
		Username:       username,
		Password:       password,
		RemainingCount: 10,
	}

	// 验证密码：至少8个字符，包含至少1个字母和1个数字
	if !passwordRegex.MatchString(password) ||
		!strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"+
			"abcdefghijklmnopqrstuvwxyz") ||
		!strings.ContainsAny(password, "0123456789") {
		return fmt.Errorf("password must be at least 8 characters, contain at least 1 letter and 1 number")
	}

	err := s.validateUsernameOrEmail(username)
	if err != nil {
		s.logger.Error("validate username or email failed", zap.Error(err))
		return err
	}

	//check if username exists
	exist, err := s.dao.CheckUserExists(s.db, registerUser.Username)
	if err != nil {
		s.logger.Error("check user exists failed", zap.Error(err))
		return err
	}
	if exist {
		return fmt.Errorf("username %s already exists", username)
	}

	//create user
	err = s.dao.CreateUser(s.db, &registerUser)
	if err != nil {
		s.logger.Error("create user failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *BuildingExtractionService) GetUserInfo(username string) (*model.User, error) {
	user, err := s.dao.GetUserByUsername(s.db, username)
	if err != nil {
		s.logger.Error("get user info failed", zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// 不返回密码
	user.Password = ""
	return user, nil
}

func (s *BuildingExtractionService) ExtractBuildings(file multipart.File, header *multipart.FileHeader, userID int) (string, string, string, error) {
	//校验用户是否还有剩余提取次数
	count, err := s.dao.CheckRemainingCount(s.db, userID)
	if err != nil {
		s.logger.Error("check remaining count failed", zap.Error(err))
		return "", "", "", fmt.Errorf("building extraction failed: %v", err)
	}
	if count <= 0 {
		s.logger.Error("user has no remaining extraction count")
		return "", "", "", fmt.Errorf("user has no remaining extraction count")
	}
	//上传图片
	if !isValidImageType(header.Filename) {
		s.logger.Error("invalid file type", zap.Error(fmt.Errorf("invalid file type. only image files are allowed")))
		return "", "", "", fmt.Errorf("invalid file type. only image files are allowed")
	}
	// 3. 验证文件大小 (限制为20MB)
	if header.Size > 20*1024*1024 {
		s.logger.Error("file size exceeds the limit (20MB)")
		return "", "", "", fmt.Errorf("file size exceeds the limit (20MB)")
	}
	// 4. 为文件生成唯一文件名
	filename := generateUniqueFilename(header.Filename)
	// 5. 确保上传目录存在
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		s.logger.Error("failed to create upload directory", zap.Error(err))
		return "", "", "", fmt.Errorf("failed to create upload directory")
	}
	// 6. 创建目标文件
	filepath := path.Join(uploadDir, filename)
	out, err := os.Create(filepath)
	if err != nil {
		s.logger.Error("failed to create file on server", zap.Error(err))
		return "", "", "", fmt.Errorf("failed to create file on server")
	}
	defer out.Close()
	// 7. 将上传的文件内容复制到目标文件
	_, err = io.Copy(out, file)
	if err != nil {
		s.logger.Error("failed to save file on server", zap.Error(err))
		return "", "", "", fmt.Errorf("failed to save file on server")
	}

	// 3. 生成输出文件名（注意：这里只要 .png，就留给 Python 来补 _mask.png）
	resultsDir := "./results"
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		s.logger.Error("failed to create results directory", zap.Error(err))
		return "", "", "", fmt.Errorf("failed to create results directory")
	}
	// 假设你想以同一个名字结尾为 .png
	outputFilename := strings.TrimSuffix(filename, path.Ext(filename)) + ".png"
	outputFilepath := path.Join(resultsDir, outputFilename)

	// 4. 调用 Python
	cmd := exec.Command("python", "ml_server/inference.py",
		"--input", filepath,
		"--output", outputFilepath,
		"--model", "ml_server/models/UANet_VGG.ckpt",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Error("failed to execute python script",
			zap.Error(err),
			zap.String("output", string(output)))
		return "", "", "", fmt.Errorf("building extraction failed: %v", err)
	}

	// 5. Go 端最终想返回 _mask.png 对应的访问 URL
	// Python 里最后实际生成的文件是 xxx_mask.png
	maskFilename := strings.TrimSuffix(outputFilename, path.Ext(outputFilename)) + "_mask.png"
	maskURL := "/results/" + maskFilename

	//6. 成功调用，用户剩余提取次数-1
	err = s.dao.DeductRemainingCount(s.db, userID)
	if err != nil {
		s.logger.Error("deduct remainning count failed", zap.Error(err))
		return "", "", "", fmt.Errorf("building extraction failed: %v", err)
	}

	resultFilepath := "results/" + maskFilename
	return maskURL, filepath, resultFilepath, nil
}

// 辅助函数：验证文件是否为有效的图片类型
func isValidImageType(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".tif":  true,
		".tiff": true,
		".bmp":  true,
	}
	return validExts[ext]
}

// 辅助函数：生成唯一文件名
func generateUniqueFilename(originalFilename string) string {
	ext := path.Ext(originalFilename)
	uuid := uuid.New().String()
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d_%s%s", timestamp, uuid, ext)
}

func (s *BuildingExtractionService) GetAllProjects(userID int) ([]model.Project, error) {
	projects, err := s.dao.GetAllProjects(s.db, userID)
	if err != nil {
		s.logger.Error("get all projects failed", zap.Error(err))
		return nil, err
	}
	return projects, nil
}

func (s *BuildingExtractionService) SaveProject(saveProjectRequest dto.SaveProjectRequest, userID int) error {
	project := model.Project{
		ProjectName: saveProjectRequest.ProjectName,
		InputImage:  saveProjectRequest.InputImage,
		OutputImage: saveProjectRequest.OutputImage,
		ModelName:   saveProjectRequest.ModelName,
		UserID:      userID,
	}
	err := s.dao.CreateProject(s.db, &project)
	if err != nil {
		s.logger.Error("save project failed", zap.Error(err))
		return err
	}
	return nil
}

func (s *BuildingExtractionService) GetProjectById(projectId string) (*model.Project, error) {
	project, err := s.dao.GetProjectById(s.db, projectId)
	if err != nil {
		s.logger.Error("get project by id failed", zap.Error(err))
		return nil, err
	}
	return project, nil
}

func (s *BuildingExtractionService) DeleteProject(projectId string) error {
	project, err := s.dao.DeleteProject(s.db, projectId)
	if err != nil {
		s.logger.Error("delete project failed", zap.Error(err))
		return err
	}

	//删除本地文件夹中存储的图片
	os.Remove(project.InputImage)
	os.Remove(project.OutputImage)
	return nil
}

func (s *BuildingExtractionService) UpdatePassword(user *model.User, currentPassword, newPassword string) error {

	//验证旧密码
	password, err := s.dao.GetUserPassword(s.db, user.Username)
	if err != nil {
		s.logger.Error("fail to get password", zap.Error(err))
		return err
	}

	if password != currentPassword {
		s.logger.Error("wrong password")
		return fmt.Errorf("wrong password")
	}

	// // 验证新密码
	// if !passwordRegex.MatchString(newPassword) ||
	// 	!strings.ContainsAny(newPassword, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"+
	// 		"abcdefghijklmnopqrstuvwxyz") ||
	// 	!strings.ContainsAny(newPassword, "0123456789") {
	// 	return fmt.Errorf("password must be at least 8 characters, contain at least 1 letter and 1 number")
	// }

	user.Password = newPassword
	return s.dao.UpdateUser(s.db, user)
}
