package dao

import (
	"building-extraction/internal/model"
	"time"

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

func (d *BuildingExtractionDao) GetUserByUserID(db *gorm.DB, userID int) (*model.User, error) {
	var user model.User
	err := db.Where("id = ?", userID).First(&user).Error
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
	err := db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"remaining_count":  gorm.Expr("remaining_count - ?", 1),
		"extraction_count": gorm.Expr("extraction_count + ?", 1),
	}).Error
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

func (d *BuildingExtractionDao) GetTotalUsers(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&model.User{}).Count(&count).Error
	if err != nil {
		d.logger.Error("get total users failed", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (d *BuildingExtractionDao) GetTotalExtractions(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&model.User{}).Select("SUM(extraction_count)").Scan(&count).Error
	if err != nil {
		d.logger.Error("get total extractions failed", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (d *BuildingExtractionDao) GetActiveUsersToday(db *gorm.DB) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := db.Model(&model.User{}).
		Where("DATE(updated_at) = ?", today).
		Count(&count).Error
	if err != nil {
		d.logger.Error("get active users today failed", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (d *BuildingExtractionDao) GetExtractionsToday(db *gorm.DB) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := db.Model(&model.Project{}).
		Where("DATE(created_at) = ?", today).
		Count(&count).Error
	if err != nil {
		d.logger.Error("get extractions today failed", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (d *BuildingExtractionDao) CreateExtractionRequest(db *gorm.DB, request *model.ExtractionRequest) error {
	err := db.Create(request).Error
	if err != nil {
		d.logger.Error("create extraction request failed", zap.Error(err))
		return err
	}
	return nil
}

func (d *BuildingExtractionDao) GetPendingRequests(db *gorm.DB) ([]model.ExtractionRequest, error) {
	var requests []model.ExtractionRequest
	err := db.Model(&model.ExtractionRequest{}).
		Where("status = ?", model.RequestStatusPending).
		Preload("User"). // 现在可以正确预加载User关联
		Find(&requests).Error
	if err != nil {
		d.logger.Error("get pending requests failed", zap.Error(err))
		return nil, err
	}
	return requests, nil
}

func (d *BuildingExtractionDao) HandleRequest(db *gorm.DB, requestID int, status model.RequestStatus, reply string) error {
	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 更新请求状态
		err := tx.Model(&model.ExtractionRequest{}).
			Where("id = ?", requestID).
			Updates(map[string]interface{}{
				"status": status,
				"reply":  reply,
			}).Error
		if err != nil {
			return err
		}

		// 2. 如果批准请求，增加用户的剩余次数
		if status == model.RequestStatusApproved {
			var request model.ExtractionRequest
			if err := tx.First(&request, requestID).Error; err != nil {
				return err
			}

			err = tx.Model(&model.User{}).
				Where("id = ?", request.UserID).
				UpdateColumn("remaining_count", gorm.Expr("remaining_count + ?", request.RequestCount)).
				Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (d *BuildingExtractionDao) GetAllUsers(db *gorm.DB) ([]model.User, error) {
	var users []model.User
	err := db.Model(&model.User{}).Find(&users).Error
	if err != nil {
		d.logger.Error("get all users failed", zap.Error(err))
		return nil, err
	}
	return users, nil
}

func (d *BuildingExtractionDao) GetUserRequests(db *gorm.DB, userID int) ([]model.ExtractionRequest, error) {
	var requests []model.ExtractionRequest
	err := db.Model(&model.ExtractionRequest{}).
		Where("user_id = ?", userID).
		Order("created_at DESC"). // 按创建时间倒序排列
		Find(&requests).Error
	if err != nil {
		d.logger.Error("get user requests failed", zap.Error(err))
		return nil, err
	}
	return requests, nil
}
