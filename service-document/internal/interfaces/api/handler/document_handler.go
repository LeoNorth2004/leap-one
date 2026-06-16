package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-document/internal/application/service"
	"leap-one/service-document/internal/domain/entity"
	"leap-one/service-document/internal/domain/repository"
	"leap-one/service-document/internal/interfaces/api/dto"
)

// DocumentHandler 文档HTTP处理�?
type DocumentHandler struct {
	docSvc      *service.DocumentService
	versionSvc  *service.VersionService
	commentSvc  *service.CommentService
	categorySvc *service.CategoryService
	kbSvc       *service.KnowledgeBaseService
	attachSvc   *service.AttachmentService
	tagSvc      *service.TagService
	logger      *zap.Logger
}

func NewDocumentHandler(
	docSvc *service.DocumentService, versionSvc *service.VersionService,
	commentSvc *service.CommentService, categorySvc *service.CategoryService,
	kbSvc *service.KnowledgeBaseService, attachSvc *service.AttachmentService,
	tagSvc *service.TagService, logger *zap.Logger,
) *DocumentHandler {
	return &DocumentHandler{docSvc, versionSvc, commentSvc, categorySvc, kbSvc, attachSvc, tagSvc, logger}
}

// ==================== 文档CRUD ====================

func (h *DocumentHandler) Create(c *gin.Context) {
	var req dto.CreateDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest("参数校验失败: "+err.Error()))
		return
	}
	doc := &entity.Document{
		ID: uuid.New(), Title: req.Title, Content: req.Content, Type: req.Type,
		CategoryID: req.CategoryID, ParentID: req.ParentID, ProductID: req.ProductID,
		ProjectID: req.ProjectID, OwnerID: req.OwnerID, Visibility: req.Visibility,
		Tags: req.Tags, IsTemplate: req.IsTemplate,
	}
	result, err := h.docSvc.Create(doc)
	if err != nil {
		h.logger.Error("创建文档失败", zap.Error(err))
		c.JSON(500, dto.InternalError("创建文档失败"))
		return
	}
	c.JSON(201, dto.Success(toDocResp(result)))
}

func (h *DocumentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	params := &repository.DocumentListParams{
		Page: page, PageSize: size, Type: c.Query("type"), Status: c.Query("status"),
		Visibility: c.Query("visibility"), Keyword: c.Query("keyword"), SortBy: c.Query("sort_by"), SortOrder: c.Query("sort_order"),
	}
	if pid := c.Query("project_id"); pid != "" {
		if u, e := uuid.Parse(pid); e == nil {
			params.ProjectID = &u
		}
	}
	list, total, err := h.docSvc.List(params)
	if err != nil {
		c.JSON(500, dto.InternalError("查询失败"))
		return
	}
	var resps []dto.DocumentResponse
	for _, d := range list {
		resps = append(resps, toDocResp(d))
	}
	c.JSON(200, dto.PageSuccess(resps, total, page, size))
}

func (h *DocumentHandler) GetTree(c *gin.Context) {
	pid, err := uuid.Parse(c.Query("project_id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效项目ID"))
		return
	}
	tree, err := h.docSvc.GetTree(pid)
	if err != nil {
		c.JSON(500, dto.InternalError("获取目录树失�?))
		return
	}
	var resps []dto.DocumentResponse
	for _, d := range tree {
		resps = append(resps, toDocRespWithChildren(d))
	}
	c.JSON(200, dto.Success(resps))
}

func (h *DocumentHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	doc, err := h.docSvc.GetByID(id)
	if err != nil {
		c.JSON(404, dto.NotFound("文档不存�?))
		return
	}
	c.JSON(200, dto.Success(toDocResp(doc)))
}

func (h *DocumentHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	var req dto.UpdateDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest("参数校验失败: "+err.Error()))
		return
	}
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Visibility != nil {
		updates["visibility"] = *req.Visibility
	}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Tags != nil {
		updates["tags"] = *req.Tags
	}
	result, err := h.docSvc.Update(id, updates)
	if err != nil {
		c.JSON(500, dto.InternalError("更新失败"))
		return
	}
	c.JSON(200, dto.Success(toDocResp(result)))
}

func (h *DocumentHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	if err := h.docSvc.Delete(id); err != nil {
		c.JSON(500, dto.InternalError("删除失败"))
		return
	}
	c.JSON(200, dto.Success(nil))
}

// ==================== 发布/版本/评论/收藏/附件 ====================

func (h *DocumentHandler) Publish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	if err := h.docSvc.Publish(id); err != nil {
		c.JSON(500, dto.InternalError("发布失败"))
		return
	}
	c.JSON(200, dto.Success(gin.H{"status": "published"}))
}

func (h *DocumentHandler) ListVersions(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	versions, err := h.versionSvc.ListVersions(id)
	if err != nil {
		c.JSON(500, dto.InternalError("获取版本失败"))
		return
	}
	var resps []dto.VersionResponse
	for _, v := range versions {
		resps = append(resps, toVersionResp(v))
	}
	c.JSON(200, dto.Success(resps))
}

