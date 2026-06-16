package repository

import("context";"github.com/google/uuid";"leap-one/service-devops/internal/domain/entity")

type RepositoryRepository interface{Create(ctx context.Context,r*entity.Repository)error;GetByID(ctx context.Context,id uuid.UUID)(*entity.Repository,error);List(ctx context.Context)([]*entity.Repository,error);Update(ctx context.Context,r*entity.Repository)error;Delete(ctx context.Context,id uuid.UUID)error}
type PipelineRepository interface{Create(ctx context.Context,p*entity.Pipeline)error;GetByID(ctx context.Context,id uuid.UUID)(*entity.Pipeline,error);ListByRepoID(ctx context.Context,repoID uuid.UUID)([]*entity.Pipeline,error);Update(ctx context.Context,p*entity.Pipeline)error;Delete(ctx context.Context,id uuid.UUID)error}
type PipelineRunRepository interface{Create(ctx context.Context,pr*entity.PipelineRun)error;GetByID(ctx context.Context,id uuid.UUID)(*entity.PipelineRun,error);ListByPipelineID(ctx context.Context,pipelineID uuid.UUID)([]*entity.PipelineRun,error);Update(ctx context.Context,pr*entity.PipelineRun)error}
type PipelineJobRepository interface{Create(ctx context.Context,j*entity.PipelineJob)error;ListByRunID(ctx context.Context,runID uuid.UUID)([]*entity.PipelineJob,error)}
type ArtifactRepository interface{Create(ctx context.Context,a*entity.Artifact)error;GetByID(ctx context.Context,id uuid.UUID)(*entity.Artifact,error);List(ctx context.Context)([]*entity.Artifact,error);Delete(ctx context.Context,id uuid.UUID)error}
type DeploymentRepository interface{Create(ctx context.Context,d*entity.Deployment)error;GetByID(ctx context.Context,id uuid.UUID)(*entity.Deployment,error);List(ctx context.Context)([]*entity.Deployment,error);Update(ctx context.Context,d*entity.Deployment)error}
type EnvVarRepository interface{Create(ctx context.Context,e*entity.EnvVar)error;List(ctx context.Context)([]*entity.EnvVar,error);Update(ctx context.Context,e*entity.EnvVar)error;Delete(ctx context.Context,id uuid.UUID)error}
