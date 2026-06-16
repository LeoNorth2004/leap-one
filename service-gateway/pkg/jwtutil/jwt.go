package jwtutil

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"leap-one/service-gateway/internal/config"
)

// TokenClaims JWT声明结构
type TokenClaims struct {
	UserID   string   `json:"userId"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// TokenPair Access Token和Refresh Token对
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // Access Token过期时间（秒）
}

// Generator JWT生成器
type Generator struct {
	cfg *config.JWTConfig
}

// NewGenerator 创建JWT生成器实例
func NewGenerator(cfg *config.JWTConfig) *Generator {
	return &Generator{cfg: cfg}
}

// GenerateToken 生成Access Token和Refresh Token对
func (g *Generator) GenerateToken(userID, username string, roles []string) (*TokenPair, error) {
	now := time.Now()

	// 生成Access Token
	accessClaims := TokenClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.cfg.Issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(g.cfg.AccessExpire)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSigned, err := accessToken.SignedString([]byte(g.cfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("签发AccessToken失败: %w", err)
	}

	// RefreshToken（仅包含基本信息，用于续期）
	refreshClaims := TokenClaims{
		UserID:   userID,
		Roles:    []string{}, // RefreshToken不携带角色信息
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.cfg.Issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(g.cfg.RefreshExpire)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshSigned, err := refreshToken.SignedString([]byte(g.cfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("签发RefreshToken失败: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessSigned,
		RefreshToken: refreshSigned,
		ExpiresIn:    int64(g.cfg.AccessExpire.Seconds()),
	}, nil
}

// ParseToken 解析并验证Token，返回Claims
func (g *Generator) ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", t.Header["alg"])
		}
		return []byte(g.cfg.Secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.Join(err, errors.New("token_expired"))
		}
		return nil, fmt.Errorf("解析Token失败: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的Token Claims")
	}

	return claims, nil
}

// RefreshToken 使用RefreshToken刷新Token对
func (g *Generator) RefreshToken(refreshTokenStr string) (*TokenPair, error) {
	claims, err := g.ParseToken(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("RefreshToken无效: %w", err)
	}
	return g.GenerateToken(claims.UserID, claims.Username, nil)
}
