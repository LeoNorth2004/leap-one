package repository_impl

import (
	"context"
	"fmt"
	"time"

	"leap-one/service-user-org/internal/domain/entity"
	"leap-one/service-user-org/internal/domain/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepositoryImpl 用户仓库实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create 创建用户（密码应在调用前已bcrypt加密）
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据ID获取用户（带缓存）
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户（用于登录验证）
func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 软删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

// List 分页查询用户列表（支持关键词搜索）
func (r *UserRepositoryImpl) List(ctx context.Context, page, pageSize int, keyword string) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.User{})

	// 关键词搜索：匹配用户名、邮箱、真实姓名、手机号
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where(
			"username LIKE ? OR email LIKE ? OR real_name LIKE ? OR phone LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdatePassword 更新用户密码
func (r *UserRepositoryImpl) UpdatePassword(ctx context.Context, id uuid.UUID, hashedPassword string) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error
}

// UpdateLastLogin 更新最后登录时间和IP地址
func (r *UserRepositoryImpl) UpdateLastLogin(ctx context.Context, id uuid.UUID, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at": now,
			"last_login_ip": ip,
		}).Error
}

// HashPassword 使用bcrypt加密密码（成本因子12）
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}
	return string(bytes), nil
}

// ValidatePassword 验证密码是否正确
func ValidatePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
