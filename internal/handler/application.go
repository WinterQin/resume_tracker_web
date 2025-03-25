package handler

import (
	"internship-manager/internal/model"
	"internship-manager/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	applicationService *service.ApplicationService
}

func NewApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{
		applicationService: &service.ApplicationService{},
	}
}

// CreateApplication 创建实习申请
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		// 必填字段
		Company  string `json:"company" binding:"required"`
		Position string `json:"position" binding:"required"`

		// 可选字段
		EventLink string `json:"event_link"`
		Notes     string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Notes == "" {
		req.Notes = "无"
	}

	// 创建申请记录
	application := model.Application{
		UserID:    userID,
		Company:   req.Company,
		Position:  req.Position,
		Status:    model.StatusSubmitted,
		EventLink: req.EventLink,
		Notes:     req.Notes,
	}

	err := h.applicationService.CreateApplicationFull(&application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// UpdateStatus 更新申请状态
func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		ID     uint                    `json:"id" binding:"required"`
		Status model.ApplicationStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	err := h.applicationService.UpdateApplicationStatus(req.ID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// UpdateApplication 更新申请信息
func (h *ApplicationHandler) UpdateApplication(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		ID        uint   `json:"id" binding:"required"`
		Company   string `json:"company" binding:"required"`
		Position  string `json:"position" binding:"required"`
		EventLink string `json:"event_link" binding:"required"`
		Notes     string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if req.Notes == "" {
		req.Notes = "无"
	}

	// 更新申请记录
	updates := map[string]interface{}{
		"company":    req.Company,
		"position":   req.Position,
		"event_link": req.EventLink,
		"notes":      req.Notes,
	}

	err := h.applicationService.UpdateApplication(req.ID, userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

//// GetApplications 获取用户的申请（分页）
//func (h *ApplicationHandler) GetApplications(c *gin.Context) {
//	userID := c.GetUint("userID")
//
//	// 获取分页参数
//	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
//	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
//
//	// 参数验证
//	if page < 1 {
//		page = 1
//	}
//	if pageSize < 1 {
//		pageSize = 10
//	}
//
//	applications, total, err := h.applicationService.GetApplicationsWithPagination(userID, page, pageSize)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"applications": applications,
//		"total":        total,
//		"current_page": page,
//		"page_size":    pageSize,
//	})
//}

// GetApplications 获取用户的申请（分页）
func (h *ApplicationHandler) GetApplications(c *gin.Context) {
	userID := c.GetUint("userID")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	searchQuery := c.DefaultQuery("search", "")

	// 获取状态筛选参数
	statusesStr := c.DefaultQuery("statuses", "")
	var statuses []string
	if statusesStr != "" {
		statuses = strings.Split(statusesStr, ",")
	}

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	applications, total, err := h.applicationService.GetApplicationsWithPagination(userID, page, pageSize, searchQuery, statuses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"applications": applications,
		"total":        total,
		"current_page": page,
		"page_size":    pageSize,
	})
}

// GetRecentApplications 获取用户最近的5条申请记录
func (h *ApplicationHandler) GetRecentApplications(c *gin.Context) {
	userID := c.GetUint("userID")

	// 获取最近5条记录
	applications, err := h.applicationService.GetRecentApplicationsByUserID(userID, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"applications": applications})
}

// GetStatistics 获取申请统计信息
func (h *ApplicationHandler) GetStatistics(c *gin.Context) {
	userID := c.GetUint("userID")
	stats, err := h.applicationService.GetApplicationStatistics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statistics": stats})
}

// DeleteApplication 删除实习申请
func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	// 将 id 转换为 uint
	applicationID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = h.applicationService.DeleteApplication(uint(applicationID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
