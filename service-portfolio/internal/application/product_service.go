package application

import (
	"context"
	"errors"
	"time"

	"leap-one/service-portfolio/internal/domain/entity"
	"leap-one/service-portfolio/internal/domain/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ProductService 产品应用服务 - 协调产品相关的业务流程
type ProductService struct {
	productRepo     repository.ProductRepository
	productLineRepo repository.ProductLineRepository
	versionRepo     repository.ProductVersionRepository
	roadmapRepo     repository.ProductRoadmapRepository
	planRepo        repository.ProductPlanRepository
	programRepo     repository.ProgramRepository
	logger          *zap.Logger
}

// NewProductService 创建产品应用服务实例
func NewProductService(
	productRepo repository.ProductRepository,
	productLineRepo repository.ProductLineRepository,
	versionRepo repository.ProductVersionRepository,
	roadmapRepo repository.ProductRoadmapRepository,
	planRepo repository.ProductPlanRepository,
	programRepo repository.ProgramRepository,
	logger *zap.Logger,
) *ProductService {
	return &ProductService{
		productRepo:     productRepo,
		productLineRepo: productLineRepo,
		versionRepo:     versionRepo,
		roadmapRepo:     roadmapRepo,
		planRepo:        planRepo,
		programRepo:     programRepo,
		logger:          logger,
	}
}

// CreateProduct 创建产品用例
func (s *ProductService) CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	// 校验编码唯一性
	if existing, _ := s.productRepo.GetByCode(ctx, product.Code); existing != nil {
		return nil, ErrProductCodeExists
	}

	// 如果指定了产品线，校验其存在性
	if product.ProductLineID != nil && *product.ProductLineID != uuid.Nil {
		if line, _ := s.productLineRepo.GetByID(ctx, *product.ProductLineID); line == nil {
			return nil, ErrProductLineNotFound
		}
	}

	// 如果指定了项目集，校验其存在性
	if product.ProgramID != nil && *product.ProgramID != uuid.Nil {
		if prog, _ := s.programRepo.GetByID(ctx, *product.ProgramID); prog == nil {
			return nil, ErrProgramNotFound
		}
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		s.logger.Error("创建产品失败", zap.Error(err), zap.String("code", product.Code))
		return nil, errors.New("创建产品失败")
	}

	s.logger.Info("产品创建成功",
		zap.String("product_id", product.ID.String()),
		zap.String("code", product.Code),
	)
	return product, nil
}

// GetProductDetail 获取产品详情
func (s *ProductService) GetProductDetail(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil || product == nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

// UpdateProduct 更新产品信息
func (s *ProductService) UpdateProduct(ctx context.Context, product *entity.Product) error {
	existing, err := s.productRepo.GetByID(ctx, product.ID)
	if err != nil || existing == nil {
		return ErrProductNotFound
	}

	// 如果修改了编码，检查唯一性
	if product.Code != existing.Code {
		if dup, _ := s.productRepo.GetByCode(ctx, product.Code); dup != nil && dup.ID != product.ID {
			return ErrProductCodeExists
		}
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		s.logger.Error("更新产品失败", zap.Error(err))
		return errors.New("更新产品失败")
	}
	return nil
}

// DeleteProduct 删除产品（软删除）
func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil || product == nil {
		return ErrProductNotFound
	}

	if err := s.productRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除产品失败", zap.Error(err))
		return errors.New("删除产品失败")
	}
	return nil
}

// ListProducts 分页查询产品列表
func (s *ProductService) ListProducts(ctx context.Context, page, pageSize int, keyword, status string, productLineID *uuid.UUID) ([]*entity.Product, int64, error) {
	return s.productRepo.List(ctx, page, pageSize, keyword, status, productLineID)
}

// ==================== 产品线管理 ====================

// CreateProductLine 创建产品线
func (s *ProductService) CreateProductLine(ctx context.Context, line *entity.ProductLine) (*entity.ProductLine, error) {
	if err := s.productLineRepo.Create(ctx, line); err != nil {
		s.logger.Error("创建产品线失败", zap.Error(err))
		return nil, errors.New("创建产品线失败")
	}
	return line, nil
}

