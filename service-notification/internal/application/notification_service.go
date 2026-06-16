package application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"leap-one/service-notification/internal/domain/entity"
	"leap-one/service-notification/internal/domain/repository"
	"go.uber.org/zap"
)

type NotificationService struct {
	notiRepo repository.NotificationRepository
	tplRepo  repository.TemplateRepository
	subRepo  repository.SubscriptionRepository
	logger   *zap.Logger
}

func NewNotificationService(notiRepo repository.NotificationRepository, tplRepo repository.TemplateRepository, subRepo repository.SubscriptionRepository, logger *zap.Logger) *NotificationService {
	return &NotificationService{notiRepo: notiRepo, tplRepo: tplRepo, subRepo: subRepo, logger: logger}
}

func (s *NotificationService) SendNotificationUseCase(ctx context.Context, templateCode string, receiverIDs []uuid.UUID, variables map[string]string, channel string) ([]*entity.Notification, error) {
	template, _ := s.tplRepo.GetByCode(ctx, templateCode)
	title := "系统通知"
	content := ""
	notifType := "system"
	if template != nil {
		title = renderTpl(template.Subject, variables)
		content = renderTpl(template.Body, variables)
		notifType = template.EventType
		if notifType == "" {
			notifType = "system"
		}
	} else {
		if t, ok := variables["title"]; ok {
			title = t
		}
		if c, ok := variables["content"]; ok {
			content = c
		}
	}
	now := time.Now()
	notifications := make([]*entity.Notification, len(receiverIDs))
	for i, rid := range receiverIDs {
		notifications[i] = &entity.Notification{ReceiverID: rid, Title: title, Content: content, Type: notifType, Channel: channel, SentAt: now}
	}
	if err := s.notiRepo.BatchCreate(ctx, notifications); err != nil {
		return nil, err
	}
	s.logger.Info("批量发送通知完成", zap.Int("count", len(receiverIDs)))
	return notifications, nil
}

func renderTpl(tpl string, vars map[string]string) string {
	result := tpl
	for k, v := range vars {
		result = strings.ReplaceAll(result, "{{"+k+"}}", v)
	}
	return result
}
