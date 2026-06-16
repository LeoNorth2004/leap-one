package application

import (
	"errors"

	"leap-one/service-project/internal/domain/repository"

	"go.uber.org/zap"
)

// 成员服务相关错误定义
var (
	ErrMemberNotFound      = errors.New("成员不存在")
	ErrMemberAlreadyExists = errors.New("该用户已是项目成员")
	ErrInvalidMemberRole   = errors.New("无效的成员角色")
)

// ProjectMemberService 项目成员应用服务
type ProjectMemberService struct {
	memberRepo  repository.ProjectMemberRepository
	projectRepo repository.ProjectRepository
	logger      *zap.Logger
}

// NewProjectMemberService 创建成员服务实例
func NewProjectMemberService(
	memberRepo repository.ProjectMemberRepository,
	projectRepo repository.ProjectRepository,
	logger *zap.Logger,
) *ProjectMemberService {
	return &ProjectMemberService{
		memberRepo:  memberRepo,
		projectRepo: projectRepo,
		logger:      logger,
	}
}
