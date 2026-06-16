package application

import (
	"errors"

	"leap-one/service-project/internal/domain/repository"

	"go.uber.org/zap"
)

// 迭代服务相关错误定义
var (
	ErrIterationNotFound    = errors.New("迭代不存在")
	ErrInvalidIterationDate = errors.New("迭代日期无效（结束日期必须大于开始日期）")
	ErrInvalidStatusChange  = errors.New("无效的迭代状态转换")
)

// IterationService 迭代/Sprint应用服务
type IterationService struct {
	iterRepo    repository.IterationRepository
	projectRepo repository.ProjectRepository
	logger      *zap.Logger
}

// NewIterationService 创建迭代服务实例
func NewIterationService(
	iterRepo repository.IterationRepository,
	projectRepo repository.ProjectRepository,
	logger *zap.Logger,
) *IterationService {
	return &IterationService{
		iterRepo:    iterRepo,
		projectRepo: projectRepo,
		logger:      logger,
	}
}
