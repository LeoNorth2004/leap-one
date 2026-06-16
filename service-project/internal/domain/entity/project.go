package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project 项目实体 - 项目管理核心模型
type Project struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(200);not null" json:"name"`                    // 项目名称
	Code        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`        // 项目编号（唯一）
	ProgramID   *uuid.UUID     `gorm:"type:uuid;index" json:"program_id"`                        // 关联项目集ID
	Description string         `gorm:"type:text" json:"description"`                             // 项目描述
	PMID        uuid.UUID      `gorm:"type:uuid;not null" json:"pm_id"`                          // 项目经理ID
	Status      string         `gorm:"type:varchar(20);default:'planning'" json:"status"`       // 状态：planning/executing/paused/completed/cancelled/archived
	Type        string         `gorm:"type:varchar(30);default:'agile'" json:"type"`            // 类型：agile/waterfall/lightweight/lifecycle
	Priority    int            `gorm:"default:3" json:"priority"`                                // 优先级 1-5
	StartDate   *time.Time     `json:"start_date"`                                               // 开始日期
	EndDate     *time.Time     `json:"end_date"`                                                 // 结束日期
	Budget      *float64       `json:"budget"`                                                   // 预算
	TemplateID  *uuid.UUID     `gorm:"type:uuid" json:"template_id"`                             // 关联模板ID
	CreatedByID uuid.UUID      `gorm:"type:uuid" json:"created_by_id"`                           // 创建人ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedByID *uuid.UUID     `gorm:"type:uuid" json:"updated_by_id"`                           // 更新人ID
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Version     int            `gorm:"default:1" json:"version"`                                 // 乐观锁版本号
}

// TableName 指定数据库表名
func (Project) TableName() string {
	return "projects"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// ProjectMember 项目成员实体
type ProjectMember struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID uuid.UUID      `gorm:"type:uuid;index;not null" json:"project_id"` // 所属项目ID
	UserID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`    // 用户ID
	Role      string         `gorm:"type:varchar(30)" json:"role"`               // 角色：pm/po/dev/qa/viewer
	JoinTime  time.Time      `json:"join_time"`                                   // 加入时间
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (ProjectMember) TableName() string {
	return "project_members"
}

// ProjectTemplate 项目模板实体
type ProjectTemplate struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(200);not null" json:"name"`           // 模板名称
	Description string         `gorm:"type:text" json:"description"`                      // 模板描述
	Type        string         `gorm:"type:varchar(30)" json:"type"`                     // 模板类型：agile/waterfall
	Config      string         `gorm:"type:text" json:"config"`                           // JSON配置（预设字段、阶段等）
	IsSystem    bool           `gorm:"default:false" json:"is_system"`                   // 是否系统预置
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (ProjectTemplate) TableName() string {
	return "project_templates"
}

// ProjectMilestone 项目里程碑实体
type ProjectMilestone struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"project_id"`   // 所属项目ID
	Name        string         `gorm:"type:varchar(200);not null" json:"name"`       // 里程碑名称
	Description string         `gorm:"type:text" json:"description"`                 // 描述
	DueDate     time.Time      `json:"due_date"`                                     // 截止日期
	Status      string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // 状态：pending/completed/overdue
	CompletedAt *time.Time     `json:"completed_at"`                                  // 完成时间
	CompletedBy *uuid.UUID     `gorm:"type:uuid" json:"completed_by"`                // 完成人ID
	SortOrder   int            `gorm:"default:0" json:"sort_order"`                   // 排序序号
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (ProjectMilestone) TableName() string {
	return "project_milestones"
}

// ProjectRisk 项目风险实体
type ProjectRisk struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"project_id"`   // 所属项目ID
	Title       string         `gorm:"type:varchar(300);not null" json:"title"`      // 风险标题
	Description string         `gorm:"type:text" json:"description"`                  // 风险描述
	Probability int            `gorm:"default:3" json:"probability"`                   // 发生概率 1-5
	Impact      int            `gorm:"default:3" json:"impact"`                       // 影响程度 1-5
	Severity    int            `json:"severity"`                                      // 严重程度 = probability * impact
	Status      string         `gorm:"type:varchar(20);default:'open'" json:"status"` // 状态：open/mitigating/closed
	OwnerID     uuid.UUID      `gorm:"type:uuid" json:"owner_id"`                     // 风险负责人ID
	Mitigation  string         `gorm:"type:text" json:"mitigation"`                   // 缓解措施
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (ProjectRisk) TableName() string {
	return "project_risks"
}

// BeforeCreate 创建前钩子：自动计算严重程度
func (r *ProjectRisk) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	r.Severity = r.Probability * r.Impact
	return nil
}

// BeforeUpdate 更新前钩子：重新计算严重程度
func (r *ProjectRisk) BeforeUpdate(tx *gorm.DB) error {
	r.Severity = r.Probability * r.Impact
	return nil
}

// CustomField 自定义字段实体
type CustomField struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID uuid.UUID      `gorm:"type:uuid;index;not null" json:"project_id"` // 所属项目ID
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`     // 字段显示名称
	FieldKey  string         `gorm:"type:varchar(50);not null" json:"field_key"` // 字段标识
	FieldType string         `gorm:"type:varchar(20)" json:"field_type"`         // 字段类型：text/number/date/select/user
	Options   string         `gorm:"type:text" json:"options"`                   // JSON选项（select类型用）
	Required  bool           `gorm:"default:false" json:"required"`              // 是否必填
	SortOrder int            `gorm:"default:0" json:"sort_order"`               // 排序序号
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (CustomField) TableName() string {
	return "custom_fields"
}
