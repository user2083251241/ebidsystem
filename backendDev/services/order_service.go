package services

import (
	"time"

	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrderService struct {
	DB *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{DB: db}
}

type CreateOrderRequest struct {
	Symbol    string  `json:"symbol"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`     // 限价单必填
	Direction string  `json:"direction"` // "buy" 或 "sell"
	OrderType string  `json:"type"`      // "market" 或 "limit"
}

// 创建订单（判断方向和角色）：
func (s *OrderService) CreateOrder(user *models.User, req CreateOrderRequest) (*models.Order, error) {
	// 权限校验（如仅允许交易员创建卖出订单）：
	if req.Direction == "sell" && user.Role != "trader" {
		return nil, fiber.NewError(fiber.StatusForbidden, "无权创建卖出订单")
	}
	order := models.LiveOrder{
		UserID:    user.ID,
		Symbol:    req.Symbol,
		Quantity:  req.Quantity,
		Price:     req.Price,
		Direction: req.Direction,
		OrderType: req.OrderType,
		Status:    "pending",
	}
	if err := s.DB.Create(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// QueryCondition 定义通用查询条件
type QueryCondition struct {
	UserID        *uint   // 用户ID过滤
	Direction     *string // 买卖方向过滤
	Status        *string // 订单状态过滤
	HideSensitive bool    // 是否隐藏敏感字段（如数量、状态）
}

// GetOrders 通用订单查询函数：
func (s *OrderService) GetOrders(conditions QueryCondition) ([]models.OrderDTO, error) {
	// 初始化查询：
	query := s.DB.Model(&models.Order{})
	// 动态注入查询条件：
	if conditions.UserID != nil {
		query = query.Where("user_id = ?", *conditions.UserID)
	}
	if conditions.Direction != nil {
		query = query.Where("direction = ?", *conditions.Direction)
	}
	if conditions.Status != nil {
		query = query.Where("status = ?", *conditions.Status)
	}
	// 执行查询：
	var orders []models.Order
	if err := query.Find(&orders).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "查询失败")
	}
	// 转换为 DTO 并脱敏：
	var dtos []models.OrderDTO
	for _, order := range orders {
		dto := models.OrderDTO{
			ID:     order.ID,
			Symbol: order.Symbol,
			Price:  order.Price,
		}
		// 根据条件显示敏感字段：
		if !conditions.HideSensitive {
			dto.Quantity = order.Quantity
			dto.Status = order.Status
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

// 创建销售草稿：
func (s *OrderService) CreateDraftOrder(salesID, sellerID uint, req CreateOrderRequest, refOrderID *uint) (*models.Order, error) {
	if err := s.checkSalesAuthorization(salesID, sellerID); err != nil {
		return nil, err
	}
	draft := models.Order{
		UserID:       req.SellerID, // 需从请求中获取卖家ID
		Symbol:       req.Symbol,
		Quantity:     req.Quantity,
		Price:        req.Price,
		Direction:    "sell",
		Status:       "draft",
		DraftBySales: user.ID,
		RefOrderID:   refOrderID, // 如果是修改建议，关联原订单ID
	}
	if err := s.DB.Create(&draft).Error; err != nil {
		return nil, err
	}
	return &draft, nil
}
func (s *OrderService) UpdateDraftOrder(draftID uint, salesID uint, updates map[string]interface{}) error {
	var draft models.Order
	if err := s.DB.Where("id = ? AND draft_by_sales = ? AND status = 'draft'", draftID, salesID).First(&draft).Error; err != nil {
		return fiber.NewError(fiber.StatusForbidden, "草稿不可修改")
	}
	return s.DB.Model(&draft).Updates(updates).Error
}

// 销售提交草稿审批：
func (s *OrderService) SubmitDraftOrder(draftID uint, salesID uint) error {
	if err := s.checkSalesAuthorization(salesID, draft.SellerID); err != nil {
		return err
	}
	return s.DB.Model(&models.Order{}).Where("id = ?", draftID).Update("status", "pending_approval").Error
}

// 卖家审批草稿：
func (s *OrderService) ApproveDraftOrder(draftID uint, sellerID uint) (*models.Order, error) {
	var draft models.Order
	if err := s.DB.Where("id = ? AND user_id = ? AND status = 'pending_approval'", draftID, sellerID).First(&draft).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "草稿不存在")
	}
	liveOrder := models.Order{
		UserID:    draft.UserID,
		Symbol:    draft.Symbol,
		Quantity:  draft.Quantity,
		Price:     draft.Price,
		Direction: "sell",
		Status:    "pending",
	}
	if err := s.DB.Create(&liveOrder).Error; err != nil {
		return nil, err
	}
	return &liveOrder, nil
}

// 检查销售是否被卖家授权：
func (s *OrderService) checkSalesAuthorization(salesID, sellerID uint) error {
	var auth models.SellerSalesAuthorization
	// 检查授权状态、有效期及审批层级：
	if err := s.DB.Where(
		"seller_id = ? AND sales_id = ? AND authorization = 'approved' AND expires_at > ? AND approval_level >= ?",
		sellerID, salesID, time.Now(), 1, // approval_level 表示需要的最低审批层级
	).First(&auth).Error; err != nil {
		return fiber.NewError(fiber.StatusForbidden, "未获得有效卖家授权或授权层级不足")
	}
	return nil
}

// 交易员查看所有方向和状态的订单：
func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	if err := s.DB.Unscoped().Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
