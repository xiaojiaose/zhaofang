package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/gin-gonic/gin"
)

type WxPayApi struct{}

// CreatePayOrder 创建支付订单
// @Tags      WxPay
// @Summary   创建支付订单
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      CreatePayOrderReq              true  "支付订单参数"
// @Success   200   {object}  response.Response{data=PayOrderResp}  "创建成功"
// @Router    /wxpay/createOrder [post]
func (w *WxPayApi) CreatePayOrder(c *gin.Context) {
	var req CreatePayOrderReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	order, payParams, err := wxPayService.CreateOrder(req.OpenID, req.Body, req.TotalFee)
	if err != nil {
		response.FailWithMessage("创建订单失败: "+err.Error(), c)
		return
	}

	resp := PayOrderResp{
		OrderNo:   order.OrderNo,
		PayParams: payParams,
	}

	response.OkWithData(resp, c)
}

// PayNotify 支付回调
// @Tags      WxPay
// @Summary   支付回调
// @accept    application/json
// @Produce   application/json
// @Router    /wxpay/notify [post]
func (w *WxPayApi) PayNotify(c *gin.Context) {
	err := wxPayService.HandleNotify(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code": "FAIL",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": "SUCCESS",
		"message": "OK",
	})
}

// QueryOrder 查询订单状态
// @Tags      WxPay
// @Summary   查询订单状态
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     orderNo  query     string                         true  "订单号"
// @Success   200      {object}  response.Response{data=system.PayOrder}  "查询成功"
// @Router    /wxpay/queryOrder [get]
func (w *WxPayApi) QueryOrder(c *gin.Context) {
	orderNo := c.Query("orderNo")
	if orderNo == "" {
		response.FailWithMessage("订单号不能为空", c)
		return
	}

	var order system.PayOrder
	err := global.GVA_DB.Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		response.FailWithMessage("订单不存在", c)
		return
	}

	response.OkWithData(order, c)
}

// QueryWxOrder 查询微信订单状态
// @Tags      WxPay
// @Summary   查询微信订单状态
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     outTradeNo  query     string                         true  "商户订单号"
// @Success   200         {object}  response.Response{data=object}  "查询成功"
// @Router    /wxpay/queryWxOrder [get]
func (w *WxPayApi) QueryWxOrder(c *gin.Context) {
	outTradeNo := c.Query("outTradeNo")
	if outTradeNo == "" {
		response.FailWithMessage("商户订单号不能为空", c)
		return
	}

	transaction, err := wxPayService.QueryOrder(outTradeNo)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}

	response.OkWithData(transaction, c)
}

// CreatePayOrderReq 创建支付订单请求
type CreatePayOrderReq struct {
	OpenID   string `json:"openId" binding:"required"`   // 用户openid
	Body     string `json:"body" binding:"required"`     // 商品描述
	TotalFee int    `json:"totalFee" binding:"required"` // 支付金额(分)
}

// PayOrderResp 支付订单响应
type PayOrderResp struct {
	OrderNo   string                 `json:"orderNo"`   // 订单号
	PayParams map[string]interface{} `json:"payParams"` // 小程序支付参数
}
