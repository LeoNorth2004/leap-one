package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-user-org/internal/domain/entity"
	"leap-one/service-user-org/internal/domain/repository"
	"leap-one/service-user-org/internal/interfaces/api/dto"
)

// DepartmentHandler 部门管理Handler
type DepartmentHandler struct {
	deptRepo repository.DepartmentRepository
	userRepo repository.UserRepository
	logger   *zap.Logger
}

// NewDepartmentHandler 创建部门管理Handler实例
func NewDepartmentHandler(
	deptRepo repository.DepartmentRepository,
	userRepo repository.UserRepository,
	logger *zap.Logger,
) *DepartmentHandler {
	return &DepartmentHandler{
		deptRepo: deptRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

// CreateDepartment 创建部门（POST /api/v1/departments）
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 校验部门编码唯一性
	if existing, _ := h.deptRepo.GetByCode(ctx, req.Code); existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "部门编码已存在"})
		return
	}

	// 计算部门层级
	level := 1
	if req.ParentID != nil && *req.ParentID != uuid.Nil {
		parent, pErr := h.deptRepo.GetByID(ctx, *req.ParentID)
		if pErr != nil || parent == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "上级部门不存在"})
			return
		}
		level = parent.Level + 1
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	dept := &entity.Department{
		Name:      req.Name,
		Code:      req.Code,
		ParentID:  req.ParentID,
		Level:     level,
		SortOrder: sortOrder,
		Status:    1,
	}

	if req.Leader != nil {
		dept.Leader = *req.Leader
	}
	if req.Phone != nil {
		dept.Phone = *req.Phone
	}
	if req.Email != nil {
		dept.Email = *req.Email
	}
	if req.Description != nil {
		dept.Description = *req.Description
	}

	if err := h.deptRepo.Create(ctx, dept); err != nil {
		h.logger.Error("创建部门失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建部门失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "部门创建成功",
		"department_id": dept.ID.String(),
	})
}

// GetDepartment 获取部门详情（GET /api/v1/departments/:id）
func (h *DepartmentHandler) GetDepartment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的部门ID格式"})
		return
	}

	ctx := c.Request.Context()
	dept, err := h.deptRepo.GetByID(ctx, id)
	if err != nil || dept == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "部门不存在"})
		return
	}

	// 统计成员数量
	memberCount, _ := h.deptRepo.CountMembers(ctx, id)

	resp := dto.DepartmentInfo{
		ID:          dept.ID,
		Name:        dept.Name,
		Code:        dept.Code,
		ParentID:    dept.ParentID,
		Level:       dept.Level,
		Leader:      dept.Leader,
		Status:      dept.Status,
		MemberCount: int(memberCount),
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateDepartment 更新部门（PUT /api/v1/departments/:id）
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的部门ID格式"})
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	dept, err := h.deptRepo.GetByID(ctx, id)
	if err != nil || dept == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "部门不存在"})
		return
	}

	if req.Name != nil {
		dept.Name = *req.Name
	}
	if req.SortOrder != nil {
		dept.SortOrder = *req.SortOrder
	}
	if req.Leader != nil {
		dept.Leader = *req.Leader
	}
	if req.Phone != nil {
		dept.Phone = *req.Phone
	}
	if req.Email != nil {
		dept.Email = *req.Email
	}
	if req.Description != nil {
		dept.Description = *req.Description
	}
	if req.Status != nil {
		dept.Status = *req.Status
	}

	if err := h.deptRepo.Update(ctx, dept); err != nil {
		h.logger.Error("更新部门失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新部门失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "部门更新成功"})
}

// DeleteDepartment 删除部门（DELETE /api/v1/departments/:id）- 检查子部门和成员
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的部门ID格式"})
		return
	}

	ctx := c.Request.Context()

	// 检查部门是否存在
	dept, dErr := h.deptRepo.GetByID(ctx, id)
	if dErr != nil || dept == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "部门不存在"})
		return
	}

	// 检查是否有子部门
	hasChildren, hcErr := h.deptRepo.HasChildren(ctx, id)
	if hcErr != nil {
		h.logger.Error("检查子部门失败", zap.Error(hcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	if hasChildren {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该部门下还有子部门，请先删除或移动子部门"})
		return
	}

	// 检查是否有成员
	memberCount, mcErr := h.deptRepo.CountMembers(ctx, id)
	if mcErr == nil && memberCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该部门下还有成员，请先移除或转移成员"})
		return
	}

	if err := h.deptRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除部门失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除部门失败"})
		return
	}

	h.logger.Info("删除部门成功",
		zap.String("operator", getStringFromContext(c, "username")),
		zap.String("dept_name", dept.Name),
	)

	c.JSON(http.StatusOK, gin.H{"message": "部门删除成功"})
}

