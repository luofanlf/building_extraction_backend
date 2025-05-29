package dao

import (
	"building-extraction/internal/model"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type BuildingExtractionDao struct {
	logger *zap.Logger
}

func NewBuildingExtractionDao() *BuildingExtractionDao {
	logger, _ := zap.NewProduction()

	// 调整堆栈级别为 Fatal
	logger = logger.WithOptions(zap.AddStacktrace(zapcore.FatalLevel))

	return &BuildingExtractionDao{
		logger: logger,
	}
}

func (d *BuildingExtractionDao) CheckUserExists(db *gorm.DB, username string) (bool, error) {
	var count int64
	err := db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		d.logger.Error("check user exists failed", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

func (d *BuildingExtractionDao) CreateUser(db *gorm.DB, user *model.User) error {
	err := db.Create(user).Error
	if err != nil {
		d.logger.Error("create user failed", zap.Error(err))
		return err
	}
	return nil
}

func (d *BuildingExtractionDao) GetUserByUsername(db *gorm.DB, username string) (*model.User, error) {
	var user model.User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		d.logger.Error("get user by username failed", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (d *BuildingExtractionDao) GetAllProjects(db *gorm.DB, userID int) ([]model.Project, error) {
	var projects []model.Project
	err := db.Model(&model.Project{}).Where("user_id = ?", userID).Find(&projects).Error
	if err != nil {
		d.logger.Error("get all projects failed", zap.Error(err))
		return nil, err
	}
	return projects, nil
}

func (d *BuildingExtractionDao) CreateProject(db *gorm.DB, project *model.Project) error {
	err := db.Create(project).Error
	if err != nil {
		d.logger.Error("create project failed", zap.Error(err))
		return err
	}
	return nil
}

func (d *BuildingExtractionDao) GetProjectById(db *gorm.DB, projectId string) (*model.Project, error) {
	var project model.Project
	err := db.Model(&model.Project{}).Where("id = ?", projectId).First(&project).Error
	if err != nil {
		d.logger.Error("get project by id failed", zap.Error(err))
		return nil, err
	}
	return &project, nil
}

func (d *BuildingExtractionDao) DeleteProject(db *gorm.DB, projectId string) (*model.Project, error) {
	var project *model.Project
	err := db.Model(&model.Project{}).Where("id = ?", projectId).First(&project).Error
	if err != nil {
		d.logger.Error("delete project failed", zap.Error(err))
		return nil, err
	}

	err = db.Model(&model.Project{}).Where("id = ?", projectId).Delete(&project).Error
	if err != nil {
		d.logger.Error("delete project failed", zap.Error(err))
		return nil, err
	}
	return project, nil
}

func (d *BuildingExtractionDao) UpdateUser(db *gorm.DB, user *model.User) error {
	// 只更新密码字段
	err := db.Model(&model.User{}).Where("id = ?", user.ID).Update("password", user.Password).Error
	if err != nil {
		d.logger.Error("update user failed", zap.Error(err))
		return err
	}
	return nil
}

func (d *BuildingExtractionDao) GetUserPassword(db *gorm.DB, username string) (string, error) {
	var user model.User
	err := db.Model(&model.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		d.logger.Error("fail to get password", zap.Error(err))
		return "", err
	}
	return user.Password, nil
}

func (d *BuildingExtractionDao) DeductRemainingCount(db *gorm.DB, userID int) error {
	err := db.Model(&model.User{}).Where("id = ?", userID).UpdateColumn("remaining_count", gorm.Expr("remaining_count - ?", 1)).Error
	if err != nil {
		d.logger.Error("deduct remaining count failed", zap.Error(err))
		return err
	}
	return nil
}

func (d *BuildingExtractionDao) CheckRemainingCount(db *gorm.DB, userID int) (int, error) {
	var user model.User
	err := db.Model(&model.User{}).Where("id = ?", userID).First(&user).Error
	if err != nil {
		d.logger.Error("get user info failed", zap.Error(err))
		return 0, err
	}
	return user.RemainingCount, nil
}
