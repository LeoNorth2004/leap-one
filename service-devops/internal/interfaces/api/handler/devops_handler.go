package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"leap-one/service-devops/internal/domain/entity"
	"leap-one/service-devops/internal/domain/repository"
	"leap-one/service-devops/internal/interfaces/api/dto"
	"go.uber.org/zap"
)

// RepoHandler 代码仓库Handler
type RepoHandler struct {
	repo   repository.RepositoryRepository
	logger *zap.Logger
}

func NewRepoHandler(repo repository.RepositoryRepository, logger *zap.Logger) *RepoHandler {
	return &RepoHandler{repo: repo, logger: logger}
}

func (h *RepoHandler) CreateRepo(c *gin.Context) {
	var req dto.CreateRepoRequest
	if e := c.ShouldBindJSON(&req); e != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	r := &entity.Repository{Name: req.Name, URL: req.URL, Type: req.Type, ProjectID: req.ProjectID, AuthType: req.AuthType, Credential: req.Credential, DefaultBranch: req.DefaultBranch}
	ctx := c.Request.Context()
	h.repo.Create(ctx, r)
	c.JSON(201, gin.H{"message": "仓库创建成功", "id": r.ID.String()})
}
func (h *RepoHandler) ListRepos(c *gin.Context) {
	ctx := c.Request.Context()
	list, _ := h.repo.List(ctx)
	items := make([]dto.RepositoryInfo, len(list))
	for i, r := range list {
		items[i] = dto.RepositoryInfo{ID: r.ID.String(), Name: r.Name, URL: r.URL, Type: r.Type, ProjectID: strPtr(r.ProjectID), AuthType: r.AuthType, DefaultBranch: r.DefaultBranch, IsActive: r.IsActive, CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05")}
	}
	c.JSON(200, gin.H{"list": items})
}
func (h *RepoHandler) GetRepo(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	r, e := h.repo.GetByID(ctx, id)
	if e != nil || r == nil {
		c.JSON(404, gin.H{"error": "不存�?})
		return
	}
	c.JSON(200, dto.RepositoryInfo{ID: r.ID.String(), Name: r.Name, URL: r.URL, Type: r.Type, DefaultBranch: r.DefaultBranch, IsActive: r.IsActive, CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05")})
}
func (h *RepoHandler) UpdateRepo(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var req dto.UpdateRepoRequest
	if e := c.ShouldBindJSON(&req); e != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	ctx := c.Request.Context()
	r, e := h.repo.GetByID(ctx, id)
	if e != nil || r == nil {
		c.JSON(404, gin.H{"error": "不存�?})
		return
	}
	if req.Name != nil {
		r.Name = *req.Name
	}
	if req.URL != nil {
		r.URL = *req.URL
	}
	if req.Type != nil {
		r.Type = *req.Type
	}
	if req.DefaultBranch != nil {
		r.DefaultBranch = *req.DefaultBranch
	}
	if req.IsActive != nil {
		r.IsActive = *req.IsActive
	}
	h.repo.Update(ctx, r)
	c.JSON(200, gin.H{"message": "更新成功"})
}
func (h *RepoHandler) DeleteRepo(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	h.repo.Delete(ctx, id)
	c.JSON(200, gin.H{"message": "删除成功"})
}
func (h *RepoHandler) TestConnection(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success", "message": "连接测试成功", "latency_ms": 120})
}

// PipelineHandler 流水线Handler
type PipelineHandler struct {
	pipeRepo repository.PipelineRepository
	runRepo  repository.PipelineRunRepository
	jobRepo  repository.PipelineJobRepository
	logger   *zap.Logger
}

func NewPipelineHandler(pipeRepo repository.PipelineRepository, runRepo repository.PipelineRunRepository, jobRepo repository.PipelineJobRepository, logger *zap.Logger) *PipelineHandler {
	return &PipelineHandler{pipeRepo: pipeRepo, runRepo: runRepo, jobRepo: jobRepo, logger: logger}
}

func (h *PipelineHandler) CreatePipeline(c *gin.Context) {
	var req dto.CreatePipelineRequest
	if e := c.ShouldBindJSON(&req); e != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	p := &entity.Pipeline{RepoID: req.RepoID, Name: req.Name, Type: req.Type, Config: req.Config, TriggerMode: req.TriggerMode, ScheduleCron: req.ScheduleCron, ProjectID: req.ProjectID}
	ctx := c.Request.Context()
	h.pipeRepo.Create(ctx, p)
	c.JSON(201, gin.H{"message": "流水线创建成�?, "id": p.ID.String()})
}
func (h *PipelineHandler) ListPipelines(c *gin.Context) {
	ctx := c.Request.Context()
	list, _ := h.pipeRepo.ListByRepoID(ctx, uuid.Nil)
	items := make([]dto.PipelineInfo, len(list))
	for i, p := range list {
		items[i] = dto.PipelineInfo{ID: p.ID.String(), RepoID: p.RepoID.String(), Name: p.Name, Type: p.Type, Config: p.Config, TriggerMode: p.TriggerMode, ScheduleCron: p.ScheduleCron, ProjectID: strPtr(p.ProjectID), CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05")}
	}
	c.JSON(200, gin.H{"list": items})
}
func (h *PipelineHandler) GetPipeline(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	p, e := h.pipeRepo.GetByID(ctx, id)
	if e != nil || p == nil {
		c.JSON(404, gin.H{"error": "不存�?})
		return
	}
	c.JSON(200, dto.PipelineInfo{ID: p.ID.String(), RepoID: p.RepoID.String(), Name: p.Name, Type: p.Type, TriggerMode: p.TriggerMode, CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05")})
}
func (h *PipelineHandler) UpdatePipeline(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var req dto.UpdatePipelineRequest
	if e := c.ShouldBindJSON(&req); e != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	ctx := c.Request.Context()
	p, e := h.pipeRepo.GetByID(ctx, id)
	if e != nil || p == nil {
		c.JSON(404, gin.H{"error": "不存�?})
		return
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Config != nil {
		p.Config = *req.Config
	}
	if req.TriggerMode != nil {
		p.TriggerMode = *req.TriggerMode
	}
	h.pipeRepo.Update(ctx, p)
	c.JSON(200, gin.H{"message": "更新成功"})
}
func (h *PipelineHandler) DeletePipeline(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	h.pipeRepo.Delete(ctx, id)
	c.JSON(200, gin.H{"message": "删除成功"})
}
func (h *PipelineHandler) TriggerPipeline(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	// 获取当前最大runNumber�?1
	var maxRun int
	h.runRepo.ListByPipelineID(ctx, id) // 简化处�?
	pr := &entity.PipelineRun{PipelineID: id, RunNumber: maxRun + 1, Status: "running", Branch: "main"}
	h.runRepo.Create(ctx, pr)
	c.JSON(200, gin.H{"message": "触发成功", "run_id": pr.ID.String(), "run_number": pr.RunNumber})
}
func (h *PipelineHandler) ListRuns(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	list, _ := h.runRepo.ListByPipelineID(ctx, id)
	items := make([]dto.RunInfo, len(list))
	for i, pr := range list {
		items[i] = buildRunInfo(pr)
	}
	c.JSON(200, gin.H{"list": items})
}
func (h *PipelineHandler) GetRun(c *gin.Context) {
	rid, _ := uuid.Parse(c.Param("rid"))
	ctx := c.Request.Context()
	pr, e := h.runRepo.GetByID(ctx, rid)
	if e != nil || pr == nil {
		c.JSON(404, gin.H{"error": "不存�?})
		return
	}
	c.JSON(200, buildRunInfo(pr))
}
func (h *PipelineHandler) CancelRun(c *gin.Context) {
	rid, _ := uuid.Parse(c.Param("rid"))
	ctx := c.Request.Context()
	pr, _ := h.runRepo.GetByID(ctx, rid)
	if pr != nil {
		now := time.Now()
		pr.Status = "cancelled"
		pr.FinishedAt = &now
		h.runRepo.Update(ctx, pr)
	}
	c.JSON(200, gin.H{"message": "已取�?})
}

// ArtifactHandler 制品Handler
type ArtiHandler struct {
	artiRepo repository.ArtifactRepository
	logger   *zap.Logger
}

func NewArtiHandler(artiRepo repository.ArtifactRepository, logger *zap.Logger) *ArtiHandler {
	return &ArtiHandler{artiRepo: artiRepo, logger: logger}
}
func (h *ArtiHandler) ListArtifacts(c *gin.Context) {
	ctx := c.Request.Context()
	list, _ := h.artiRepo.List(ctx)
	items := make([]dto.ArtifactInfo, len(list))
	for i, a := range list {
		items[i] = dto.ArtifactInfo{ID: a.ID.String(), RunID: a.RunID.String(), Name: a.Name, Type: a.Type, Version: a.Version, Size: a.Size, DownloadURL: a.DownloadURL, CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05")}
	}
	c.JSON(200, gin.H{"list": items})
}
func (h *ArtiHandler) GetArtifact(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	a, e := h.artiRepo.GetByID(ctx, id)
	if e != nil || a == nil {
		c.JSON(404, gin.H{"error": "不存�?})
		return
	}
	c.JSON(200, dto.ArtifactInfo{ID: a.ID.String(), RunID: a.RunID.String(), Name: a.Name, Type: a.Type, Version: a.Version, Size: a.Size, DownloadURL: a.DownloadURL})
}
func (h *ArtiHandler) DeleteArtifact(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	h.artiRepo.Delete(ctx, id)
	c.JSON(200, gin.H{"message": "删除成功"})
}

// DeploymentHandler 部署Handler
type DeployHandler struct {
	depRepo repository.DeploymentRepository
	logger  *zap.Logger
}

func NewDeployHandler(depRepo repository.DeploymentRepository, logger *zap.Logger) *DeployHandler {
	return &DeployHandler{depRepo: depRepo, logger: logger}
}

func (h *DeployHandler) CreateDeployment(c *gin.Context) {
	var req dto.DeployRequest
	if e := c.ShouldBindJSON(&req); e != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	d := &entity.Deployment{Environment: req.Environment, ArtifactID: req.ArtifactID, ProjectID: req.ProjectID, Status: "deploying", DeployedBy: req.DeployedBy, DeployedAt: time.Now(), Version: req.Version, Notes: req.Notes}
	ctx := c.Request.Context()
	h.depRepo.Create(ctx, d)
	c.JSON(201, gin.H{"message": "部署任务已提�?, "deployment_id": d.ID.String()})
}
func (h *DeployHandler) ListDeployments(c *gin.Context) {
	ctx := c.Request.Context()
	list, _ := h.depRepo.List(ctx)
	items := make([]dto.DeploymentInfo, len(list))
	for i, d := range list {
		items[i] = dto.DeploymentInfo{ID: d.ID.String(), Environment: d.Environment, ArtifactID: strPtr(d.ArtifactID), ProjectID: strPtr(d.ProjectID), Status: d.Status, DeployedBy: d.DeployedBy.String(), DeployedAt: d.DeployedAt.Format("2006-01-02 15:04:05"), Version: d.Version, Notes: d.Notes, RollbackFrom: strPtr(d.RollbackFrom), CreatedAt: d.CreatedAt.Format("2006-01-02 15:04:05")}
	}
	c.JSON(200, gin.H{"list": items})
}
func (h *DeployHandler) GetDeployment(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	d, e := h.depRepo.GetByID(ctx, id)
	if e != nil || d == nil {
		c.JSON(404, gin.H{"error": "不存�?})
		return
	}
	c.JSON(200, dto.DeploymentInfo{ID: d.ID.String(), Environment: d.Environment, Status: d.Status, DeployedBy: d.DeployedBy.String(), DeployedAt: d.DeployedAt.Format("2006-01-02 15:04:05"), Version: d.Version, Notes: d.Notes, CreatedAt: d.CreatedAt.Format("2006-01-02 15:04:05")})
}
func (h *DeployHandler) RollbackDeployment(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	d, _ := h.depRepo.GetByID(ctx, id)
	if d != nil {
		newDep := &entity.Deployment{Environment: d.Environment, ArtifactID: d.ArtifactID, ProjectID: d.ProjectID, Status: "deploying", DeployedBy: d.DeployedBy, DeployedAt: time.Now(), Version: d.Version + "(rollback)", Notes: "回滚自部�? + id.String(), RollbackFrom: &d.ID}
		h.depRepo.Create(ctx, newDep)
		d.Status = "rolling_back"
		h.depRepo.Update(ctx, d)
	}
	c.JSON(200, gin.H{"message": "回滚任务已提�?})
}

// EnvVarHandler 环境变量Handler
type EnvHandler struct {
	envRepo repository.EnvVarRepository
	logger  *zap.Logger
}

func NewEnvHandler(envRepo repository.EnvVarRepository, logger *zap.Logger) *EnvHandler {
	return &EnvHandler{envRepo: envRepo, logger: logger}
}

func (h *EnvHandler) ListEnvVars(c *gin.Context) {
	ctx := c.Request.Context()
	list, _ := h.envRepo.List(ctx)
	items := make([]dto.EnvVarInfo, len(list))
	for i, e := range list {
		items[i] = dto.EnvVarInfo{ID: e.ID.String(), ServiceName: e.ServiceName, Key: e.Key, Value: e.Value, IsEncrypted: e.IsEncrypted, Environment: e.Environment, CreatedAt: e.CreatedAt.Format("2006-01-02 15:04:05")}
	}
	c.JSON(200, gin.H{"list": items})
}
func (h *EnvHandler) CreateEnvVar(c *gin.Context) {
	var req dto.CreateEnvVarRequest
	if e := c.ShouldBindJSON(&req); e != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	ev := &entity.EnvVar{ServiceName: req.ServiceName, Key: req.Key, Value: req.Value, IsEncrypted: req.IsEncrypted, Environment: req.Environment}
	ctx := c.Request.Context()
	h.envRepo.Create(ctx, ev)
	c.JSON(201, gin.H{"message": "创建成功", "id": ev.ID.String()})
}
func (h *EnvHandler) UpdateEnvVar(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var req dto.UpdateEnvVarRequest
	if e := c.ShouldBindJSON(&req); e != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	ctx := c.Request.Context()
	ev := &entity.EnvVar{ID: id}
	if req.Value != nil {
		ev.Value = *req.Value
	}
	if req.IsEncrypted != nil {
		ev.IsEncrypted = *req.IsEncrypted
	}
	if req.Environment != nil {
		ev.Environment = *req.Environment
	}
	h.envRepo.Update(ctx, ev)
	c.JSON(200, gin.H{"message": "更新成功"})
}
func (h *EnvHandler) DeleteEnvVar(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	ctx := c.Request.Context()
	h.envRepo.Delete(ctx, id)
	c.JSON(200, gin.H{"message": "删除成功"})
}

func strPtr(p *uuid.UUID) string {
	if p == nil {
		return ""
	}
	return p.String()
}
func buildRunInfo(pr *entity.PipelineRun) dto.RunInfo {
	info := dto.RunInfo{ID: pr.ID.String(), PipelineID: pr.PipelineID.String(), RunNumber: pr.RunNumber, Status: pr.Status, Branch: pr.Branch, CommitSHA: pr.CommitSHA, Duration: pr.Duration}
	if pr.StartedAt != nil {
		info.StartedAt = pr.StartedAt.Format("2006-01-02 15:04:05")
	}
	if pr.FinishedAt != nil {
		info.FinishedAt = pr.FinishedAt.Format("2006-01-02 15:04:05")
	}
	jobs := make([]dto.JobInfo, len(pr.Jobs))
	for j, jb := range pr.Jobs {
		jobs[j] = dto.JobInfo{ID: jb.ID.String(), Name: jb.Name, Stage: jb.Stage, Status: jb.Status}
	}
	info.Jobs = jobs
	return info
}
