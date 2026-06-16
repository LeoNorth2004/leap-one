package application

import (
	"context"
	"errors"
	"time"

	"leap-one/service-user-org/internal/domain/entity"
	"leap-one/service-user-org/internal/domain/repository"
	"leap-one/service-user-org/internal/infrastructure/repository_impl"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 用户服务相关错误定义
var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrUserAlreadyExists = errors.New("用户名已存在")
	ErrInvalidPassword   = errors.New("密码无效")
	ErrUserDisabled      = errors.New("账号已被禁用")
)

// UserService 用户应用服务 - 协调用户相关的业务流程
type UserService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	logger   *zap.Logger
}

// NewUserService 创建用户应用服务实例
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		logger:   logger,
	}
}

// RegisterUseCase 注册用例
// 完整流程：校验唯一性 → 密码加密 → 创建用户 → 分配默认角色 → 返回结果
func (s *UserService) RegisterUseCase(ctx context.Context, username, password, email, realName string) (*entity.User, error) {
	// 1. 校验用户名唯一性
	existing, _ := s.userRepo.GetByUsername(ctx, username)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	// 2. 密码bcrypt加密（成本因子12）
	hashedPwd, err := repository_impl.HashPassword(password)
	if err != nil {
		s.logger.Error("注册时密码加密失败", zap.Error(err))
		return nil, errors.New("系统内部错误")
	}

	// 3. 创建用户实体
	user := &entity.User{
		Username: username,
		Password: hashedPwd,
		Email:    email,
		RealName: realName,
		Status:   1,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("注册时创建用户失败", zap.Error(err), zap.String("username", username))
		return nil, errors.New("创建用户失败")
	}

	// 4. 分配默认角色（如"viewer"查看者角色）
	defaultRole, _ := s.roleRepo.GetByCode(ctx, "viewer")
	if defaultRole != nil {
		// TODO: 通过UserRole关联表分配默认角色
		s.logger.Info("为新用户分配默认角色",
			zap.String("user_id", user.ID.String()),
			zap.String("role_code", defaultRole.Code),
		)
	}

	s.logger.Info("用户注册成功",
		zap.String("user_id", user.ID.String()),
		zap.String("username", username),
	)

	return user, nil
}

// LoginUseCase 登录用例
// 完整流程：查询用户 → 验证密码 → 检查状态 → 更新登录信息 → 返回用户信息
func (s *UserService) LoginUseCase(ctx context.Context, username, password, clientIP string) (*entity.User, error) {
	// 1. 查询用户
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	// 2. 验证密码
	if err := repository_impl.ValidatePassword(password, user.Password); err != nil {
		s.logger.Warn("登录密码验证失败",
			zap.String("username", username),
			zap.String("client_ip", clientIP),
		)
		return nil, ErrInvalidPassword
	}

	// 3. 检查状态
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	// 4. 更新最后登录时间和IP
	now := time.Now()
	if loginErr := s.userRepo.UpdateLastLogin(ctx, user.ID, clientIP); loginErr != nil {
		s.logger.Warn("更新最后登录时间失败", zap.Error(loginErr))
		user.LastLoginAt = &now // 内存中更新，即使数据库更新失败也不影响登录流程
		user.LastLoginIP = clientIP
	}

	return user, nil
}

// GetUserDetailUseCase 获取用户详情用例（含角色和组信息）
func (s *UserService) GetUserDetailUseCase(ctx context.Context, userID uuid.UUID) (*entity.User, []string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, nil, ErrUserNotFound
	}

	// TODO: 获取用户的角色列表和组列表
	roleCodes := []string{}

	return user, roleCodes, nil
}

// ChangePasswordUseCase 修改密码用例
func (s *UserService) ChangePasswordUseCase(ctx context.Context, userID uuid.UUID, oldPwd, newPwd string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	// 验证旧密码
	if err := repository_impl.ValidatePassword(oldPwd, user.Password); err != nil {
		return ErrInvalidPassword
	}

	// 加密新密码
	hashedNewPwd, err := repository_impl.HashPassword(newPwd)
	if err != nil {
		return errors.New("密码加密失败")
	}

	return s.userRepo.UpdatePassword(ctx, userID, hashedNewPwd)
}
