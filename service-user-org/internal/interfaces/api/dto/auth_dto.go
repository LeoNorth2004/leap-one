package dto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名
	Password string `json:"password" binding:"required,min=6,max=50"` // 密码
}

// LoginResponse 登录响应（含JWT令牌）
type LoginResponse struct {
	Token     string    `json:"token"`      // JWT访问令牌
	ExpiresAt time.Time `json:"expires_at"` // 令牌过期时间
	UserInfo  UserInfo  `json:"user_info"`  // 用户基本信息
}

// UserInfo 用户基本信息（不含敏感信息）
type UserInfo struct {
	ID           string   `json:"id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Phone        string   `json:"phone"`
	RealName     string   `json:"real_name"`
	Avatar       string   `json:"avatar"`
	Status       int8     `json:"status"`
	DepartmentID string   `json:"department_id,omitempty"`
	Roles        []string `json:"roles,omitempty"` // 角色编码列表
}

// JWTClaims JWT令牌声明结构
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Roles    string `json:"roles"` // 逗号分隔的角色编码列表
	jwt.RegisteredClaims
}
