package repository_impl

import (
	"errors"
	"time"
)

// 通用错误定义
var (
	ErrSystemTemplateCannotDelete = errors.New("系统预置模板不允许删除")
)

// NowFunc 当前时间函数（方便测试时替换）
var NowFunc = func() time.Time {
	return time.Now()
}

// buildOrderClause 构建排序子句
func buildOrderClause(sortBy, sortOrder, defaultSort, defaultOrder string) string {
	// 允许的排序列白名单
	allowedSortFields := map[string]string{
		"created_at":   "created_at",
		"updated_at":   "updated_at",
		"name":        "name",
		"code":        "code",
		"priority":    "priority",
		"status":      "status",
		"start_date":  "start_date",
		"end_date":    "end_date",
		"due_date":    "due_date",
		"sort_order":  "sort_order",
		"severity":    "severity",
	}

	field, ok := allowedSortFields[sortBy]
	if !ok || field == "" {
		field = defaultSort
	}

	order := "ASC"
	if sortOrder == "desc" || sortOrder == "DESC" {
		order = "DESC"
	} else if defaultOrder == "DESC" {
		order = "DESC"
	}

	return field + " " + order
}
