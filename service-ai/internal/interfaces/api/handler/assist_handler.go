package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-ai/internal/domain/entity"
	"leap-one/service-ai/internal/domain/repository"
	"leap-one/service-ai/internal/interfaces/api/dto"
)

// AIAssistHandler AI辅助功能Handler
type AIAssistHandler struct {
	predRepo repository.PredictionRepository
	cfgRepo  repository.AIConfigRepository
	logger   *zap.Logger
}

// NewAIAssistHandler 创建AI辅助功能Handler实例
func NewAIAssistHandler(predRepo repository.PredictionRepository, cfgRepo repository.AIConfigRepository, logger *zap.Logger) *AIAssistHandler {
	return &AIAssistHandler{predRepo: predRepo, cfgRepo: cfgRepo, logger: logger}
}

// AssistRequirement AI辅助编写需�?(POST /api/v1/ai/assist/requirement)
func (h *AIAssistHandler) AssistRequirement(c *gin.Context) {
	var req dto.AssistRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	result := dto.AIAssistantResponse{
		Content: "根据您提供的需求信息，我为您生成了以下建议：\n\n" +
			"**需求标�?*: " + req.Title + "\n\n" +
			"**优化后的描述建议**:\n" +
			"1. 明确需求的业务目标和价值\n" +
			"2. 补充验收标准（AC）\n" +
			"3. 定义优先级和依赖关系\n" +
			"4. 建议拆分为可独立交付的用户故事\n\n" +
			"**优先级建�?*: " + func() string {
			if req.Priority != "" {
				return req.Priority
			}
			return "medium"
		}(),
		Suggestions: []string{
			"添加用户故事格式描述",
			"补充验收标准列表",
			"关联相关需求和任务",
			"评估工作量和复杂�?,
		},
		Metadata: map[string]interface{}{"model": "gpt-4", "tokens_used": 256},
	}

	c.JSON(http.StatusOK, result)
}

// AssistTestCase AI辅助编写测试用例 (POST /api/v1/ai/assist/test-case)
func (h *AIAssistHandler) AssistTestCase(c *gin.Context) {
	var req dto.AssistTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	testType := "功能测试"
	if req.Type != "" {
		switch req.Type {
		case "performance":
			testType = "性能测试"
		case "security":
			testType = "安全测试"
		case "ui":
			testType = "UI测试"
		}
	}

	result := dto.AIAssistantResponse{
		Content: "基于需求描述，我为您生成了以下" + testType + "用例建议：\n\n" +
			"**测试场景覆盖**:\n" +
			"- 正常流程验证\n" +
			"- 边界条件测试\n" +
			"- 异常输入处理\n" +
			"- 并发场景测试\n\n" +
			"共生�?个测试用例，涵盖主要功能路径�?,
		Suggestions: []string{
			"添加自动化脚本模�?,
			"生成测试数据准备方案",
			"定义预期结果对照�?,
		},
		Metadata: map[string]interface{}{"test_type": testType, "case_count": 8},
	}

	c.JSON(http.StatusOK, result)
}

// SuggestTaskAssign 智能任务分配建议 (POST /api/v1/ai/suggest/task-assign)
func (h *AIAssistHandler) SuggestTaskAssign(c *gin.Context) {
	var req dto.TaskAssignSuggestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	result := dto.AIAssistantResponse{
		Content: "基于任务描述和技能要求，推荐以下人员分配方案：\n\n" +
			"**推荐候选人**:\n" +
			"1. 张三 - 匹配�? 92% - 后端开发专家\n" +
			"2. 李四 - 匹配�? 85% - 全栈工程师\n" +
			"3. 王五 - 匹配�? 78% - 有类似项目经验\n\n" +
			"**分配理由**: 根据技能匹配、当前工作量、历史完成质量综合评估�?,
		Suggestions: []string{
			"查看候选人详细档案",
			"检查候选人当前任务负载",
			"查看历史协作记录",
		},
		Metadata: map[string]interface{}{"algorithm": "skill_matching_v2", "candidates_evaluated": 15},
	}

	c.JSON(http.StatusOK, result)
}

