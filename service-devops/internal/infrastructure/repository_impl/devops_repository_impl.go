package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-devops/internal/domain/entity"
	"leap-one/service-devops/internal/domain/repository"
	"gorm.io/gorm"
)

type RepoRepoImpl struct{ db *gorm.DB }

func NewRepositoryRepo(db *gorm.DB) repository.RepositoryRepository { return &RepoRepoImpl{db: db} }
func (r *RepoRepoImpl) Create(ctx context.Context, x *entity.Repository) error {
	return r.db.WithContext(ctx).Create(x).Error
}
func (r *RepoRepoImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Repository, error) {
	var x entity.Repository
	err := r.db.WithContext(ctx).First(&x, "id=?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}
func (r *RepoRepoImpl) List(ctx context.Context) ([]*entity.Repository, error) {
	var list []*entity.Repository
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&list).Error
	return list, err
}
func (r *RepoRepoImpl) Update(ctx context.Context, x *entity.Repository) error {
	return r.db.WithContext(ctx).Save(x).Error
}
func (r *RepoRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Repository{}, "id=?", id).Error
}

type PipelineRepoImpl struct{ db *gorm.DB }

func NewPipelineRepo(db *gorm.DB) repository.PipelineRepository { return &PipelineRepoImpl{db: db} }
func (r *PipelineRepoImpl) Create(ctx context.Context, p *entity.Pipeline) error {
	return r.db.WithContext(ctx).Create(p).Error
}
func (r *PipelineRepoImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Pipeline, error) {
	var p entity.Pipeline
	err := r.db.WithContext(ctx).First(&p, "id=?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
func (r *PipelineRepoImpl) ListByRepoID(ctx context.Context, repoID uuid.UUID) ([]*entity.Pipeline, error) {
	var list []*entity.Pipeline
	err := r.db.WithContext(ctx).Where("repo_id=?", repoID).Order("created_at DESC").Find(&list).Error
	return list, err
}
func (r *PipelineRepoImpl) Update(ctx context.Context, p *entity.Pipeline) error {
	return r.db.WithContext(ctx).Save(p).Error
}
func (r *PipelineRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Pipeline{}, "id=?", id).Error
}

type RunRepoImpl struct{ db *gorm.DB }

func NewRunRepo(db *gorm.DB) repository.PipelineRunRepository { return &RunRepoImpl{db: db} }
func (r *RunRepoImpl) Create(ctx context.Context, pr *entity.PipelineRun) error {
	return r.db.WithContext(ctx).Create(pr).Error
}
func (r *RunRepoImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.PipelineRun, error) {
	var pr entity.PipelineRun
	err := r.db.WithContext(ctx).Preload("Jobs").First(&pr, "id=?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pr, nil
}
func (r *RunRepoImpl) ListByPipelineID(ctx context.Context, pipelineID uuid.UUID) ([]*entity.PipelineRun, error) {
	var list []*entity.PipelineRun
	err := r.db.WithContext(ctx).Where("pipeline_id=?", pipelineID).Order("run_number DESC").Find(&list).Error
	return list, err
}
func (r *RunRepoImpl) Update(ctx context.Context, pr *entity.PipelineRun) error {
	return r.db.WithContext(ctx).Save(pr).Error
}

type JobRepoImpl struct{ db *gorm.DB }

func NewJobRepo(db *gorm.DB) repository.PipelineJobRepository { return &JobRepoImpl{db: db} }
func (r *JobRepoImpl) Create(ctx context.Context, j *entity.PipelineJob) error {
	return r.db.WithContext(ctx).Create(j).Error
}
func (r *JobRepoImpl) ListByRunID(ctx context.Context, runID uuid.UUID) ([]*entity.PipelineJob, error) {
	var list []*entity.PipelineJob
	err := r.db.WithContext(ctx).Where("run_id=?", runID).Order("created_at ASC").Find(&list).Error
	return list, err
}

type ArtifactRepoImpl struct{ db *gorm.DB }

func NewArtifactRepo(db *gorm.DB) repository.ArtifactRepository { return &ArtifactRepoImpl{db: db} }
func (r *ArtifactRepoImpl) Create(ctx context.Context, a *entity.Artifact) error {
	return r.db.WithContext(ctx).Create(a).Error
}
func (r *ArtifactRepoImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Artifact, error) {
	var a entity.Artifact
	err := r.db.WithContext(ctx).First(&a, "id=?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}
func (r *ArtifactRepoImpl) List(ctx context.Context) ([]*entity.Artifact, error) {
	var list []*entity.Artifact
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&list).Error
	return list, err
}
func (r *ArtifactRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Artifact{}, "id=?", id).Error
}

type DeploymentRepoImpl struct{ db *gorm.DB }

func NewDeploymentRepo(db *gorm.DB) repository.DeploymentRepository {
	return &DeploymentRepoImpl{db: db}
}
func (r *DeploymentRepoImpl) Create(ctx context.Context, d *entity.Deployment) error {
	return r.db.WithContext(ctx).Create(d).Error
}
func (r *DeploymentRepoImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Deployment, error) {
	var d entity.Deployment
	err := r.db.WithContext(ctx).First(&d, "id=?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}
func (r *DeploymentRepoImpl) List(ctx context.Context) ([]*entity.Deployment, error) {
	var list []*entity.Deployment
	err := r.db.WithContext(ctx).Order("deployed_at DESC").Find(&list).Error
	return list, err
}
func (r *DeploymentRepoImpl) Update(ctx context.Context, d *entity.Deployment) error {
	return r.db.WithContext(ctx).Save(d).Error
}

type EnvVarRepoImpl struct{ db *gorm.DB }

func NewEnvVarRepo(db *gorm.DB) repository.EnvVarRepository { return &EnvVarRepoImpl{db: db} }
func (r *EnvVarRepoImpl) Create(ctx context.Context, e *entity.EnvVar) error {
	return r.db.WithContext(ctx).Create(e).Error
}
func (r *EnvVarRepoImpl) List(ctx context.Context) ([]*entity.EnvVar, error) {
	var list []*entity.EnvVar
	err := r.db.WithContext(ctx).Order("service_name ASC,key ASC").Find(&list).Error
	return list, err
}
func (r *EnvVarRepoImpl) Update(ctx context.Context, e *entity.EnvVar) error {
	return r.db.WithContext(ctx).Save(e).Error
}
func (r *EnvVarRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.EnvVar{}, "id=?", id).Error
}
