package services

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/champNoob/ebidsystem/backend/utils"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrderService struct {
	db        *gorm.DB
	validator *validator.Validate
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		db:        db,
		validator: validator.New(),
	}
}

// ---------------------- 卖家订单操作 ----------------------
// 创建卖家订单：
func (s *OrderService) CreateSellerOrder(user *models.User, req CreateSellerOrderRequest) (*models.LiveOrder, error) {
	// 检查用户角色是否为卖家：
	if user.Role != "seller" {
		return nil, fiber.NewError(fiber.StatusForbidden, "invalid role")
	}
	// 验证 req 结构体：
	if err := utils.ValidateStruct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	// 构造订单：
	order := models.LiveOrder{
		BaseOrder: models.BaseOrder{
			Symbol:    req.Symbol,
			Quantity:  req.Quantity,
			Price:     req.Price,
			CreatorID: user.ID,
		},
		Direction: "sell",
		Status:    "pending",
	}
	// 正式创建订单，并存入数据库：
	if err := s.db.Create(&order).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "create failed")
	}

	return &order, nil
}

// 更新卖家订单：
func (s *OrderService) UpdateSellerOrder(sellerID uint, orderID uint, req UpdateSellerOrderRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.LiveOrder{}).
			Where("id = ? AND creator_id = ? AND status = 'pending'", orderID, sellerID).
			Updates(map[string]interface{}{
				"quantity": req.Quantity,
				"price":    req.Price,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fiber.NewError(fiber.StatusNotFound, "订单不存在或不可修改")
		}
		return nil
	})
}

// 获取卖家订单：
func (s *OrderService) GetSellerOrders(sellerID uint) ([]models.LiveOrder, error) {
	if s.db == nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "数据库实例未初始化")
	}
	var orders []models.LiveOrder
	query := s.db.Where("creator_id = ? AND direction = 'sell'", sellerID).Find(&orders)
	log.Printf("Executing SQL query: %v", query.Statement.SQL.String())
	log.Printf("Query conditions: %v", query.Statement.Vars)
	if err := query.Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve seller orders")
	}
	return orders, nil
}

// 取消单个卖家订单：
func (s *OrderService) CancelSellerOrder(sellerID uint, orderID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.LiveOrder{}).
			Where("id = ? AND creator_id = ? AND status = 'pending'", orderID, sellerID).
			Updates(map[string]interface{}{
				"status":     "cancelled",
				"deleted_at": gorm.Expr("CURRENT_TIMESTAMP"),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fiber.NewError(fiber.StatusNotFound, "无有效订单可取消")
		}
		return nil
	})
}

// 批量取消卖家订单：
func (s *OrderService) BatchCancelSellerOrders(sellerID uint, orderIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&models.LiveOrder{}).
			Where("creator_id = ? AND id IN ? AND status = 'pending'", sellerID, orderIDs).
			Updates(map[string]interface{}{
				"status":     "cancelled",
				"deleted_at": gorm.Expr("CURRENT_TIMESTAMP"),
			}).Error
	})
}

// ---------------------- 销售草稿操作 ----------------------
func (s *OrderService) CreateDraftOrder(salesID uint, req CreateDraftRequest) (*models.DraftOrder, error) {
	if err := s.checkSalesAuthorization(salesID, req.SellerID); err != nil {
		return nil, err
	}

	draft := models.DraftOrder{
		BaseOrder: models.BaseOrder{
			Symbol:    req.Symbol,
			Quantity:  req.Quantity,
			Price:     req.Price,
			CreatorID: salesID,
		},
		RefOrderID: req.RefOrderID,
		Status:     "draft",
	}

	if err := s.db.Create(&draft).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "create draft failed")
	}

	return &draft, nil
}

// 更新草稿订单
func (s *OrderService) UpdateDraftOrder(salesID uint, draftID uint, req UpdateDraftRequest) error {
	result := s.db.Model(&models.DraftOrder{}).
		Where("id = ? AND draft_by_sales = ?", draftID, salesID).
		Updates(map[string]interface{}{
			"quantity": req.Quantity,
			"price":    req.Price,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "草稿不存在或无权修改")
	}
	return nil
}

// 提交草稿审批
func (s *OrderService) SubmitDraftOrder(salesID uint, draftID uint) error {
	var draft models.DraftOrder
	if err := s.db.Where("id = ? AND draft_by_sales = ?", draftID, salesID).First(&draft).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "draft not found")
	}

	return s.db.Model(&draft).Update("status", "pending_approval").Error
}

// 获取已授权草稿：
func (s *OrderService) GetAuthorizedDrafts(salesID uint) ([]models.DraftOrder, error) {
	var drafts []models.DraftOrder
	err := s.db.
		Joins("JOIN seller_sales_authorizations ON seller_sales_authorizations.seller_id = draft_orders.creator_id").
		Where("seller_sales_authorizations.sales_id = ? AND authorization = 'approved'", salesID).
		Find(&drafts).Error
	return drafts, err
}

// 删除草稿：
func (s *OrderService) DeleteDraft(salesID uint, draftID uint) error {
	return s.db.Where("id = ? AND draft_by_sales = ?", draftID, salesID).Delete(&models.DraftOrder{}).Error
}

// ---------------------- 客户操作 ----------------------