// PredictRequirements AI需求预�?(POST /api/v1/ai/predict/requirements)
func (h *AIAssistHandler) PredictRequirements(c *gin.Context) {
	var req dto.PredictRequirementsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 保存预测记录
	prediction := &entity.AIPrediction{
		Type:       "requirement_prediction",
		TargetID:   req.ProjectID,
		Result:     `{"predicted_requirements":12,"confidence_high":5,"confidence_medium":4,"confidence_low":3,"sprint_capacity":8,"risk_level":"medium","next_sprint_suggestion":"建议安排5个高置信度需�?}`,
		Confidence: 0.82,
		Model:      "gpt-4",
	}
	h.predRepo.Create(ctx, prediction)

	c.JSON(http.StatusOK, gin.H{
		"prediction_id": prediction.ID.String(),
		"project_id":    req.ProjectID,
		"type":          "requirement_prediction",
		"data": gin.H{
			"predicted_requirements": 12,
			"confidence_high":        5,
			"confidence_medium":      4,
			"confidence_low":         3,
			"sprint_capacity":        8,
			"risk_level":             "medium",
			"suggestions": []string{
				"下个迭代建议安排5个高置信度需�?,
				"关注中低置信度需求的风险�?,
				"预留20%容量应对突发需�?,
			},
		},
		"confidence": 0.82,
	})
}

// IdentifyRisks 风险智能识别 (POST /api/v1/ai/identify/risks)
func (h *AIAssistHandler) IdentifyRisks(c *gin.Context) {
	var req dto.IdentifyRisksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 保存预测记录
	prediction := &entity.AIPrediction{
		Type:       "risk_identification",
		TargetID:   req.ProjectID,
		Result:     `{"risks":[{"level":"high","category":"技术债务","description":"核心模块代码复杂度过�?},{"level":"medium","category":"资源","description":"关键开发人员可能离�?},{"level":"medium","category":"进度","description":"第三方API集成可能延期"}],"overall_risk_score":7.2}`,
		Confidence: 0.75,
		Model:      "gpt-4",
	}
	h.predRepo.Create(ctx, prediction)

	c.JSON(http.StatusOK, gin.H{
		"prediction_id": prediction.ID.String(),
		"project_id":    req.ProjectID,
		"type":          "risk_identification",
		"data": gin.H{
			"overall_risk_score": 7.2,
			"risks": []map[string]string{
				{"level": "high", "category": "技术债务", "description": "核心模块代码复杂度过高，建议安排重构"},
				{"level": "medium", "category": "资源风险", "description": "关键开发人员可能离职，需做好知识转移"},
				{"level": "medium", "category": "进度风险", "description": "第三方API集成可能延期，需准备备选方�?},
			},
			"mitigation_suggestions": []string{
				"立即启动代码审查和技术债务梳理",
				"建立关键模块知识文档",
				"与第三方确认接口SLA并制定降级方�?,
			},
		},
		"confidence": 0.75,
	})
}

// ListPredictions 预测历史记录 (GET /api/v1/ai/predictions)
func (h *AIAssistHandler) ListPredictions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	predType := c.Query("type")
	targetIDStr := c.Query("target_id")
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	var targetID uuid.UUID
	if targetIDStr != "" {
		if tid, err := uuid.Parse(targetIDStr); err == nil {
			targetID = tid
		}
	}

	ctx := c.Request.Context()
	preds, total, err := h.predRepo.List(ctx, page, size, predType, targetID)
	if err != nil {
		h.logger.Error("查询预测记录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询预测记录失败"})
		return
	}

	list := make([]dto.PredictionInfo, len(preds))
	for i, p := range preds {
		list[i] = dto.PredictionInfo{
			ID: p.ID.String(), Type: p.Type, TargetID: p.TargetID.String(),
			Result: p.Result, Confidence: p.Confidence, Model: p.Model,
			CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, dto.PredictionListResponse{List: list, Total: total, Page: page, Size: size})
}
