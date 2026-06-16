package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 网关服务全局配置
type Config struct {
	Server    ServerConfig             `mapstructure:"server"`
	JWT       JWTConfig                `mapstructure:"jwt"`
	RateLimit RateLimitConfig          `mapstructure:"rate_limit"`
	Services  map[string]ServiceConfig `mapstructure:"services"`
	CORS      CORSConfig               `mapstructure:"cors"`
}

// ServerConfig HTTP服务器配置
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// JWTConfig JWT认证配置
type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	AccessExpire  time.Duration `mapstructure:"access_expire"`
	RefreshExpire time.Duration `mapstructure:"refresh_expire"`
	Issuer        string        `mapstructure:"issuer"`
	SkipPaths     []string      `mapstructure:"skip_paths"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool    `mapstructure:"enabled"`
	RPS               float64 `mapstructure:"rps"`
	Burst             int     `mapstructure:"burst"`
	RequestsPerMinute int     `mapstructure:"requests_per_minute"`
}

// ServiceConfig 下游微服务地址配置
type ServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// CORSConfig 跨域配置
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	MaxAge           int      `mapstructure:"max_age"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

// GetServiceURL 根据服务名获取完整的服务地址
func (c *Config) GetServiceURL(serviceName string) string {
	svc, ok := c.Services[serviceName]
	if !ok {
		return ""
	}
	return fmt.Sprintf("http://%s:%d", svc.Host, svc.Port)
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if configPath != "" {
		v.AddConfigPath(configPath)
	}
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	v.SetEnvPrefix("LEAP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("警告: 未找到配置文件，将使用默认值和环境变量")
		} else {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)

	v.SetDefault("jwt.secret", "leap-one-default-secret-change-in-production")
	v.SetDefault("jwt.access_expire", "2h")
	v.SetDefault("jwt.refresh_expire", 168*time.Hour)
	v.SetDefault("jwt.issuer", "leap-one-gateway")
	v.SetDefault("jwt.skip_paths", []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/refresh-token",
		"/healthz",
	})

	v.SetDefault("rate_limit.enabled", true)
	v.SetDefault("rate_limit.rps", 100)
	v.SetDefault("rate_limit.burst", 200)
	v.SetDefault("rate_limit.requests_per_minute", 100)

	v.SetDefault("cors.allow_origins", []string{"*"})
	v.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"})
	v.SetDefault("cors.allow_headers", []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "X-User-ID"})
	v.SetDefault("cors.expose_headers", []string{"Content-Length", "Content-Type"})
	v.SetDefault("cors.max_age", 86400)
	v.SetDefault("cors.allow_credentials", true)

	serviceNames := []string{
		"user-org", "portfolio", "project", "task",
		"requirement", "quality", "devops", "document",
		"kanban", "bi", "ai", "notification", "search", "config",
	}
	basePort := 8001
	for i, name := range serviceNames {
		v.SetDefault(fmt.Sprintf("services.%s.host", name), "localhost")
		v.SetDefault(fmt.Sprintf("services.%s.port", name), basePort+i)
	}
}
