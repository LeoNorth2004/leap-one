package application

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-devops/internal/domain/entity"
	"leap-one/service-devops/internal/domain/repository"
	"go.uber.org/zap"
)

type DevOpsService struct {
	repoRepo repository.RepositoryRepository
	pipeRepo repository.PipelineRepository
	runRepo  repository.PipelineRunRepository
	logger   *zap.Logger
}

func NewDevOpsService(repoRepo repository.RepositoryRepository, pipeRepo repository.PipelineRepository, runRepo repository.PipelineRunRepository, logger *zap.Logger) *DevOpsService {
	return &DevOpsService{repoRepo: repoRepo, pipeRepo: pipeRepo, runRepo: runRepo, logger: logger}
}
func (s *DevOpsService) TriggerPipelineUseCase(ctx context.Context, pipelineID uuid.UUID, triggeredBy uuid.UUID, branch string) (*entity.PipelineRun, error) {
	// 获取当前最大runNumber
	runs, _ := s.runRepo.ListByPipelineID(ctx, pipelineID)
	maxNum := 0
	for _, r := range runs {
		if r.RunNumber > maxNum {
			maxNum = r.RunNumber
		}
	}
	pr := &entity.PipelineRun{PipelineID: pipelineID, RunNumber: maxNum + 1, Status: "running", TriggeredBy: triggeredBy, Branch: branch}
	if err := s.runRepo.Create(ctx, pr); err != nil {
		return nil, err
	}
	s.logger.Info("流水线触发成�?, zap.String("pipeline_id", pipelineID.String()), zap.Int("run_number", pr.RunNumber))
	return pr, nil
}
func (s *DevOpsService) CreateDeploymentUseCase(ctx context.Context, environment string, artifactID *uuid.UUID, projectID *uuid.UUID, deployedBy uuid.UUID, version, notes string) (*entity.Deployment, error) {
	d := &entity.Deployment{Environment: environment, ArtifactID: artifactID, ProjectID: projectID, Status: "deploying", DeployedBy: deployedBy, Version: version, Notes: notes}
	if err := createDeploy(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

// createDeploy 辅助函数（避免循环依赖）
var createDeploy func(ctx context.Context, d *entity.Deployment) error
