package dto

// CreateProgramRequest 创建项目集请求
type CreateProgramRequest struct {
	Name        string   `json:"name" binding:"required,max=200"` // 项目集名称
	Code        string   `json:"code" binding:"required,max=50"`  // 编号
	Description string   `json:"description"`                     // 描述
	ParentID    *string  `json:"parent_id"`                       // 父项目集ID
	OwnerID     string   `json:"owner_id" binding:"required"`     // 负责人ID
	Budget      *float64 `json:"budget"`                          // 预算
	StartDate   *string  `json:"start_date"`                      // 开始日期（YYYY-MM-DD）
	EndDate     *string  `json:"end_date"`                        // 结束日期（YYYY-MM-DD）
	Priority    int      `json:"priority" binding:"min=1,max=5"`  // 优先级 1-5
}

// UpdateProgramRequest 更新项目集请求
type UpdateProgramRequest struct {
	Name        *string  `json:"name" binding:"omitempty,max=200"`
	Code        *string  `json:"code" binding:"omitempty,max=50"`
	Description *string  `json:"description"`
	ParentID    *string  `json:"parent_id"`
	OwnerID     *string  `json:"owner_id"`
	Status      *string  `json:"status" binding:"omitempty,oneof=active paused completed cancelled"`
	Budget      *float64 `json:"budget"`
	StartDate   *string  `json:"start_date"`
	EndDate     *string  `json:"end_date"`
	Priority    *int     `json:"priority" binding:"omitempty,min=1,max=5"`
}

// ProgramInfo 项目集基本信息
type ProgramInfo struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Code        string        `json:"code"`
	Description string        `json:"description"`
	ParentID    *string       `json:"parent_id"`
	OwnerID     string        `json:"owner_id"`
	Status      string        `json:"status"`
	Budget      *float64      `json:"budget"`
	StartDate   *string       `json:"start_date"`
	EndDate     *string       `json:"end_date"`
	Priority    int           `json:"priority"`
	Children    []ProgramInfo `json:"children,omitempty"`
}

// ProgramListResponse 项目集列表响应（分页）
type ProgramListResponse struct {
	List  []ProgramInfo `json:"list"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// ProgramDetailResponse 项目集详情响应
type ProgramDetailResponse struct {
	ProgramInfo
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ProgramStatisticsResponse 项目集统计信息响应
type ProgramStatisticsResponse struct {
	ProgramID       string `json:"program_id"`
	TotalProducts   int64  `json:"total_products"`   // 关联产品数
	ActiveProducts  int64  `json:"active_products"`  // 活跃产品数
	TotalMilestones int64  `json:"total_milestones"` // 里程碑总数
	DoneMilestones  int64  `json:"done_milestones"`  // 已完成里程碑数
	TotalRisks      int64  `json:"total_risks"`      // 风险总数
	OpenRisks       int64  `json:"open_risks"`       // 未关闭风险数
	ChildCount      int64  `json:"child_count"`      // 子项目集数量
}

// CreateMilestoneRequest 创建里程碑请求
type CreateMilestoneRequest struct {
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"` // YYYY-MM-DD
	Status      string `json:"status" binding:"omitempty,oneof=pending completed overdue"`
}

// MilestoneInfo 里程碑信息
type MilestoneInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	DueDate     *string `json:"due_date"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

// CreateRiskRequest 创建风险请求
type CreateRiskRequest struct {
	Title       string  `json:"title" binding:"required,max=300"`
	Description string  `json:"description"`
	Probability string  `json:"probability" binding:"omitempty,oneof=low medium high"`
	Impact      string  `json:"impact" binding:"omitempty,oneof=low medium high"`
	Status      string  `json:"status" binding:"omitempty,oneof=open mitigating closed"`
	OwnerID     *string `json:"owner_id"`
}

// RiskInfo 风险项信息
type RiskInfo struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Probability string  `json:"probability"`
	Impact      string  `json:"impact"`
	Status      string  `json:"status"`
	OwnerID     *string `json:"owner_id"`
	CreatedAt   string  `json:"created_at"`
}
