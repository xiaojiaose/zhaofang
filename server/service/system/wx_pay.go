package system

import (
	"context"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"io"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"go.uber.org/zap"
)

type WxPayService struct{}

var (
	appID     = "your_appid"
	mchID     = "your_mch_id"
	notifyURL = "https://your-domain.com/api/wxpay/notify"
	// 商户私钥文件路径
	privateKeyPath = "path/to/apiclient_key.pem"
	// 商户证书序列号
	mchCertificateSerialNumber = "your_certificate_serial_number"
	// APIv3密钥
	apiV3Key = "your_apiv3_key"

	clientOnce sync.Once
	client     *core.Client
)

// getClient 获取微信支付客户端
func (w *WxPayService) getClient() (*core.Client, error) {
	clientOnce.Do(func() {
		mchPrivateKey, err := utils.LoadPrivateKeyWithPath(privateKeyPath)
		if err != nil {
			global.GVA_LOG.Error("load merchant private key error", zap.Error(err))
			return
		}

		ctx := context.Background()
		opts := []core.ClientOption{
			option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, apiV3Key),
		}

		client, err = core.NewClient(ctx, opts...)
		if err != nil {
			global.GVA_LOG.Error("new wechat pay client error", zap.Error(err))
			return
		}
	})

	return client, nil
}

// CreateOrder 创建支付订单
func (w *WxPayService) CreateOrder(openID, body string, totalFee int) (*system.PayOrder, map[string]interface{}, error) {
	orderNo := fmt.Sprintf("%d", time.Now().UnixNano())

	order := &system.PayOrder{
		OrderNo:    orderNo,
		OutTradeNo: orderNo,
		OpenID:     openID,
		TotalFee:   totalFee,
		Body:       body,
		Status:     0,
	}

	err := global.GVA_DB.Create(order).Error
	if err != nil {
		return nil, nil, err
	}

	payParams, err := w.createJSAPIOrder(order)
	if err != nil {
		return nil, nil, err
	}

	return order, payParams, nil
}

// createJSAPIOrder 创建JSAPI支付订单
func (w *WxPayService) createJSAPIOrder(order *system.PayOrder) (map[string]interface{}, error) {
	wclient, err := w.getClient()
	if err != nil {
		return nil, err
	}

	svc := jsapi.JsapiApiService{Client: wclient}

	req := jsapi.PrepayRequest{
		Appid:       core.String(appID),
		Mchid:       core.String(mchID),
		Description: core.String(order.Body),
		OutTradeNo:  core.String(order.OutTradeNo),
		NotifyUrl:   core.String(notifyURL),
		Amount: &jsapi.Amount{
			Total: core.Int64(int64(order.TotalFee)),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(order.OpenID),
		},
	}

	resp, _, err := svc.PrepayWithRequestPayment(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"timeStamp": resp.TimeStamp,
		"nonceStr":  resp.NonceStr,
		"package":   resp.Package,
		"signType":  resp.SignType,
		"paySign":   resp.PaySign,
	}, nil
}

// HandleNotify 处理支付回调
func (w *WxPayService) HandleNotify(c *gin.Context) error {
	body, _ := io.ReadAll(c.Request.Body)
	global.GVA_LOG.Info("1 微信支付回调 开始", zap.String("body", string(body)))

	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(privateKeyPath)
	if err != nil {
		global.GVA_LOG.Error("load merchant private key error", zap.Error(err))
		return err
	}
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(c, mchPrivateKey, mchCertificateSerialNumber, mchID, apiV3Key)
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchID)
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler := notify.NewNotifyHandler(apiV3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))

	transaction := new(payments.Transaction)
	notifyReq, err := handler.ParseNotifyRequest(context.Background(), c.Request, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		global.GVA_LOG.Error("2 微信支付回调验签失败", zap.Error(err))
		fmt.Println(err)
		return err
	}
	// 处理通知内容
	fmt.Println(notifyReq.Summary)
	fmt.Println(transaction.TransactionId)

	global.GVA_LOG.Info("3 微信支付回调解密成功", zap.String("trade_state", *transaction.TradeState), zap.String("transaction_id", *transaction.TransactionId), zap.String("out_trade_no", *transaction.OutTradeNo), zap.String("body", transaction.String()))

	// 验证支付状态
	if *transaction.TradeState != "SUCCESS" {
		global.GVA_LOG.Error("4 微信支付失败", zap.String("trade_state", *transaction.TradeState))
		return fmt.Errorf("payment not success: %s", *transaction.TradeState)
	}

	// 更新订单状态
	now := time.Now()
	err = global.GVA_DB.Model(&system.PayOrder{}).
		Where("out_trade_no = ?", *transaction.OutTradeNo).
		Updates(map[string]interface{}{
			"status":         1,
			"transaction_id": *transaction.TransactionId,
			"pay_time":       &now,
		}).Error

	if err != nil {
		global.GVA_LOG.Error("update order failed", zap.Error(err))
		return err
	}

	return nil
}

// QueryOrder 查询订单状态
func (w *WxPayService) QueryOrder(outTradeNo string) (*payments.Transaction, error) {
	client, err := w.getClient()
	if err != nil {
		return nil, err
	}

	svc := jsapi.JsapiApiService{Client: client}

	req := jsapi.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(mchID),
	}

	resp, _, err := svc.QueryOrderByOutTradeNo(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
