package services

import "time"

/* 卖家 */

// 卖家创建订单请求：
type CreateSellerOrderRequest struct {
	Symbol    string  `json:"symbol" validate:"required,alphanum"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,gt=0"`
	OrderType string  `json:"type" validate:"required,oneof=market limit"`
}

// 卖家更新订单请求：
type UpdateSellerOrderRequest struct {
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required,gt=0"`
}

// 创建授权请求：
type CreateAuthorizationRequest struct {
	SalesID   uint      `json:"sales_id" validate:"required"`
	ExpiresAt time.Time `json:"expires_at" validate:"required"`
}

/* 销售 */

// 创建草稿请求：
type CreateDraftRequest struct {
	SellerID   uint    `json:"seller_id" validate:"required"`
	Symbol     string  `json:"symbol" validate:"required,alphanum"`
	Quantity   int     `json:"quantity" validate:"required,min=1"`
	Price      float64 `json:"price" validate:"required,gt=0"`
	RefOrderID *uint   `json:"ref_order_id"` // 关联原订单ID（可选）
}

// 更新草稿请求：
type UpdateDraftRequest struct {
	Quantity int     `json:"quantity" validate:"min=1"`
	Price    float64 `json:"price" validate:"gt=0"`
}

/* 客户 */

// 客户创建订单请求：
type CreateBuyOrderRequest struct {
	Symbol   string  `json:"symbol" validate:"required,alphanum"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required,gt=0"`
}

/* 所有用户 */

// 用户注册请求：
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // 角色: client/sales/trader
}

// 用户登录请求：
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
