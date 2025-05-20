package entity

type Stock struct {
	Symbol string  `gorm:"primaryKey"`         // 股票代码（主键）
	Name   string  `gorm:"not null"`           // 股票名称
	Price  float64 `gorm:"type:decimal(10,2)"` // 当前价格
}

// 表名自定义（可选）
func (Stock) TableName() string {
	return "stocks"
}