func (h *DocumentHandler) GetVersion(c *gin.Context) {
	docID, _ := uuid.Parse(c.Param("id"))
	ver, _ := strconv.Atoi(c.Param("vid"))
	v, err := h.versionSvc.GetVersion(docID, ver)
	if err != nil {
		c.JSON(404, dto.NotFound("版本不存�?))
		return
	}
	c.JSON(200, dto.Success(toVersionResp(v)))
}

func (h *DocumentHandler) Restore(c *gin.Context) {
	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	targetVer, _ := strconv.Atoi(c.Query("version"))
	doc, err := h.docSvc.GetByID(docID)
	if err != nil {
		c.JSON(404, dto.NotFound("文档不存�?))
		return
	}
	if err := h.versionSvc.RestoreToVersion(docID, targetVer, doc, h.docSvc); err != nil {
		c.JSON(500, dto.InternalError("恢复失败"))
		return
	}
	c.JSON(200, dto.Success(gin.H{"restored_to_version": targetVer}))
}

func (h *DocumentHandler) AddComment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	var req dto.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	comment := &entity.DocumentComment{ID: uuid.New(), DocumentID: id, UserID: req.UserID, Content: req.Content, Position: req.Position, ParentID: req.ParentID}
	if err := h.commentSvc.Add(comment); err != nil {
		c.JSON(500, dto.InternalError("添加评论失败"))
		return
	}
	c.JSON(201, dto.Success(dto.CommentResponse{ID: comment.ID, DocumentID: id, UserID: req.UserID, Content: req.Content}))
}

func (h *DocumentHandler) ListComments(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	comments, err := h.commentSvc.List(id)
	if err != nil {
		c.JSON(500, dto.InternalError("获取评论失败"))
		return
	}
	var resps []dto.CommentResponse
	for _, cm := range comments {
		resps = append(resps, toCommentResp(cm))
	}
	c.JSON(200, dto.Success(resps))
}

func (h *DocumentHandler) Favorite(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetHeader("X-User-ID"))
	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	if err := h.docSvc.AddFavorite(userID, docID); err != nil {
		c.JSON(500, dto.InternalError("收藏失败"))
		return
	}
	c.JSON(200, dto.Success(gin.H{"favorited": true}))
}

func (h *DocumentHandler) Unfavorite(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetHeader("X-User-ID"))
	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	if err := h.docSvc.RemoveFavorite(userID, docID); err != nil {
		c.JSON(500, dto.InternalError("取消收藏失败"))
		return
	}
	c.JSON(200, dto.Success(gin.H{"favorited": false}))
}

func (h *DocumentHandler) UploadAttachment(c *gin.Context) {
	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	var req dto.AttachRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	attach := &entity.DocumentAttachment{ID: uuid.New(), DocumentID: docID, FileName: req.FileName, FileSize: req.FileSize, FileType: req.FileType, FileURL: req.FileURL}
	if err := h.attachSvc.Upload(attach); err != nil {
		c.JSON(500, dto.InternalError("上传失败"))
		return
	}
	c.JSON(201, dto.Success(dto.AttachResponse{ID: attach.ID, DocumentID: docID, FileName: req.FileName, FileSize: req.FileSize, FileType: req.FileType, FileURL: req.FileURL}))
}

func (h *DocumentHandler) Search(c *gin.Context) {
	docs, err := h.docSvc.Search(c.Query("q"))
	if err != nil {
		c.JSON(500, dto.InternalError("搜索失败"))
		return
	}
	var resps []dto.DocumentResponse
	for _, d := range docs {
		resps = append(resps, toDocResp(d))
	}
	c.JSON(200, dto.Success(resps))
}

// ==================== 分类 CRUD ====================

func (h *DocumentHandler) CreateCategory(c *gin.Context) {
	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	cat := &entity.DocumentCategory{ID: uuid.New(), Name: req.Name, ParentID: req.ParentID, SortOrder: req.SortOrder}
	if err := h.categorySvc.Create(cat); err != nil {
		c.JSON(500, dto.InternalError("创建分类失败"))
		return
	}
	c.JSON(201, dto.Success(dto.CategoryResponse{ID: cat.ID, Name: req.Name, ParentID: req.ParentID, SortOrder: req.SortOrder}))
}

func (h *DocumentHandler) ListCategories(c *gin.Context) {
	cats, err := h.categorySvc.List()
	if err != nil {
		c.JSON(500, dto.InternalError("获取分类失败"))
		return
	}
	var resps []dto.CategoryResponse
	for _, cat := range cats {
		resps = append(resps, dto.CategoryResponse{ID: cat.ID, Name: cat.Name, ParentID: cat.ParentID, SortOrder: cat.SortOrder})
	}
	c.JSON(200, dto.Success(resps))
}

func (h *DocumentHandler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	cat := &entity.DocumentCategory{ID: id, Name: req.Name, ParentID: req.ParentID, SortOrder: req.SortOrder}
	if err := h.categorySvc.Update(cat); err != nil {
		c.JSON(500, dto.InternalError("更新分类失败"))
		return
	}
	c.JSON(200, dto.Success(nil))
}

func (h *DocumentHandler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	if err := h.categorySvc.Delete(id); err != nil {
		c.JSON(500, dto.InternalError("删除分类失败"))
		return
	}
	c.JSON(200, dto.Success(nil))
}

// ==================== 知识�?CRUD ====================

func (h *DocumentHandler) CreateKB(c *gin.Context) {
	var req dto.KBRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	ownerID, _ := uuid.Parse(c.GetHeader("X-User-ID"))
	kb := &entity.KnowledgeBase{ID: uuid.New(), Name: req.Name, Description: req.Description, OwnerID: ownerID, IsPublic: req.IsPublic}
	if err := h.kbSvc.Create(kb); err != nil {
		c.JSON(500, dto.InternalError("创建知识库失�?))
		return
	}
	c.JSON(201, dto.Success(dto.KBResponse{ID: kb.ID, Name: req.Name, Description: req.Description, OwnerID: ownerID, IsPublic: req.IsPublic}))
}

func (h *DocumentHandler) ListKBs(c *gin.Context) {
	ownerID, _ := uuid.Parse(c.DefaultQuery("owner_id", "00000000-0000-0000-0000-000000000000"))
	kbs, err := h.kbSvc.List(ownerID)
	if err != nil {
		c.JSON(500, dto.InternalError("获取知识库列表失�?))
		return
	}
	var resps []dto.KBResponse
	for _, kb := range kbs {
		resps = append(resps, dto.KBResponse{ID: kb.ID, Name: kb.Name, Description: kb.Description, OwnerID: kb.OwnerID, IsPublic: kb.IsPublic})
	}
	c.JSON(200, dto.Success(resps))
}

func (h *DocumentHandler) GetKB(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	kb, err := h.kbSvc.GetByID(id)
	if err != nil {
		c.JSON(404, dto.NotFound("知识库不存在"))
		return
	}
	c.JSON(200, dto.Success(dto.KBResponse{ID: kb.ID, Name: kb.Name, Description: kb.Description, OwnerID: kb.OwnerID, IsPublic: kb.IsPublic}))
}

func (h *DocumentHandler) UpdateKB(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	var req dto.KBRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	kb := &entity.KnowledgeBase{ID: id, Name: req.Name, Description: req.Description, IsPublic: req.IsPublic}
	if err := h.kbSvc.Update(kb); err != nil {
		c.JSON(500, dto.InternalError("更新知识库失�?))
		return
	}
	c.JSON(200, dto.Success(nil))
}

func (h *DocumentHandler) DeleteKB(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	if err := h.kbSvc.Delete(id); err != nil {
		c.JSON(500, dto.InternalError("删除知识库失�?))
		return
	}
	c.JSON(200, dto.Success(nil))
}

// ==================== 模板 CRUD（复用文档接口）====================

func (h *DocumentHandler) ListTemplates(c *gin.Context) {
	params := &repository.DocumentListParams{Page: 1, PageSize: 100, IsTemplate: boolPtr(true)}
	list, _, _ := h.docSvc.List(params)
	var resps []dto.DocumentResponse
	for _, d := range list {
		resps = append(resps, toDocResp(d))
	}
	c.JSON(200, dto.Success(resps))
}

// 转换函数
func toDocResp(d *entity.Document) dto.DocumentResponse {
	return dto.DocumentResponse{ID: d.ID, Title: d.Title, Content: d.Content, Type: d.Type, CategoryID: d.CategoryID, ParentID: d.ParentID, ProductID: d.ProductID, ProjectID: d.ProjectID, OwnerID: d.OwnerID, Status: d.Status, Visibility: d.Visibility, Version: d.Version, Tags: d.Tags, IsTemplate: d.IsTemplate, OrderIndex: d.OrderIndex, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt}
}
func toDocRespWithChildren(d *entity.Document) dto.DocumentResponse {
	r := toDocResp(d)
	for i := range d.Children {
		r.Children = append(r.Children, toDocRespWithChildren(&d.Children[i]))
	}
	return r
}
func toVersionResp(v *entity.DocumentVersion) dto.VersionResponse {
	return dto.VersionResponse{ID: v.ID, DocumentID: v.DocumentID, VersionNo: v.VersionNo, Title: v.Title, Content: v.Content, ChangeNote: v.ChangeNote, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt}
}
func toCommentResp(cm *entity.DocumentComment) dto.CommentResponse {
	return dto.CommentResponse{ID: cm.ID, DocumentID: cm.DocumentID, UserID: cm.UserID, Content: cm.Content, Position: cm.Position, CreatedAt: cm.CreatedAt}
}
func boolPtr(b bool) *bool { return &b }