// ListDepartments 分页查询部门列表（GET /api/v1/departments）
func (h *DepartmentHandler) ListDepartments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	ctx := c.Request.Context()
	depts, total, err := h.deptRepo.List(ctx, page, size)
	if err != nil {
		h.logger.Error("查询部门列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询部门列表失败"})
		return
	}

	list := make([]*dto.DepartmentInfo, len(depts))
	for i, d := range depts {
		memberCount, _ := h.deptRepo.CountMembers(ctx, d.ID)
		list[i] = &dto.DepartmentInfo{
			ID:          d.ID,
			Name:        d.Name,
			Code:        d.Code,
			ParentID:    d.ParentID,
			Level:       d.Level,
			Leader:      d.Leader,
			Status:      d.Status,
			MemberCount: int(memberCount),
		}
	}

	c.JSON(http.StatusOK, dto.DepartmentListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// GetDepartmentTree 获取部门树形结构（GET /api/v1/departments/tree）
func (h *DepartmentHandler) GetDepartmentTree(c *gin.Context) {
	ctx := c.Request.Context()

	tree, err := h.deptRepo.GetTree(ctx)
	if err != nil {
		h.logger.Error("获取部门树失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取部门树失败"})
		return
	}

	resp := buildDeptTreeResponse(tree)
	c.JSON(http.StatusOK, gin.H{
		"tree":    resp,
		"message": "获取部门树成功",
	})
}

// GetDepartmentMembers 部门成员列表（GET /api/v1/departments/:id/members）
func (h *DepartmentHandler) GetDepartmentMembers(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的部门ID格式"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	ctx := c.Request.Context()
	users, total, err := h.deptRepo.GetDepartmentMembers(ctx, id, page, size)
	if err != nil {
		h.logger.Error("获取部门成员失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	list := make([]dto.UserInfo, len(users))
	for i, u := range users {
		list[i] = dto.UserInfo{
			ID:       u.ID.String(),
			Username: u.Username,
			RealName: u.RealName,
			Email:    u.Email,
			Phone:    u.Phone,
			Avatar:   u.Avatar,
			Status:   u.Status,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// MoveDepartment 移动部门（变更父级）（PUT /api/v1/departments/:id/move）
func (h *DepartmentHandler) MoveDepartment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的部门ID格式"})
		return
	}

	var req struct {
		ParentID *uuid.UUID `json:"parent_id"` // 新父级ID，Nil表示移到顶级
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 不能将自己设为自己的子级
	if req.ParentID != nil && *req.ParentID == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能将部门移动到自己下面"})
		return
	}

	ctx := c.Request.Context()
	var parentID uuid.UUID
	if req.ParentID != nil {
		parentID = *req.ParentID
	}
	if err := h.deptRepo.Move(ctx, id, parentID); err != nil {
		h.logger.Error("移动部门失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移动部门失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "部门移动成功"})
}

// buildDeptTreeResponse 将实体树转换为DTO响应树
func buildDeptTreeResponse(depts []*entity.Department) []*dto.DepartmentTreeResponse {
	if len(depts) == 0 {
		return []*dto.DepartmentTreeResponse{}
	}

	result := make([]*dto.DepartmentTreeResponse, len(depts))
	for i, d := range depts {
		result[i] = &dto.DepartmentTreeResponse{
			ID:        d.ID,
			Name:      d.Name,
			Code:      d.Code,
			ParentID:  d.ParentID,
			Level:     d.Level,
			SortOrder: d.SortOrder,
			Leader:    d.Leader,
			Status:    d.Status,
			Children:  buildDeptTreeResponsePtr(d.Children),
		}
	}
	return result
}

// buildDeptTreeResponsePtr 将值切片转换为指针切片后构建响应树
func buildDeptTreeResponsePtr(depts []entity.Department) []*dto.DepartmentTreeResponse {
	if len(depts) == 0 {
		return nil
	}
	ptrs := make([]*entity.Department, len(depts))
	for i := range depts {
		ptrs[i] = &depts[i]
	}
	return buildDeptTreeResponse(ptrs)
}
