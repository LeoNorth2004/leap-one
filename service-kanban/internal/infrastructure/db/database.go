package db

import (
	"fmt"

	"leap-one/service-kanban/internal/config"
	"leap-one/service-kanban/internal/domain/entity"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase 创建数据库连�?func NewDatabase(cfg *config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormConfig)
	if err != nil { return nil, fmt.Errorf("无法连接数据�? %w", err) }
	sqlDB, err := db.DB()
	if err != nil { return nil, fmt.Errorf("获取底层数据库连接失�? %w", err) }
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	log.Info("数据库连接成�?, zap.String("host", cfg.Host), zap.Int("port", cfg.Port), zap.String("dbname", cfg.DBName))
	return db, nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB, logger *zap.Logger) error {
	err := db.AutoMigrate(
		&entity.KanbanBoard{},
		&entity.KanbanColumn{},
		&entity.KanbanCard{},
		&entity.KanbanSwimlane{},
		&entity.KanbanCardMoveHistory{},
	)
	if err != nil { logger.Error("数据库自动迁移失�?, zap.Error(err)); return err }
	logger.Info("数据库表结构迁移完成")
	return nil
}
