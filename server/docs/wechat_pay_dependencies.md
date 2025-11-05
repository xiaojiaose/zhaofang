# 微信支付SDK依赖安装

## 安装依赖

在项目根目录执行：

```bash
cd server
go get github.com/wechatpay-apiv3/wechatpay-go
```

## 配置说明

在 `server/service/system/wx_pay.go` 中修改以下配置：

```go
var (
    appID     = "wx1234567890123456"                    // 小程序AppID
    mchID     = "1234567890"                            // 商户号
    notifyURL = "https://your-domain.com/api/wxpay/notify" // 支付回调地址
    
    // 商户私钥文件路径 (从微信商户平台下载)
    privateKeyPath = "/path/to/apiclient_key.pem"
    
    // 商户证书序列号 (从微信商户平台获取)
    mchCertificateSerialNumber = "1234567890ABCDEF1234567890ABCDEF12345678"
    
    // APIv3密钥 (在微信商户平台设置)
    apiV3Key = "your_32_char_apiv3_key_here_123456"
)
```

## 证书文件准备

1. 登录微信商户平台
2. 下载API证书 (apiclient_key.pem, apiclient_cert.pem)
3. 将证书文件放到服务器安全目录
4. 更新 `privateKeyPath` 为实际路径

## API接口

### 创建支付订单
```
POST /api/wxpay/createOrder
Content-Type: application/json

{
    "openId": "用户openid",
    "body": "商品描述", 
    "totalFee": 100
}
```

### 支付回调 (微信调用)
```
POST /api/wxpay/notify
```

### 查询本地订单
```
GET /api/wxpay/queryOrder?orderNo=订单号
```

### 查询微信订单
```
GET /api/wxpay/queryWxOrder?outTradeNo=商户订单号
```

## 小程序调用示例

```javascript
// 创建订单并支付
wx.request({
    url: 'https://your-domain.com/api/wxpay/createOrder',
    method: 'POST',
    header: {
        'Authorization': 'Bearer your_token'
    },
    data: {
        openId: 'user_openid',
        body: '商品名称',
        totalFee: 100  // 1元 = 100分
    },
    success: (res) => {
        if (res.data.code === 0) {
            // 调用微信支付
            wx.requestPayment({
                ...res.data.data.payParams,
                success: () => {
                    wx.showToast({title: '支付成功'});
                },
                fail: () => {
                    wx.showToast({title: '支付失败', icon: 'error'});
                }
            });
        }
    }
});
```

## 优势

使用官方SDK的优势：
- 自动处理签名验证
- 支持APIv3最新接口
- 内置证书管理
- 更好的错误处理
- 官方维护，安全可靠
