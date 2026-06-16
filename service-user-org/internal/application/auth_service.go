package application

import (
	"context"
	"errors"

	"leap-one/service-user-org/internal/domain/repository"

	"go.uber.org/zap"
)

// 认证服务相关错误定义
var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrTokenExpired       = errors.New("令牌已过期")
	ErrTokenRevoked       = errors.New("令牌已被撤销")
)

// AuthService 认证应用服务 - 协调认证相关业务流程
// 注意：Token生成不由这里做（由API网关负责），这里只做凭据验证
type AuthService struct {
	userRepo repository.UserRepository
	permRepo repository.PermissionRepository
	logger   *zap.Logger
}

// NewAuthService 创建认证应用服务实例
func NewAuthService(
	userRepo repository.UserRepository,
	permRepo repository.PermissionRepository,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		permRepo: permRepo,
		logger:   logger,
	}
}

// ValidateCredentials 验证用户凭据（用户名+密码）
// 用于网关调用本服务进行认证验证
func (s *AuthService) ValidateCredentials(ctx context.Context, username, password string) (userID string, roles []string, err error) {
	// 委托给UserService的LoginUseCase完成核心验证逻辑
	userService := NewUserService(s.userRepo, nil, s.logger)
	user, loginErr := userService.LoginUseCase(ctx, username, password, "")
	if loginErr != nil {
		return "", nil, loginErr
	}

	// 获取用户的权限编码作为角色信息
	roles, permErr := s.permRepo.GetUserPermissions(ctx, user.ID)
	if permErr != nil {
		s.logger.Warn("获取用户权限失败，返回空权限列表", zap.Error(permErr))
		roles = []string{}
	}

	return user.ID.String(), roles, nil
}

// CheckPermission 检查用户是否拥有指定权限
func (s *AuthService) CheckPermission(ctx context.Context, userIDStr, permissionCode string) (bool, error) {
	userID, parseErr := parseUUID(userIDStr)
	if parseErr != nil {
		return false, parseErr
	}

	hasPermission, err := s.permRepo.CheckPermission(ctx, userID, permissionCode)
	if err != nil {
		s.logger.Error("检查权限失败", zap.Error(err))
		return false, err
	}

	return hasPermission, nil
}

// GetUserPermissions 获取用户的所有权限编码
func (s *AuthService) GetUserPermissions(ctx context.Context, userIDStr string) ([]string, error) {
	userID, parseErr := parseUUID(userIDStr)
	if parseErr != nil {
		return nil, parseErr
	}

	codes, err := s.permRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		s.logger.Error("获取用户权限失败", zap.Error(err))
		return nil, err
	}

	return codes, nil
}