// 客户创建买入订单：
func (s *OrderService) CreateClientBuyOrder(clientID uint, req CreateBuyOrderRequest) (*models.LiveOrder, error) {
	order := models.LiveOrder{
		BaseOrder: models.BaseOrder{
			Symbol:    req.Symbol,
			Quantity:  req.Quantity,
			Price:     req.Price,
			CreatorID: clientID,
		},
		Direction: "buy",
		Status:    "pending",
	}

	if err := s.db.Create(&order).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create buy order")
	}

	return &order, nil
}

// 获取可卖订单：
func (s *OrderService) GetAvailableSellOrders() ([]models.OrderDTO, error) {
	strategy := &ClientOrdersStrategy{}
	return s.QueryOrders(strategy)
}

// ---------------------- 交易员操作 --------------------

// 获取所有订单：
func (s *OrderService) GetAllOrders() ([]models.LiveOrder, error) {
	var orders []models.LiveOrder
	if err := s.db.Unscoped().Find(&orders).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve orders")
	}
	return orders, nil
}

// 紧急撤单：
func (s *OrderService) EmergencyCancelOrder(orderID uint, traderID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var order models.LiveOrder
		if err := tx.Unscoped().First(&order, orderID).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "order not found")
		}

		// 记录审计日志
		if err := tx.Create(&models.AuditLog{
			UserID:    traderID,
			Action:    "emergency_cancel",
			OrderID:   orderID,
			Timestamp: time.Now(),
		}).Error; err != nil {
			return err
		}

		return tx.Model(&order).Update("status", "cancelled").Error
	})
}

// ---------------------- 卖家-销售授权 ----------------------

// 创建销售授权：
func (s *OrderService) CreateSalesAuthorization(sellerID uint, req CreateAuthorizationRequest) (*models.SellerSalesAuthorization, error) {
	auth := models.SellerSalesAuthorization{
		SellerID:      sellerID,
		SalesID:       req.SalesID,
		Authorization: "pending",
		ExpiresAt:     req.ExpiresAt,
	}

	if err := s.db.Create(&auth).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "authorization creation failed")
	}

	return &auth, nil
}

// ---------------------- 所有角色 -------------------------

// 取消用户所有未完成订单：
func (s *OrderService) CancelUserUnfinishedOrders(userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 取消用户创建的正式订单
		if err := tx.Model(&models.LiveOrder{}).
			Where("creator_id = ? AND status = ?", userID, "pending").
			Updates(map[string]interface{}{
				"status": "cancelled",
			}).Error; err != nil {
			log.Printf("取消用户创建的订单失败 | 用户ID: %d | 错误: %v", userID, err)
			return fiber.NewError(fiber.StatusInternalServerError, "取消订单失败："+err.Error())
		}

		// 取消用户相关的草稿订单
		if err := tx.Model(&models.DraftOrder{}).
			Where("creator_id = ? AND status = ?", userID, "draft").
			Updates(map[string]interface{}{
				"status": "cancelled",
			}).Error; err != nil {
			// 如果表不存在，跳过这个错误（因为用户可能没有草稿订单）
			if !strings.Contains(err.Error(), "doesn't exist") {
				log.Printf("取消用户相关草稿订单失败 | 用户ID: %d | 错误: %v", userID, err)
				return fiber.NewError(fiber.StatusInternalServerError, "取消草稿订单失败："+err.Error())
			}
			log.Printf("草稿订单表不存在，跳过 | 用户ID: %d", userID)
		}

		// 取消用户的授权关系
		if err := tx.Where("seller_id = ? OR sales_id = ?", userID, userID).
			Delete(&models.SellerSalesAuthorization{}).Error; err != nil {
			// 如果表不存在，跳过这个错误（因为用户可能没有授权关系）
			if !strings.Contains(err.Error(), "doesn't exist") {
				log.Printf("取消用户授权关系失败 | 用户ID: %d | 错误: %v", userID, err)
				return fiber.NewError(fiber.StatusInternalServerError, "取消授权关系失败："+err.Error())
			}
			log.Printf("授权关系表不存在，跳过 | 用户ID: %d", userID)
		}

		log.Printf("已成功取消所有相关订单和授权 | 用户ID: %d", userID)
		return nil
	})
}

// ---------------------- 通用查询引擎 ----------------------
func (s *OrderService) QueryOrders(strategy QueryStrategy) ([]models.OrderDTO, error) {
	var orders []models.LiveOrder
	query := s.db.Model(&models.LiveOrder{})

	query = strategy.Apply(query)

	if err := query.Find(&orders).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "query failed")
	}

	// DTO转换
	converter := strategy.GetDTOConverter()
	dtos := make([]models.OrderDTO, len(orders))
	for i, o := range orders {
		dtos[i] = converter(o)
	}

	return dtos, nil
}

// ---------------------- 工具函数 ----------------------
func (s *OrderService) checkSalesAuthorization(salesID, sellerID uint) error {
	var auth models.SellerSalesAuthorization
	err := s.db.Where(
		"seller_id = ? AND sales_id = ? AND authorization = 'approved' AND expires_at > ?",
		sellerID, salesID, time.Now(),
	).First(&auth).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusForbidden, "unauthorized access")
	}

	return err
}

// 验证请求结构体：
func (s *OrderService) Validate(req interface{}) error {
	return s.validator.Struct(req)
}
