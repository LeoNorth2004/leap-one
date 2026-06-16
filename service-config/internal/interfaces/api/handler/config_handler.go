package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"leap-one/service-config/internal/domain/entity"
	"leap-one/service-config/internal/domain/repository"
	"leap-one/service-config/internal/interfaces/api/dto"
	"go.uber.org/zap"
)

// ConfigHandler 系统配置Handler
type ConfigHandler struct {
	cfgRepo repository.SystemConfigRepository
	logger  *zap.Logger
}

func NewConfigHandler(cfgRepo repository.SystemConfigRepository, logger *zap.Logger) *ConfigHandler {
	return &ConfigHandler{cfgRepo: cfgRepo, logger: logger}
}

func (h *ConfigHandler) ListConfigs(c *gin.Context) {
	category := c.Query("category")
	ctx := c.Request.Context()
	var items []*entity.SystemConfig
	var err error
	if category != "" {
		items, err = h.cfgRepo.ListByCategory(ctx, category)
	} else {
		items, err = h.cfgRepo.ListAll(ctx)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	result := make([]dto.ConfigItem, len(items))
	for i, cfg := range items {
		result[i] = buildConfigItem(cfg)
	}
	c.JSON(http.StatusOK, gin.H{"list": result})
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {
	ctx := c.Request.Context()
	cfg, err := h.cfgRepo.GetByCategoryAndKey(ctx, c.Param("category"), c.Param("key"))
	if err != nil || cfg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "配置不存�?})
		return
	}
	c.JSON(http.StatusOK, buildConfigItem(cfg))
}

func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var req dto.UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	ctx := c.Request.Context()
	cfg, err := h.cfgRepo.GetByCategoryAndKey(ctx, c.Param("category"), c.Param("key"))
	if err != nil || cfg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "不存�?})
		return
	}
	cfg.Value = req.Value
	h.cfgRepo.Update(ctx, cfg)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *ConfigHandler) BatchUpdateConfigs(c *gin.Context) {
	var req dto.BatchUpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	ctx := c.Request.Context()
	configs := make([]*entity.SystemConfig, len(req.Configs))
	for i, u := range req.Configs {
		cfg, _ := h.cfgRepo.GetByCategoryAndKey(ctx, u.Category, u.Key)
		if cfg == nil {
			cfg = &entity.SystemConfig{Category: u.Category, Key: u.Key, ValueType: "string"}
		}
		cfg.Value = u.Value
		configs[i] = cfg
	}
	h.cfgRepo.BatchUpdate(ctx, configs)
	c.JSON(http.StatusOK, gin.H{"message": "批量更新成功"})
}

func (h *ConfigHandler) GetConfigGroups(c *gin.Context) {
	ctx := c.Request.Context()
	groups, err := h.cfgRepo.GetGroups(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": "查询失败"})
		return
	}
	result := make(map[string][]dto.ConfigItem)
	for cat, items := range groups {
		lst := make([]dto.ConfigItem, len(items))
		for i, item := range items {
			lst[i] = buildConfigItem(item)
		}
		result[cat] = lst
	}
	c.JSON(http.StatusOK, dto.ConfigGroupResponse{Groups: result})
}

// FeatureFlagHandler 功能开关Handler
type FlagHandler struct {
	flagRepo repository.FeatureFlagRepository
	logger   *zap.Logger
}

func NewFlagHandler(flagRepo repository.FeatureFlagRepository, logger *zap.Logger) *FlagHandler {
	return &FlagHandler{flagRepo: flagRepo, logger: logger}
}

