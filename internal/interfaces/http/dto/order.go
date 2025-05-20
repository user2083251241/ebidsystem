package dto

import (
	"time"
)

// 创建订单请求
type CreateOrderRequest struct {
	Symbol    string  `json:"symbol" validate:"required,alphanum,max=255"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,gt=0"`
	OrderType string  `json:"type" validate:"required,oneof=market limit"`
}

// 更新订单请求
type UpdateOrderRequest struct {
	Quantity int     `json:"quantity" validate:"min=1"`
	Price    float64 `json:"price" validate:"gt=0"`
}

// 订单响应
type OrderResponse struct {
	ID        uint      `json:"id"`
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// 订单统计响应
type OrderStatsResponse struct {
	TotalOrders   int       `json:"total_orders"`
	AveragePrice  float64   `json:"average_price"`
	TotalQuantity int       `json:"total_quantity"`
	LastUpdated   time.Time `json:"last_updated"`
}
