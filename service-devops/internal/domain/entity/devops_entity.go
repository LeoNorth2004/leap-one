package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name          string         `gorm:"size:200;not null" json:"name"`
	URL           string         `gorm:"size:500;not null" json:"url"`
	Type          string         `gorm:"size:30;default:'gitlab'" json:"type"` // github/gitlab/gitee/git
	ProjectID     *uuid.UUID     `json:"project_id"`
	AuthType      string         `gorm:"size:20;default:'ssh'" json:"auth_type"` // ssh/token/password
	Credential    string         `gorm:"type:text" json:"-"`                     // 加密凭证（不序列化）
	DefaultBranch string         `gorm:"size:100;default:main" json:"default_branch"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Repository) TableName() string { return "repositories" }
func (r *Repository) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type Pipeline struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	RepoID       uuid.UUID      `gorm:"index;not null" json:"repo_id"`
	Name         string         `gorm:"size:200;not null" json:"name"`
	Type         string         `gorm:"size:30;default:'ci'" json:"type"`             // ci/cd/custom
	Config       string         `gorm:"type:text" json:"config"`                      // YAML/JSON流水线配�?
	TriggerMode  string         `gorm:"size:20;default:'manual'" json:"trigger_mode"` // manual/webhook/schedule
	ScheduleCron string         `gorm:"size:100" json:"schedule_cron"`
	ProjectID    *uuid.UUID     `json:"project_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Runs         []PipelineRun  `gorm:"foreignKey:PipelineID" json:"runs,omitempty"`
}

func (Pipeline) TableName() string { return "pipelines" }
func (p *Pipeline) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type PipelineRun struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	PipelineID  uuid.UUID      `gorm:"index;not null" json:"pipeline_id"`
	RunNumber   int            `gorm:"not null" json:"run_number"`
	Status      string         `gorm:"size:20;default:'pending'" json:"status"` // pending/running/success/failed/cancelled
	TriggeredBy uuid.UUID      `json:"triggered_by"`
	Branch      string         `gorm:"size:100" json:"branch"`
	CommitSHA   string         `gorm:"size:100" json:"commit_sha"`
	StartedAt   *time.Time     `json:"started_at"`
	FinishedAt  *time.Time     `json:"finished_at"`
	Duration    int64          `json:"duration"` // �?
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Jobs        []PipelineJob  `gorm:"foreignKey:RunID" json:"jobs,omitempty"`
}

func (PipelineRun) TableName() string { return "pipeline_runs" }
func (pr *PipelineRun) BeforeCreate(tx *gorm.DB) error {
	if pr.ID == uuid.Nil {
		pr.ID = uuid.New()
	}
	return nil
}

type PipelineJob struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RunID      uuid.UUID  `gorm:"index;not null" json:"run_id"`
	Name       string     `gorm:"size:200;not null" json:"name"`
	Stage      string     `gorm:"size:50" json:"stage"`
	Status     string     `gorm:"size:20;default:'pending'" json:"status"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	LogContent string     `gorm:"type:text" json:"log_content"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (PipelineJob) TableName() string { return "pipeline_jobs" }
func (pj *PipelineJob) BeforeCreate(tx *gorm.DB) error {
	if pj.ID == uuid.Nil {
		pj.ID = uuid.New()
	}
	return nil
}

type Artifact struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	RunID       uuid.UUID      `gorm:"index" json:"run_id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Type        string         `gorm:"size:30;default:'binary'" json:"type"` // binary/docker/image/package
	Version     string         `gorm:"size:100;not null" json:"version"`
	Size        int64          `json:"size"`
	StoragePath string         `gorm:"size:500" json:"storage_path"`
	DownloadURL string         `gorm:"size:500" json:"download_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Artifact) TableName() string { return "artifacts" }
func (a *Artifact) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type Deployment struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Environment  string         `gorm:"size:50;not null" json:"environment"` // dev/test/staging/prod
	ArtifactID   *uuid.UUID     `json:"artifact_id"`
	ProjectID    *uuid.UUID     `json:"project_id"`
	Status       string         `gorm:"size:20;default:'deploying'" json:"status"` // deploying/success/failed/rolling_back
	DeployedBy   uuid.UUID      `gorm:"not null" json:"deployed_by"`
	DeployedAt   time.Time      `json:"deployed_at"`
	Version      string         `gorm:"size:100" json:"version"`
	Notes        string         `gorm:"type:text" json:"notes"`
	RollbackFrom *uuid.UUID     `json:"rollback_from"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Deployment) TableName() string { return "deployments" }
func (d *Deployment) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

type EnvVar struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ServiceName string         `gorm:"size:100;not null;index" json:"service_name"`
	Key         string         `gorm:"size:200;not null" json:"key"`
	Value       string         `gorm:"type;text" json:"value"`
	IsEncrypted bool           `gorm:"default:false" json:"is_encrypted"`
	Environment string         `gorm:"size:20;default:'all'" json:"environment"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EnvVar) TableName() string { return "env_vars" }
func (e *EnvVar) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