func (h *FlagHandler) ListFlags(c *gin.Context) {
	ctx := c.Request.Context()
	list, _ := h.flagRepo.List(ctx)
	items := make([]dto.FeatureFlagInfo, len(list))
	for i, f := range list {
		items[i] = dto.FeatureFlagInfo{ID: f.ID.String(), Key: f.Key, Name: f.Name, Description: f.Description, Enabled: f.Enabled, Rules: f.Rules, CreatedAt: f.CreatedAt.Format("2006-01-02 15:04:05")}
	}
	c.JSON(http.StatusOK, gin.H{"list": items})
}
func (h *FlagHandler) GetFlag(c *gin.Context) {
	ctx := c.Request.Context()
	f, err := h.flagRepo.GetByKey(ctx, c.Param("key"))
	if err != nil || f == nil {
		c.JSON(404, gin.H{"error": "不存�?})
		return
	}
	c.JSON(200, dto.FeatureFlagInfo{ID: f.ID.String(), Key: f.Key, Name: f.Name, Description: f.Description, Enabled: f.Enabled, Rules: f.Rules})
}
func (h *FlagHandler) UpdateFlag(c *gin.Context) {
	var req dto.UpdateFeatureFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	ctx := c.Request.Context()
	f, err := h.flagRepo.GetByKey(ctx, c.Param("key"))
	if err != nil || f == nil {
		c.JSON(404, gin.H{"error": "不存�?})
		return
	}
	if req.Name != nil {
		f.Name = *req.Name
	}
	if req.Description != nil {
		f.Description = *req.Description
	}
	if req.Enabled != nil {
		f.Enabled = *req.Enabled
	}
	if req.Rules != nil {
		f.Rules = *req.Rules
	}
	h.flagRepo.Update(ctx, f)
	c.JSON(200, gin.H{"message": "更新成功"})
}
func (h *FlagHandler) CreateFlag(c *gin.Context) {
	var req dto.CreateFeatureFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	f := &entity.FeatureFlag{Key: req.Key, Name: req.Name, Description: req.Description, Enabled: req.Enabled, Rules: req.Rules}
	ctx := c.Request.Context()
	h.flagRepo.Create(ctx, f)
	c.JSON(201, gin.H{"message": "创建成功", "flag_id": f.ID.String()})
}
func (h *FlagHandler) DeleteFlag(c *gin.Context) {
	ctx := c.Request.Context()
	h.flagRepo.Delete(ctx, c.Param("key"))
	c.JSON(200, gin.H{"message": "删除成功"})
}

// AuditLogHandler 审计日志Handler
type AuditLogHandler struct {
	logRepo repository.AuditLogRepository
	logger  *zap.Logger
}

func NewAuditLogHandler(logRepo repository.AuditLogRepository, logger *zap.Logger) *AuditLogHandler {
	return &AuditLogHandler{logRepo: logRepo, logger: logger}
}

func (h *AuditLogHandler) ListLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	action := c.Query("action")
	resource := c.Query("resource")
	userIDStr := c.Query("user_id")
	userID, _ := uuid.Parse(userIDStr)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	ctx := c.Request.Context()
	list, total, _ := h.logRepo.List(ctx, page, size, userID, action, resource)
	items := make([]dto.AuditLogInfo, len(list))
	for i, l := range list {
		items[i] = dto.AuditLogInfo{ID: l.ID.String(), UserID: l.UserID.String(), Action: l.Action, Resource: l.Resource, ResourceID: l.ResourceID.String(), Detail: l.Detail, IPAddress: l.IPAddress, UserAgent: l.UserAgent, CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05")}
	}
	c.JSON(200, dto.AuditLogListResponse{List: items, Total: total, Page: page, Size: size})
}
func (h *AuditLogHandler) GetLog(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	log, err := h.logRepo.GetByID(ctx, id)
	if err != nil || log == nil {
		c.JSON(404, gin.H{"error": "不存�?})
		return
	}
	c.JSON(200, dto.AuditLogInfo{ID: log.ID.String(), UserID: log.UserID.String(), Action: log.Action, Resource: log.Resource, ResourceID: log.ResourceID.String(), Detail: log.Detail, IPAddress: log.IPAddress, UserAgent: log.UserAgent, CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05")})
}

func buildConfigItem(cfg *entity.SystemConfig) dto.ConfigItem {
	return dto.ConfigItem{ID: cfg.ID.String(), Category: cfg.Category, Key: cfg.Key, Value: cfg.Value, ValueType: cfg.ValueType, IsEncrypted: cfg.IsEncrypted, IsPublic: cfg.IsPublic, Description: cfg.Description, SortOrder: cfg.SortOrder, CreatedAt: cfg.CreatedAt.Format("2006-01-02 15:04:05"), UpdatedAt: cfg.UpdatedAt.Format("2006-01-02 15:04:05")}
}