// UpdateProductLine 更新产品线
func (s *ProductService) UpdateProductLine(ctx context.Context, line *entity.ProductLine) error {
	existing, err := s.productLineRepo.GetByID(ctx, line.ID)
	if err != nil || existing == nil {
		return ErrProductLineNotFound
	}
	return s.productLineRepo.Update(ctx, line)
}

// DeleteProductLine 删除产品线
func (s *ProductService) DeleteProductLine(ctx context.Context, id uuid.UUID) error {
	existing, err := s.productLineRepo.GetByID(ctx, id)
	if err != nil || existing == nil {
		return ErrProductLineNotFound
	}
	return s.productLineRepo.Delete(ctx, id)
}

// ListProductLines 获取全部产品线列表
func (s *ProductService) ListProductLines(ctx context.Context) ([]*entity.ProductLine, error) {
	return s.productLineRepo.List(ctx)
}

// GetProductLineDetail 获取产品线详情
func (s *ProductService) GetProductLineDetail(ctx context.Context, id uuid.UUID) (*entity.ProductLine, error) {
	line, err := s.productLineRepo.GetByID(ctx, id)
	if err != nil || line == nil {
		return nil, ErrProductLineNotFound
	}
	return line, nil
}

// ==================== 版本管理 ====================

// CreateVersion 创建产品版本
func (s *ProductService) CreateVersion(ctx context.Context, version *entity.ProductVersion) error {
	// 检查产品是否存在
	if _, err := s.productRepo.GetByID(ctx, version.ProductID); err != nil {
		return ErrProductNotFound
	}

	if err := s.versionRepo.Create(ctx, version); err != nil {
		s.logger.Error("创建版本失败", zap.Error(err))
		return errors.New("创建版本失败")
	}
	return nil
}

// ListVersions 获取产品的版本列表
func (s *ProductService) ListVersions(ctx context.Context, productID uuid.UUID) ([]*entity.ProductVersion, error) {
	// 检查产品是否存在
	if _, err := s.productRepo.GetByID(ctx, productID); err != nil {
		return nil, ErrProductNotFound
	}
	return s.versionRepo.ListByProductID(ctx, productID)
}

// GetVersionDetail 获取版本详情
func (s *ProductService) GetVersionDetail(ctx context.Context, id uuid.UUID) (*entity.ProductVersion, error) {
	version, err := s.versionRepo.GetByID(ctx, id)
	if err != nil || version == nil {
		return nil, ErrVersionNotFound
	}
	return version, nil
}

// UpdateVersion 更新版本信息
func (s *ProductService) UpdateVersion(ctx context.Context, version *entity.ProductVersion) error {
	existing, err := s.versionRepo.GetByID(ctx, version.ID)
	if err != nil || existing == nil {
		return ErrVersionNotFound
	}
	return s.versionRepo.Update(ctx, version)
}

// ReleaseVersion 发布版本（将状态更新为released）
func (s *ProductService) ReleaseVersion(ctx context.Context, id uuid.UUID, releaseDate *time.Time) error {
	version, err := s.versionRepo.GetByID(ctx, id)
	if err != nil || version == nil {
		return ErrVersionNotFound
	}
	version.Status = "released"
	version.ReleaseDate = releaseDate
	return s.versionRepo.Update(ctx, version)
}

// ==================== 路线图管理 ====================

// CreateRoadmapItem 创建路线图项
func (s *ProductService) CreateRoadmapItem(ctx context.Context, item *entity.ProductRoadmapItem) error {
	// 检查产品是否存在
	if _, err := s.productRepo.GetByID(ctx, item.ProductID); err != nil {
		return ErrProductNotFound
	}

	if err := s.roadmapRepo.Create(ctx, item); err != nil {
		s.logger.Error("创建路线图项失败", zap.Error(err))
		return errors.New("创建路线图项失败")
	}
	return nil
}

