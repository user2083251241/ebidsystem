// 单个撤单（软删除，使用 DELETE 方法）：
func CancelOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	// 从 JWT 中获取卖家 ID:
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	var order models.Order
	// 检查订单是否存在：
	if err := db.Where("id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "订单不存在"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "查询订单失败"})
	}
	// 检查订单归属权：
	if order.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "无权操作此订单"})
	}
	// 检查订单状态：
	if order.Status == "cancelled" {
		return c.Status(400).JSON(fiber.Map{"error": "订单已取消，不可重复操作"})
	}
	// 添加事物：
	tx := db.Begin()
	// 执行软删除（设置 DeletedAt）：
	if err := db.Delete(&order).Error; err != nil {
		tx.Rollback() //回滚事务
		return c.Status(500).JSON(fiber.Map{"error": "取消失败"})
	}
	// 更新状态为 "cancelled"（需使用 Unscoped 更新软删除记录）：
	if err := db.Unscoped().Model(&order).Update("status", "cancelled").Error; err != nil {
		tx.Rollback() //回滚事务
		return c.Status(500).JSON(fiber.Map{"error": "状态更新失败"})
	}
	// 提交事务：
	tx.Commit()
	// 返回成功信息：
	return c.JSON(fiber.Map{"message": "订单已取消"})
}

// 批量撤单（软删除，使用 POST 方法）：
func BatchCancelOrders(c *fiber.Ctx) error {
	type BatchCancelRequest struct {
		OrderIDs []uint `json:"order_ids"`
	}
	var req BatchCancelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "请求格式错误"})
	}
	// 获取当前用户 ID：
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	// 开启事物：
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 检查每个订单状态：
	for _, orderID := range req.OrderIDs {
		var order models.Order
		// 检查订单是否存在：
		if err := tx.Where("id = ?", orderID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return c.Status(404).JSON(fiber.Map{
					"error": fmt.Sprintf("订单 %d 不存在", orderID),
				})
			}
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "查询订单失败"})
		}
		// 检查订单归属权：
		if order.UserID != userID {
			tx.Rollback()
			return c.Status(403).JSON(fiber.Map{
				"error": fmt.Sprintf("订单 %d 无权操作", orderID),
			})
		}
		// 检查订单状态：
		if order.Status == "cancelled" {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("订单 %d 已取消，不可重复操作", orderID),
			})
		}
	}
	// 批量软删除并更新状态：
	if err := tx.Model(&models.Order{}).
		Where("user_id = ? AND id IN ?", userID, req.OrderIDs).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("CURRENT_TIMESTAMP"),
			"status":     "cancelled",
		}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "批量取消失败"})
	}

	tx.Commit()
	return c.JSON(fiber.Map{"message": "批量取消成功"})
}

// 卖家查看自己订单（隐藏软删除）：
func GetSellerOrders(c *fiber.Ctx) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	var orders []models.Order
	// 默认查询会排除 DeletedAt 非空的记录
	if err := db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询订单失败"})
	}
	return c.JSON(orders)
}