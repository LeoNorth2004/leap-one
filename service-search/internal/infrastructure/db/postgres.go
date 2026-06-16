package db

import (
	"context"
	"fmt"
	"time"

	"leap-one/service-search/internal/config"
	"leap-one/service-search/internal/domain/entity"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func InitPostgreSQL(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	gormLog := &GormZapLogger{logger: logger, slowThreshold: 200 * time.Millisecond, logLevel: gormlogger.Info}
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{Logger: gormLog})
	if err != nil {
		return nil, fmt.Errorf("无法连接数据�?%w", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	logger.Info("PostgreSQL连接初始化成�?, zap.String("host", cfg.Database.Host), zap.Int("port", cfg.Database.Port), zap.String("dbname", cfg.Database.DBName))
	return db, nil
}
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.SearchDocument{}, &entity.SavedSearch{}, &entity.SearchHistory{}); err != nil {
		return fmt.Errorf("数据库自动迁移失�?%w", err)
	}
	return nil
}

type GormZapLogger struct {
	logger        *zap.Logger
	slowThreshold time.Duration
	logLevel      gormlogger.LogLevel
}

func (l *GormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	nl := *l
	nl.logLevel = level
	return &nl
}
func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Info {
		l.logger.Sugar().Infof(msg, data...)
	}
}
func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Warn {
		l.logger.Sugar().Warnf(msg, data...)
	}
}
func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Error {
		l.logger.Sugar().Errorf(msg, data...)
	}
}
func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql)}
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