// ListRoadmapItems 获取产品的路线图列表
func (s *ProductService) ListRoadmapItems(ctx context.Context, productID uuid.UUID) ([]*entity.ProductRoadmapItem, error) {
	// 检查产品是否存在
	if _, err := s.productRepo.GetByID(ctx, productID); err != nil {
		return nil, ErrProductNotFound
	}
	return s.roadmapRepo.ListByProductID(ctx, productID)
}

// GetRoadmapItemDetail 获取路线图项详情
func (s *ProductService) GetRoadmapItemDetail(ctx context.Context, id uuid.UUID) (*entity.ProductRoadmapItem, error) {
	item, err := s.roadmapRepo.GetByID(ctx, id)
	if err != nil || item == nil {
		return nil, ErrRoadmapItemNotFound
	}
	return item, nil
}

// UpdateRoadmapItem 更新路线图项
func (s *ProductService) UpdateRoadmapItem(ctx context.Context, item *entity.ProductRoadmapItem) error {
	existing, err := s.roadmapRepo.GetByID(ctx, item.ID)
	if err != nil || existing == nil {
		return ErrRoadmapItemNotFound
	}
	return s.roadmapRepo.Update(ctx, item)
}

// DeleteRoadmapItem 删除路线图项
func (s *ProductService) DeleteRoadmapItem(ctx context.Context, id uuid.UUID) error {
	existing, err := s.roadmapRepo.GetByID(ctx, id)
	if err != nil || existing == nil {
		return ErrRoadmapItemNotFound
	}
	return s.roadmapRepo.Delete(ctx, id)
}

// ReorderRoadmapItems 重新排序路线图项
func (s *ProductService) ReorderRoadmapItems(ctx context.Context, productID uuid.UUID, itemIDs []uuid.UUID) error {
	// 检查产品是否存在
	if _, err := s.productRepo.GetByID(ctx, productID); err != nil {
		return ErrProductNotFound
	}
	return s.roadmapRepo.UpdateSortOrder(ctx, productID, itemIDs)
}

// ==================== 计划管理 ====================

// CreatePlan 创建产品计划
func (s *ProductService) CreatePlan(ctx context.Context, plan *entity.ProductPlan) error {
	// 检查产品是否存在
	if _, err := s.productRepo.GetByID(ctx, plan.ProductID); err != nil {
		return ErrProductNotFound
	}

	if err := s.planRepo.Create(ctx, plan); err != nil {
		s.logger.Error("创建计划失败", zap.Error(err))
		return errors.New("创建计划失败")
	}
	return nil
}

// ListPlans 获取产品的计划列表
func (s *ProductService) ListPlans(ctx context.Context, productID uuid.UUID) ([]*entity.ProductPlan, error) {
	// 检查产品是否存在
	if _, err := s.productRepo.GetByID(ctx, productID); err != nil {
		return nil, ErrProductNotFound
	}
	return s.planRepo.ListByProductID(ctx, productID)
}

// ListAllPlans 分页查询全部计划
func (s *ProductService) ListAllPlans(ctx context.Context, page, pageSize int) ([]*entity.ProductPlan, int64, error) {
	return s.planRepo.List(ctx, page, pageSize)
}

// UpdatePlan 更新计划
func (s *ProductService) UpdatePlan(ctx context.Context, plan *entity.ProductPlan) error {
	existing, err := s.planRepo.GetByID(ctx, plan.ID)
	if err != nil || existing == nil {
		return ErrPlanNotFound
	}
	return s.planRepo.Update(ctx, plan)
}

// GetPlanDetail 获取计划详情
func (s *ProductService) GetPlanDetail(ctx context.Context, id uuid.UUID) (*entity.ProductPlan, error) {
	plan, err := s.planRepo.GetByID(ctx, id)
	if err != nil || plan == nil {
		return nil, ErrPlanNotFound
	}
	return plan, nil
}

// DeletePlan 删除计划
func (s *ProductService) DeletePlan(ctx context.Context, id uuid.UUID) error {
	existing, err := s.planRepo.GetByID(ctx, id)
	if err != nil || existing == nil {
		return ErrPlanNotFound
	}
	return s.planRepo.Delete(ctx, id)
}
