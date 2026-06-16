package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-notification/internal/domain/entity"
	"leap-one/service-notification/internal/domain/repository"
	"leap-one/service-notification/internal/interfaces/api/dto"
)

// NotificationHandler 通知消息管理Handler
type NotificationHandler struct {
	notiRepo repository.NotificationRepository
	subRepo  repository.SubscriptionRepository
	logger   *zap.Logger
}

func NewNotificationHandler(notiRepo repository.NotificationRepository, subRepo repository.SubscriptionRepository, logger *zap.Logger) *NotificationHandler {
	return &NotificationHandler{notiRepo: notiRepo, subRepo: subRepo, logger: logger}
}

// ListNotifications 通知列表 (GET /api/v1/notifications)
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	receiverIDStr := c.Query("receiver_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	unreadOnly := c.DefaultQuery("unread_only", "false") == "true"
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	receiverID, err := uuid.Parse(receiverIDStr)
	if err != nil || receiverIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的receiver_id"})
		return
	}

	ctx := c.Request.Context()
	list, total, err := h.notiRepo.ListByReceiver(ctx, receiverID, page, size, unreadOnly)
	if err != nil {
		h.logger.Error("查询通知列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	unreadCount, _ := h.notiRepo.CountUnread(ctx, receiverID)

	items := make([]dto.NotificationInfo, len(list))
	for i, n := range list {
		items[i] = buildNotificationInfo(n)
	}

	c.JSON(http.StatusOK, dto.NotificationListResponse{List: items, Total: total, Page: page, Size: size, UnreadCount: unreadCount})
}

// GetUnreadCount 未读数量 (GET /api/v1/notifications/unread-count)
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	receiverIDStr := c.Query("receiver_id")
	receiverID, err := uuid.Parse(receiverIDStr)
	if err != nil || receiverIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的receiver_id"})
		return
	}

	ctx := c.Request.Context()
	count, err := h.notiRepo.CountUnread(ctx, receiverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, dto.UnreadCountResponse{UnreadCount: count})
}

// MarkAsRead 标记已读 (PUT /api/v1/notifications/:id/read)
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的通知ID格式"})
		return
	}
	ctx := c.Request.Context()
	if err := h.notiRepo.MarkAsRead(ctx, id); err != nil {
		h.logger.Error("标记已读失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已标记为已读"})
}

// MarkAllAsRead 全部标记已读 (PUT /api/v1/notifications/read-all)
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	var req struct {
		ReceiverID uuid.UUID `json:"receiver_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	ctx := c.Request.Context()
	if err := h.notiRepo.MarkAllAsRead(ctx, req.ReceiverID); err != nil {
		h.logger.Error("全部标记已读失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已全部标记为已读"})
}

// DeleteNotification 删除通知 (DELETE /api/v1/notifications/:id)
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的通知ID格式"})
		return
	}
	ctx := c.Request.Context()
	if err := h.notiRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除通知失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "通知删除成功"})
}

// GetSettings 获取订阅设置 (GET /api/v1/notifications/settings)
func (h *NotificationHandler) GetSettings(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的user_id"})
		return
	}

	ctx := c.Request.Context()
	subs, err := h.subRepo.ListByUser(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	items := make([]dto.SubscriptionInfo, len(subs))
	for i, s := range subs {
		items[i] = dto.SubscriptionInfo{
			ID: s.ID.String(), UserID: s.UserID.String(), EventType: s.EventType,
			Channel: s.Channel, Enabled: s.Enabled,
		}
	}
	c.JSON(http.StatusOK, gin.H{"subscriptions": items})
}

// UpdateSettings 更新订阅设置 (PUT /api/v1/notifications/settings)
func (h *NotificationHandler) UpdateSettings(c *gin.Context) {
	var req dto.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	ctx := c.Request.Context()
	subs := make([]*entity.NotificationSubscription, len(req.Subscriptions))
	for i, s := range req.Subscriptions {
		userID, _ := uuid.Parse(s.UserID)
		subs[i] = &entity.NotificationSubscription{
			UserID: userID, EventType: s.EventType, Channel: s.Channel, Enabled: s.Enabled,
		}
	}
	if err := h.subRepo.BatchUpsert(ctx, subs); err != nil {
		h.logger.Error("更新订阅设置失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "订阅设置更新成功"})
}

func buildNotificationInfo(n *entity.Notification) dto.NotificationInfo {
	info := dto.NotificationInfo{
		ID: n.ID.String(), ReceiverID: n.ReceiverID.String(), Title: n.Title,
		Content: n.Content, Type: n.Type, Channel: n.Channel, ActionURL: n.ActionURL,
		IsRead: n.IsRead, SentAt: n.SentAt.Format("2006-01-02 15:04:05"), CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if n.SenderID != nil {
		info.SenderID = n.SenderID.String()
	}
	if n.ReadAt != nil {
		info.ReadAt = n.ReadAt.Format("2006-01-02 15:04:05")
	}
	return info
}
