package controllers

import (
	"github.com/champNoob/ebidsystem/backend/middleware"
	"github.com/champNoob/ebidsystem/backend/models"

	// "github.com/champNoob/ebidsystem/backend/services"
	"github.com/gofiber/fiber/v2"
	"github.com/user2083251241/ebidsystem/services"
	"github.com/user2083251241/ebidsystem/utils"
)

// controllers/sales.go
type SalesController struct {
	orderService *services.OrderService
}

func NewSalesController(os *services.OrderService) *SalesController {
	return &SalesController{orderService: os}
}

// SalesCreateDraft 销售创建草稿订单
func (sc *SalesController) SalesCreateDraftOrder(c *fiber.Ctx) error {
	user, _ := middleware.GetCurrentUserFromJWT(c)
	var req services.CreateDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	draft, err := sc.orderService.CreateDraftOrder(user.ID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(draft)
}

// SalesSubmitDraft 销售提交草稿审批
func (sc *SalesController) SalesSubmitDraftOrder(c *fiber.Ctx) error {
	user, _ := middleware.GetCurrentUserFromJWT(c)
	draftID, _ := c.ParamsInt("id")
	if err := sc.orderService.SubmitDraft(user.ID, uint(draftID)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(200)
}

// 销售查看已授权的卖家订单：
func GetAuthorizedOrders(c *fiber.Ctx) error {
	// 获取当前销售 ID
	salesID, err := getCurrentUserIDFromDB(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "从 JWT 中获取卖家 ID 失败"})
	}

	// 查询已授权的卖家订单：
	var orders []models.Order
	if err := db.Joins("JOIN seller_sales_authorizations ON seller_sales_authorizations.seller_id = orders.user_id").
		Where("seller_sales_authorizations.sales_id = ? AND seller_sales_authorizations.authorization = 'approved'", salesID).
		Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询订单失败"})
	}

	return c.JSON(orders)
}

// 销售查看已授权的草稿订单
func (sc *SalesController) GetDrafts(c *fiber.Ctx) error {
	user, _ := middleware.GetCurrentUserFromJWT(c)
	// 构造查询条件：仅显示草稿状态且关联当前销售
	conditions := services.QueryCondition{
		Status:        utils.Ptr("draft"),
		HideSensitive: true, // 对销售隐藏卖家ID等字段
	}
	// 附加 JOIN 条件（需在服务层扩展或通过额外参数处理）
	// 此处简化逻辑，实际需结合 JOIN
	orders, err := sc.orderService.GetOrders(conditions)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}
