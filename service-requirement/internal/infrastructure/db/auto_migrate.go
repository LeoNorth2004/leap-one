package db

import (
	"leap-one/service-requirement/internal/domain/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB, logger *zap.Logger) error {
	err := db.AutoMigrate(
		&entity.Requirement{},
		&entity.RequirementRelation{},
		&entity.RequirementReview{},
		&entity.RequirementReviewParticipant{},
		&entity.RequirementChangeLog{},
	)
	if err != nil {
		logger.Error("数据库自动迁移失�?, zap.Error(err))
		return err
	}
	logger.Info("数据库表结构迁移完成")
	return nil
}
