package system

import (
	"gorm.io/gorm"
	"time"
)

// PayOrder 支付订单
type PayOrder struct {
	ID            uint   `gorm:"primarykey" json:"ID"`
	OrderNo       string `json:"orderNo" gorm:"uniqueIndex"`       // 订单号
	OutTradeNo    string `json:"outTradeNo" gorm:"uniqueIndex"`    // 商户订单号
	TransactionID string `json:"transactionId"`                    // 微信支付订单号
	OpenID        string `json:"openId"`                           // 用户openid
	TotalFee      int    `json:"totalFee"`                         // 支付金额(分)
	Body          string `json:"body"`                             // 商品描述
	Status        int    `json:"status" gorm:"default:0"`          // 订单状态 0:待支付 1:已支付 2:已取消
	PayTime       *time.Time `json:"payTime"`                      // 支付时间
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PayOrder) TableName() string {
	return "pay_orders"
}
