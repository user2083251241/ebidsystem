package controllers

import (
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// 销售查看已授权的卖家订单
func GetAuthorizedOrders(c *fiber.Ctx) error {
	// 获取当前销售 ID
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	salesID := uint(claims["user_id"].(float64))

	// 查询已授权的卖家订单
	var orders []models.Order
	if err := db.Joins("JOIN seller_sales_authorizations ON seller_sales_authorizations.seller_id = orders.user_id").
		Where("seller_sales_authorizations.sales_id = ? AND seller_sales_authorizations.authorization = 'approved'", salesID).
		Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询订单失败"})
	}

	return c.JSON(orders)
}

// 销售创建草稿订单
func CreateDraftOrder(c *fiber.Ctx) error {
	type DraftRequest struct {
		Symbol   string  `json:"symbol"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}
	var req DraftRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}

	// 获取当前销售 ID 并验证授权
	salesID := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	var auth models.SellerSalesAuthorization
	if err := db.Where("sales_id = ? AND authorization = 'approved'", salesID).First(&auth).Error; err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "未获得卖家授权"})
	}

	// 创建草稿订单
	draft := models.Order{
		UserID:       auth.SellerID,
		Symbol:       req.Symbol,
		Quantity:     req.Quantity,
		Price:        req.Price,
		Direction:    "sell",
		Status:       "draft",
		DraftBySales: salesID,
	}
	if err := db.Create(&draft).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "草稿创建失败"})
	}

	return c.JSON(draft)
}

// 销售修改草稿订单
func UpdateDraftOrder(c *fiber.Ctx) error {
	draftID := c.Params("id")
	type UpdateDraftRequest struct {
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}
	var req UpdateDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}

	// 验证草稿归属
	var draft models.Order
	if err := db.Where("id = ? AND status = 'draft'", draftID).First(&draft).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "草稿不存在或不可修改"})
	}

	// 更新草稿
	if err := db.Model(&draft).Updates(models.Order{
		Quantity: req.Quantity,
		Price:    req.Price,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新草稿失败"})
	}

	return c.JSON(draft)
}

// 销售提交草稿给卖家审批
func SubmitDraft(c *fiber.Ctx) error {
	draftID := c.Params("id")

	// 验证草稿状态
	var draft models.Order
	if err := db.Where("id = ? AND status = 'draft'", draftID).First(&draft).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "草稿不存在或不可提交"})
	}

	// 更新状态为待审批
	if err := db.Model(&draft).Update("status", "pending_approval").Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "提交失败"})
	}

	return c.JSON(fiber.Map{"message": "草稿已提交，等待卖家审批"})
}
