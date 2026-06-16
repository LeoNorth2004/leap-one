package application

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-config/internal/domain/entity"
	"leap-one/service-config/internal/domain/repository"
	"go.uber.org/zap"
)

type ConfigService struct{cfgRepo repository.SystemConfigRepository; flagRepo repository.FeatureFlagRepository; auditRepo repository.AuditLogRepository; logger *zap.Logger}
func NewConfigService(cfgRepo repository.SystemConfigRepository,flagRepo repository.FeatureFlagRepository,auditRepo repository.AuditLogRepository,logger *zap.Logger)*ConfigService{
	return &ConfigService{cfgRepo:cfgRepo,flagRepo:flagRepo,auditRepo:auditRepo,logger:logger}
}
func(s*ConfigService)GetConfigUseCase(ctx context.Context,category,key string)(string,error){cfg,err:=s.cfgRepo.GetByCategoryAndKey(ctx,category,key);if err!=nil{return "",err};if cfg==nil{return "",nil};return cfg.Value,nil}
func(s*ConfigService)IsFeatureEnabledUseCase(ctx context.Context,key string)(bool,error){return s.flagRepo.IsEnabled(ctx,key)}
func(s*ConfigService)RecordAuditUseCase(ctx context.Context,userID uuid.UUID,action,resource string,resourceID uuid.UUID,detail string)error{
log:=&entity.AuditLog{UserID:userID,Action:action,Resource:resource,ResourceID:resourceID,Detail:detail}
return s.auditRepo.Create(ctx,log)
}
