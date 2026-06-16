package dto

import "github.com/google/uuid"

// RepositoryInfo 仓库信息
type RepositoryInfo struct {
	ID            string
	Name          string
	URL           string
	Type          string
	ProjectID     string
	AuthType      string
	DefaultBranch string
	IsActive      bool
	CreatedAt     string
}

// CreateRepoRequest 创建仓库请求
type CreateRepoRequest struct {
	Name          string     `json:"name" binding:"required,max=200"`
	URL           string     `json:"url" binding:"required,max=500"`
	Type          string     `json:"type" binding:"omitempty,oneof=github gitlab gitee git"`
	ProjectID     *uuid.UUID `json:"project_id"`
	AuthType      string     `json:"auth_type"`
	Credential    string     `json:"credential"`
	DefaultBranch string     `json:"default_branch"`
}

// UpdateRepoRequest 更新仓库请求
type UpdateRepoRequest struct {
	Name          *string
	URL           *string
	Type          *string
	AuthType      *string
	Credential    *string
	DefaultBranch *string
	IsActive      *bool
}

// PipelineInfo 流水线信�?
type PipelineInfo struct {
	ID           string
	RepoID       string
	Name         string
	Type         string
	Config       string
	TriggerMode  string
	ScheduleCron string
	ProjectID    string
	CreatedAt    string
}

// CreatePipelineRequest 创建流水线请�?
type CreatePipelineRequest struct {
	RepoID       uuid.UUID  `json:"repo_id binding:"required"`
	Name         string     `json:"name binding:"required,max=200"`
	Type         string     `json:"type"`
	Config       string     `json:"config"`
	TriggerMode  string     `json:"trigger_mode"`
	ScheduleCron string     `json:"schedule_cron"`
	ProjectID    *uuid.UUID `json:"project_id"`
}

// UpdatePipelineRequest 更新流水线请�?
type UpdatePipelineRequest struct {
	Name         *string
	Type         *string
	Config       *string
	TriggerMode  *string
	ScheduleCron *string
}

// RunInfo 执行记录信息
type RunInfo struct {
	ID          string
	PipelineID  string
	RunNumber   int
	Status      string
	TriggeredBy string
	Branch      string
	CommitSHA   string
	StartedAt   string
	FinishedAt  string
	Duration    int64
	Jobs        []JobInfo `json:"jobs,omitempty"`
}

// JobInfo Job信息
type JobInfo struct {
	ID         string
	Name       string
	Stage      string
	Status     string
	StartedAt  string
	FinishedAt string
}

// ArtifactInfo 制品信息
type ArtifactInfo struct {
	ID          string
	RunID       string
	Name        string
	Type        string
	Version     string
	Size        int64
	DownloadURL string
	CreatedAt   string
}

// DeploymentInfo 部署信息
type DeploymentInfo struct {
	ID           string
	Environment  string
	ArtifactID   string
	ProjectID    string
	Status       string
	DeployedBy   string
	DeployedAt   string
	Version      string
	Notes        string
	RollbackFrom string
	CreatedAt    string
}

// DeployRequest 执行部署请求
type DeployRequest struct {
	Environment string     `json:"environment binding:"required,oneof=dev test staging prod"`
	ArtifactID  *uuid.UUID `json:"artifact_id"`
	ProjectID   *uuid.UUID `json:"project_id"`
	DeployedBy  uuid.UUID  `json:"deployed_by binding:"required"`
	Version     string     `json:"version"`
	Notes       string     `json:"notes"`
}

// EnvVarInfo 环境变量信息
type EnvVarInfo struct {
	ID          string
	ServiceName string
	Key         string
	Value       string
	IsEncrypted bool
	Environment string
	CreatedAt   string
}

// CreateEnvVarRequest 创建环境变量请求
type CreateEnvVarRequest struct {
	ServiceName string `json:"service_name binding:"required,max=100"`
	Key         string `json:"key binding:"required,max=200"`
	Value       string `json:"value"`
	IsEncrypted bool   `json:"is_encrypted"`
	Environment string `json:"environment"`
}

// UpdateEnvVarRequest 更新环境变量请求
type UpdateEnvVarRequest struct {
	Value       *string
	IsEncrypted *bool
	Environment *string
}
