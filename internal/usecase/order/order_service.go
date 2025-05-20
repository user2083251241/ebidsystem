package orderusecase

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/user2083251241/ebidsystem/internal/domain/entity"
	"github.com/user2083251241/ebidsystem/internal/domain/repository"
	"github.com/user2083251241/ebidsystem/internal/interfaces/http/dto"
	"gorm.io/gorm"
)

type OrderService struct {
	db        *gorm.DB
	validator *validator.Validate
	orderRepo repository.OrderRepository
}

func NewOrderService(db *gorm.DB, orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		db:        db,
		validator: validator.New(),
		orderRepo: orderRepo,
	}
}

// ---------------------- 卖家订单操作 ----------------------
// 创建卖家订单：
func (s *OrderService) CreateSellerOrder(user *entity.User, req dto.CreateOrderRequest) (*entity.LiveOrder, error) {
	// 检查用户角色是否为卖家：
	if user.Role != "seller" {
		return nil, fiber.NewError(fiber.StatusForbidden, "invalid role")
	}
	// 验证 req 结构体：
	if err := s.validator.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	// 构造订单：
	order := &entity.LiveOrder{
		BaseOrder: entity.BaseOrder{
			Symbol:    req.Symbol,
			Quantity:  req.Quantity,
			Price:     req.Price,
			CreatorID: user.ID,
		},
		Direction: "sell",
		Status:    entity.StatusPending,
	}
	// 正式创建订单：
	if err := s.orderRepo.CreateLive(context.Background(), order); err != nil {
		return nil, err
	}
	return order, nil
}

// ---------------------- 客户订单操作 ----------------------
// 创建客户订单：
func (s *OrderService) CreateBuyOrder(user *entity.User, req dto.CreateOrderRequest) (*entity.LiveOrder, error) {
	// 检查用户角色是否为客户：
	if user.Role != "client" {
		return nil, fiber.NewError(fiber.StatusForbidden, "invalid role")
	}
	// 验证 req 结构体：
	if err := s.validator.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	// 构造订单：
	order := &entity.LiveOrder{
		BaseOrder: entity.BaseOrder{
			Symbol:    req.Symbol,
			Quantity:  req.Quantity,
			Price:     req.Price,
			CreatorID: user.ID,
		},
		Direction: "buy",
		Status:    entity.StatusPending,
	}
	// 正式创建订单：
	if err := s.orderRepo.CreateLive(context.Background(), order); err != nil {
		return nil, err
	}
	return order, nil
}

// ---------------------- 销售订单操作 ----------------------
// 创建草稿订单：
func (s *OrderService) CreateDraftOrder(user *entity.User, req dto.CreateDraftRequest) (*entity.DraftOrder, error) {
	// 检查用户角色是否为销售：
	if user.Role != "sales" {
		return nil, fiber.NewError(fiber.StatusForbidden, "invalid role")
	}
	// 验证 req 结构体：
	if err := s.validator.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	// 构造草稿订单：
	order := &entity.DraftOrder{
		BaseOrder: entity.BaseOrder{
			Symbol:    req.Symbol,
			Quantity:  req.Quantity,
			Price:     req.Price,
			CreatorID: user.ID,
		},
		RefOrderID: req.RefOrderID,
		Status:     entity.StatusDraft,
	}
	// 正式创建草稿订单：
	if err := s.orderRepo.CreateDraft(context.Background(), order); err != nil {
		return nil, err
	}
	return order, nil
}
