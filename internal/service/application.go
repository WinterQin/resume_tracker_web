package service

import (
	"errors"
	"gorm.io/gorm"
	"internship-manager/internal/model"
	"internship-manager/pkg/database"
	"time"
)

type ApplicationService struct{}

// CreateApplication 创建实习申请记录
func (s *ApplicationService) CreateApplication(userID uint, company, position string) error {
	application := model.Application{
		UserID:    userID,
		Company:   company,
		Position:  position,
		Status:    model.StatusSubmitted,
		ApplyDate: time.Now(),
	}
	return database.DB.Create(&application).Error
}

// GetRecentApplicationsByUserID 获取用户最近的n条申请记录
func (s *ApplicationService) GetRecentApplicationsByUserID(userID uint, limit int) ([]model.Application, error) {
	var applications []model.Application

	result := database.DB.Where("user_id = ? AND status != ?", userID, model.StatusRejected).
		Order("updated_at DESC").
		Limit(limit).
		Find(&applications)

	if result.Error != nil {
		return nil, result.Error
	}

	return applications, nil
}

// CreateApplicationFull 创建完整的实习申请记录
func (s *ApplicationService) CreateApplicationFull(application *model.Application) error {
	return database.DB.Create(application).Error
}

// UpdateApplication 更新申请记录
func (s *ApplicationService) UpdateApplication(id uint, userID uint, updates map[string]interface{}) error {
	result := database.DB.Model(&model.Application{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(updates)

	if result.RowsAffected == 0 {
		return errors.New("申请记录不存在或无权限更新")
	}

	return result.Error
}

// DeleteApplication 删除实习申请记录
func (s *ApplicationService) DeleteApplication(id uint, userID uint) error {
	// 先检查记录是否存在且属于该用户
	var application model.Application
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&application)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("申请记录不存在或无权限删除")
		}
		return result.Error
	}

	// 执行删除操作
	return database.DB.Delete(&application).Error
}

// UpdateApplicationStatus 更新申请状态
func (s *ApplicationService) UpdateApplicationStatus(id uint, status model.ApplicationStatus) error {
	result := database.DB.Model(&model.Application{}).Where("id = ?", id).Update("status", status)
	if result.RowsAffected == 0 {
		return errors.New("申请记录不存在")
	}
	return result.Error
}

// UpdateNextEvent 更新下一个面试/笔试事件
func (s *ApplicationService) UpdateNextEvent(id uint, eventTime time.Time, eventType, eventLink string) error {
	updates := map[string]interface{}{
		"next_event": eventTime,
		"event_type": eventType,
		"event_link": eventLink,
	}
	result := database.DB.Model(&model.Application{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return errors.New("申请记录不存在")
	}
	return result.Error
}

// GetApplicationsByUserID 获取用户的所有申请记录
func (s *ApplicationService) GetApplicationsByUserID(userID uint) ([]model.Application, error) {
	var applications []model.Application
	err := database.DB.Where("user_id = ?", userID).Find(&applications).Error
	return applications, err
}

// GetApplicationStatistics 获取申请统计信息
func (s *ApplicationService) GetApplicationStatistics(userID uint) (map[string]int, error) {
	var stats = make(map[string]int)
	var applications []model.Application

	err := database.DB.Where("user_id = ?", userID).Find(&applications).Error
	if err != nil {
		return nil, err
	}

	// 统计各状态数量
	for _, app := range applications {
		stats[string(app.Status)]++
	}

	return stats, nil
}

// GetUpcomingEvents 获取即将到来的面试/笔试事件
func (s *ApplicationService) GetUpcomingEvents(userID uint) ([]model.Application, error) {
	var applications []model.Application
	now := time.Now()
	err := database.DB.Where("user_id = ? AND next_event > ?", userID, now).
		Order("next_event asc").
		Find(&applications).Error
	return applications, err
}

// GetApplicationsWithPagination 获取分页的申请记录
func (s *ApplicationService) GetApplicationsWithPagination(userID uint, page, pageSize int) ([]model.Application, int64, error) {
	var applications []model.Application
	var total int64

	// 先获取总记录数
	if err := database.DB.Model(&model.Application{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	result := database.DB.Where("user_id = ?", userID).
		Order("updated_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&applications)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return applications, total, nil
}
