# 微信支付对接使用说明

## 配置说明

在 `server/service/system/wx_pay.go` 中修改配置：

```go
var wxPayConfig = system.WxPayConfig{
    AppID:     "your_appid",        // 小程序AppID
    MchID:     "your_mch_id",       // 商户号
    APIKey:    "your_api_key",      // API密钥
    NotifyURL: "https://your-domain.com/api/wxpay/notify", // 支付回调地址
}
```

## API接口

### 1. 创建支付订单
```
POST /api/wxpay/createOrder
```

请求参数：
```json
{
    "openId": "用户openid",
    "body": "商品描述",
    "totalFee": 100  // 支付金额(分)
}
```

响应：
```json
{
    "code": 0,
    "data": {
        "orderNo": "订单号",
        "payParams": {
            "timeStamp": "时间戳",
            "nonceStr": "随机字符串", 
            "package": "prepay_id=xxx",
            "signType": "MD5",
            "paySign": "签名"
        }
    }
}
```

### 2. 查询订单状态
```
GET /api/wxpay/queryOrder?orderNo=订单号
```

### 3. 支付回调（微信调用）
```
POST /api/wxpay/notify
```

## 小程序端调用示例

```javascript
// 1. 调用后端创建订单
wx.request({
    url: 'https://your-domain.com/api/wxpay/createOrder',
    method: 'POST',
    data: {
        openId: 'user_openid',
        body: '商品名称',
        totalFee: 100
    },
    success: (res) => {
        if (res.data.code === 0) {
            // 2. 调用微信支付
            wx.requestPayment({
                ...res.data.data.payParams,
                success: () => {
                    console.log('支付成功');
                },
                fail: () => {
                    console.log('支付失败');
                }
            });
        }
    }
});
```

## 数据库表结构

系统会自动创建 `pay_orders` 表存储订单信息。

## 注意事项

1. 需要在微信商户平台配置支付回调地址
2. 确保服务器支持HTTPS
3. 正式环境需要配置真实的商户信息
4. 建议添加订单超时处理逻辑
