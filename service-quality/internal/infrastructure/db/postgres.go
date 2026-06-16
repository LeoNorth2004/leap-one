package db

import (
	"context"
	"fmt"
	"time"

	"leap-one/service-quality/internal/config"
	"leap-one/service-quality/internal/domain/entity"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// InitPostgreSQL 初始化PostgreSQL数据库连�?// 使用GORM连接PostgreSQL，配置连接池和日�?func InitPostgreSQL(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	// 配置GORM自定义Logger（集成Zap�?	gormLog := &GormZapLogger{
		logger:        logger,
		slowThreshold: 200 * time.Millisecond, // 慢查询阈值：200ms
		logLevel:      gormlogger.Info,
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: gormLog,
	})
	if err != nil {
		return nil, fmt.Errorf("无法连接数据�? %w", err)
	}

	// 获取底层SQL DB以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失�? %w", err)
	}

	// 连接池配置：最大打开连接�?100，最大空闲连接数=10，连接最大存活时�?30分钟
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	logger.Info("PostgreSQL连接初始化成�?,
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("dbname", cfg.Database.DBName),
	)

	return db, nil
}

// AutoMigrate 自动迁移所有实体表到数据库
// 在服务启动时调用，确保表结构与实体定义一�?func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		// 核心业务实体�?		&entity.TestCase{},
		&entity.TestSuite{},
		&entity.TestCaseSuiteRelation{},
		&entity.TestPlan{},
		&entity.TestPlanCase{},
		&entity.Bug{},
		&entity.BugComment{},
		&entity.BugAttachment{},
		&entity.BugHistory{},
		&entity.BugWorkflow{},
		&entity.BugWorkflowTransition{},
		&entity.TestEnvironment{},
	)
	if err != nil {
		return fmt.Errorf("数据库自动迁移失�? %w", err)
	}
	return nil
}

// GormZapLogger GORM日志适配�?- 将GORM日志输出到Zap
type GormZapLogger struct {
	logger        *zap.Logger
	slowThreshold time.Duration
	logLevel      gormlogger.LogLevel
}

// LogMode 设置日志级别（实现gormlogger.Interface接口�?func (l *GormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info 记录信息级别日志
func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Info {
		l.logger.Sugar().Infof(msg, data...)
	}
}

// Warn 记录警告级别日志
func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Warn {
		l.logger.Sugar().Warnf(msg, data...)
	}
}

// Error 记录错误级别日志
func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Error {
		l.logger.Sugar().Errorf(msg, data...)
	}
}

// Trace 记录SQL执行日志（含慢查询标记）
func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}

	switch {
	case err != nil && l.logLevel >= gormlogger.Error:
		fields = append(fields, zap.Error(err))
		l.logger.Error("SQL错误", fields...)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= gormlogger.Warn:
		l.logger.Warn("慢查询检�?, fields...)
	case l.logLevel == gormlogger.Info:
		l.logger.Debug("SQL执行", fields...)
	}
}
