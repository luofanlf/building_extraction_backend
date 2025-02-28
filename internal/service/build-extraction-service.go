package service

import (
	"building-extraction/internal/dao"
	"building-extraction/internal/model"
	"fmt"
	"io"
	"mime/multipart"
	"os"
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
		Username: username,
		Password: password,
	}

	// 验证密码：至少8个字符，包含至少1个字母和1个数字
	if !passwordRegex.MatchString(password) ||
		!strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"+
			"abcdefghijklmnopqrstuvwxyz") ||
		!strings.ContainsAny(password, "0123456789") {
		return fmt.Errorf("password must be at least 8 characters, contain at least 1 letter and 1 number")
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

func (s *BuildingExtractionService) UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	if !isValidImageType(header.Filename) {
		s.logger.Error("invalid file type", zap.Error(fmt.Errorf("invalid file type. only image files are allowed")))
		return "", fmt.Errorf("invalid file type. only image files are allowed")
	}

	// 3. 验证文件大小 (限制为20MB)
	if header.Size > 20*1024*1024 {
		s.logger.Error("file size exceeds the limit (20MB)")
		return "", fmt.Errorf("file size exceeds the limit (20MB)")
	}

	// 4. 为文件生成唯一文件名
	filename := generateUniqueFilename(header.Filename)

	// 5. 确保上传目录存在
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		s.logger.Error("failed to create upload directory", zap.Error(err))
		return "", fmt.Errorf("failed to create upload directory")
	}

	// 6. 创建目标文件
	filepath := path.Join(uploadDir, filename)
	out, err := os.Create(filepath)
	if err != nil {
		s.logger.Error("failed to create file on server", zap.Error(err))
		return "", fmt.Errorf("failed to create file on server")
	}
	defer out.Close()

	// 7. 将上传的文件内容复制到目标文件
	_, err = io.Copy(out, file)
	if err != nil {
		s.logger.Error("failed to save file on server", zap.Error(err))
		return "", fmt.Errorf("failed to save file on server")
	}

	// 8. 获取文件的相对URL路径
	fileURL := "/uploads/" + filename
	return fileURL, nil
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
