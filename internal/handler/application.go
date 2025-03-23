package handler

import (
	"net/http"
	"strconv"
	"time"

	"internship-manager/internal/model"
	"internship-manager/internal/service"

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
		Company   string `json:"company" binding:"required"`
		Position  string `json:"position" binding:"required"`
		EventLink string `json:"event_link" binding:"required"`

		// 可选字段
		Location    string `json:"location"`
		Salary      string `json:"salary"`
		ContactInfo string `json:"contact_info"`
		Notes       string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 处理可选字段的默认值
	if req.Location == "" {
		req.Location = "无"
	}
	if req.Salary == "" {
		req.Salary = "无"
	}
	if req.ContactInfo == "" {
		req.ContactInfo = "无"
	}
	if req.Notes == "" {
		req.Notes = "无"
	}

	// 创建申请记录
	application := model.Application{
		UserID:      userID,
		Company:     req.Company,
		Position:    req.Position,
		Status:      model.StatusSubmitted,
		EventLink:   req.EventLink,
		Location:    req.Location,
		Salary:      req.Salary,
		ContactInfo: req.ContactInfo,
		Notes:       req.Notes,
		ApplyDate:   time.Now(),
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

// UpdateEvent 更新面试/笔试事件
func (h *ApplicationHandler) UpdateEvent(c *gin.Context) {
	var req struct {
		ID        uint      `json:"id" binding:"required"`
		EventTime time.Time `json:"event_time" binding:"required"`
		EventType string    `json:"event_type" binding:"required"`
		EventLink string    `json:"event_link"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	err := h.applicationService.UpdateNextEvent(req.ID, req.EventTime, req.EventType, req.EventLink)
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
		ID          uint   `json:"id" binding:"required"`
		Company     string `json:"company" binding:"required"`
		Position    string `json:"position" binding:"required"`
		EventLink   string `json:"event_link" binding:"required"`
		Location    string `json:"location"`
		Salary      string `json:"salary"`
		ContactInfo string `json:"contact_info"`
		Notes       string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 处理可选字段的默认值
	if req.Location == "" {
		req.Location = "无"
	}
	if req.Salary == "" {
		req.Salary = "无"
	}
	if req.ContactInfo == "" {
		req.ContactInfo = "无"
	}
	if req.Notes == "" {
		req.Notes = "无"
	}

	// 更新申请记录
	updates := map[string]interface{}{
		"company":      req.Company,
		"position":     req.Position,
		"event_link":   req.EventLink,
		"location":     req.Location,
		"salary":       req.Salary,
		"contact_info": req.ContactInfo,
		"notes":        req.Notes,
	}

	err := h.applicationService.UpdateApplication(req.ID, userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetApplications 获取用户的所有申请
// GetApplications 获取用户的申请（分页）
func (h *ApplicationHandler) GetApplications(c *gin.Context) {
	userID := c.GetUint("userID")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	applications, total, err := h.applicationService.GetApplicationsWithPagination(userID, page, pageSize)
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

// GetUpcomingEvents 获取即将到来的面试/笔试事件
func (h *ApplicationHandler) GetUpcomingEvents(c *gin.Context) {
	userID := c.GetUint("userID")
	events, err := h.applicationService.GetUpcomingEvents(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
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
