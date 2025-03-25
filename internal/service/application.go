package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"internship-manager/internal/model"
	"internship-manager/pkg/database"
)

type ApplicationService struct{}

// GetRecentApplicationsByUserID 获取用户最近的n条申请记录
func (s *ApplicationService) GetRecentApplicationsByUserID(userID uint, limit int) ([]model.Application, error) {
	var applications []model.Application

	// 正向枚举查询条件
	validStatuses := []model.ApplicationStatus{
		model.StatusSubmitted,
		model.StatusWritten,
		model.StatusInterview,
		model.StatusAccepted,
	}

	// 显式指定查询字段（包含排序需要的updated_at）
	result := database.DB.Select(
		"company",
		"position",
		"status",
		"event_link",
		"updated_at", // 必须包含排序字段
	).Where(
		"user_id = ? AND status IN (?)",
		userID,
		validStatuses,
	).Where(
		"deleted_at IS NULL",
	).Order(
		"updated_at DESC",
	).Limit(
		limit,
	).Find(&applications)

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

// UpdateApplicationStatus 更新状态
func (s *ApplicationService) UpdateApplicationStatus(id uint, status model.ApplicationStatus) error {
	result := database.DB.Model(&model.Application{}).Where("id = ?", id).Update("status", status)
	if result.RowsAffected == 0 {
		return errors.New("申请记录不存在")
	}
	return result.Error
}

// GetApplicationStatistics 获取申请统计信息
func (s *ApplicationService) GetApplicationStatistics(userID uint) (map[string]int, error) {
	var stats = make(map[string]int)

	var result []struct {
		Status string `gorm:"column:status"`
		Count  int    `gorm:"column:count"`
	}

	err := database.DB.Raw(`
        SELECT 
            status, 
            COUNT(*) AS count 
        FROM 
            applications 
        WHERE 
            user_id = ? 
            AND status IN (?, ?, ?, ?,?) 
            AND deleted_at IS NULL 
        GROUP BY 
            status
    `, userID, model.StatusAccepted, model.StatusRejected, model.StatusInterview, model.StatusWritten, model.StatusSubmitted).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	for _, row := range result {
		stats[row.Status] = row.Count
	}

	// 确保所有状态都有默认值
	stats[string(model.StatusAccepted)] = stats[string(model.StatusAccepted)]
	stats[string(model.StatusRejected)] = stats[string(model.StatusRejected)]
	stats[string(model.StatusInterview)] = stats[string(model.StatusInterview)]
	stats[string(model.StatusWritten)] = stats[string(model.StatusWritten)]
	stats[string(model.StatusSubmitted)] = stats[string(model.StatusSubmitted)]
	fmt.Println(stats)
	return stats, nil
}

//// GetApplicationsWithPagination 获取分页的申请记录
//func (s *ApplicationService) GetApplicationsWithPagination(userID uint, page, pageSize int) ([]model.Application, int64, error) {
//	var applications []model.Application
//	var total int64
//
//	// 先获取总记录数
//	if err := database.DB.Model(&model.Application{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
//		return nil, 0, err
//	}
//
//	// 获取分页数据
//	offset := (page - 1) * pageSize
//	result := database.DB.Where("user_id = ?", userID).
//		Order("updated_at DESC").
//		Offset(offset).
//		Limit(pageSize).
//		Find(&applications)
//
//	if result.Error != nil {
//		return nil, 0, result.Error
//	}
//
//	return applications, total, nil
//}

// GetApplicationsWithPagination 获取分页的申请记录
func (s *ApplicationService) GetApplicationsWithPagination(userID uint, page, pageSize int, searchQuery string, statuses []string) ([]model.Application, int64, error) {
	var applications []model.Application
	var total int64

	// 构建基础查询
	baseQuery := database.DB.Table("applications USE INDEX (idx_user_deleted_status_updated)").
		Select(
			"id",
			"company",
			"position",
			"status",
			"event_link",
			"updated_at",
		).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	// 如果有搜索关键词，添加公司名称搜索条件
	if searchQuery != "" {
		baseQuery = baseQuery.Where("company LIKE ?", "%"+searchQuery+"%")
	}

	// 如果有状态筛选，添加状态条件
	if len(statuses) > 0 {
		baseQuery = baseQuery.Where("status IN ?", statuses)
	}

	// 获取总记录数（使用克隆的查询以避免影响主查询）
	countQuery := baseQuery.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据（不包含 notes 字段）
	offset := (page - 1) * pageSize
	var basicResults []model.Application
	if err := baseQuery.Order("updated_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&basicResults).Error; err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回
	if len(basicResults) == 0 {
		return basicResults, total, nil
	}

	// 获取记录的 ID 列表
	var ids []uint
	for _, app := range basicResults {
		ids = append(ids, app.ID)
	}

	// 查询这些记录的 notes 字段
	var notesResults []model.Application
	if err := database.DB.Model(&model.Application{}).
		Select("id", "notes").
		Where("id IN ?", ids).
		Find(&notesResults).Error; err != nil {
		return nil, 0, err
	}

	// 将 notes 字段合并到结果中
	notesMap := make(map[uint]string)
	for _, note := range notesResults {
		notesMap[note.ID] = note.Notes
	}

	// 组装最终结果
	applications = basicResults
	for i := range applications {
		applications[i].Notes = notesMap[applications[i].ID]
	}

	return applications, total, nil
}
