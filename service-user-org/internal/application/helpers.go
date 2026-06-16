package application

import (
	"github.com/google/uuid"
)

// parseUUID 安全解析UUID字符串
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
