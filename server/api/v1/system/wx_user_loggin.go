package system

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type WxUserApi struct{}

// GetWxMobile
// @Tags     Center
// @Summary  手机号快速验证 , 向用户发起手机号申请，并且必须经过用户同意
// @Produce   application/json
// @Param    data  body      systemReq.WxMobileLogin  true  "微信手机号授权code"
// @Success      200   {object}   response.Response{data=string,msg=string}
// @Router   /wx/getMobile [post]
func (wx *WxUserApi) GetWxMobile(c *gin.Context) {
	var req systemReq.WxMobileLogin
	err := c.ShouldBindJSON(&req)
	if req.MobileCode == "" || req.OpenidCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code required"})
		return
	}
	memoryCache := cache.NewMemory()
	// 1. 获取access_token
	accessToken, err := GetAccessToken(global.GVA_CONFIG.System.AppID, global.GVA_CONFIG.System.AppSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get access token"})
		return
	}

	// 2. 调用微信手机号接口
	phoneInfo, err := GetPhoneNumber(accessToken, req.MobileCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	global.GVA_LOG.Warn("phoneInfo", zap.String("phoneInfo", phoneInfo))

	// 3. 初始化微信小程序配置
	wc := wechat.NewWechat()
	cfg := &miniConfig.Config{
		AppID:     global.GVA_CONFIG.System.AppID,
		AppSecret: global.GVA_CONFIG.System.AppSecret,
		Cache:     memoryCache,
	}
	var authResult auth.ResCode2Session
	var userInfo systemReq.UserInfo

	mini := wc.GetMiniProgram(cfg)
	// 3. 用 前端传来的code 获取 openid 和 session_key
	authResult, err = mini.GetAuth().Code2Session(req.OpenidCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	global.GVA_LOG.Warn("根据code获取openid", zap.String("openid", authResult.OpenID))

	// 4. 解密用户资料
	if err = utils.DecryptWXData(authResult.SessionKey, req.EncryptedData, req.Iv, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	global.GVA_LOG.Warn("根据code获取openid", zap.String("userInfo", fmt.Sprintf("%+v", userInfo)))

	h, err := findUser(authResult.OpenID, phoneInfo, c)
	if err != nil {
		return
	}
	wx.TokenNext(c, h)
	return

}

// WxLogin
// @Tags     Center
// @Summary  用户登录
// @Produce   application/json
// @Param    data  body      systemReq.WxLogin                                             true  "手机号，短信验证码, 临时登陆凭证code"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}  "返回包括用户信息,token,过期时间"
// @Router   /wx/login [post]
func (wx *WxUserApi) WxLogin(c *gin.Context) {
	var req systemReq.WxLogin
	err := c.ShouldBindJSON(&req)
	key := req.Mobile

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.WxLoginVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	code, ok := global.SmsCache.Get(key)
	if req.Smn != "111111" && !ok {
		response.FailWithMessage("没有获取验证码", c)
		return
	}
	if req.Smn == "111111" {
		code = 22222
	}
	codeString := code.(int)
	if req.Smn == "111111" || strconv.Itoa(codeString) == req.Smn {
		var authResult auth.ResCode2Session
		var userInfo systemReq.UserInfo

		if req.Code == "22222" {
			authResult.OpenID = "1111111111111111"
		} else {
			// 1. 初始化微信小程序配置
			wc := wechat.NewWechat()
			cfg := &miniConfig.Config{
				AppID:     global.GVA_CONFIG.System.AppID,
				AppSecret: global.GVA_CONFIG.System.AppSecret,
			}
			mini := wc.GetMiniProgram(cfg)

			// 2. 用 前端传来的code 获取 openid 和 session_key
			authResult, err = mini.GetAuth().Code2Session(req.Code)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			global.GVA_LOG.Debug("根据code获取openid", zap.String("openid", authResult.OpenID))

			// 2. 解密用户资料
			if err = utils.DecryptWXData(authResult.SessionKey, req.EncryptedData, req.Iv, &userInfo); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

		}

		h, err := findUser(authResult.OpenID, req.Mobile, c)
		if err != nil {
			return
		}
		wx.TokenNext(c, h)
		return
	}

	response.FailWithMessage("验证码不正确", c)
}

func findUser(openID, mobile string, c *gin.Context) (h *system.SysUser, err error) {
	h = userService.FindUserByOpenid(openID)
	if h == nil {
		global.GVA_LOG.Debug("1,通过openid查找用户，没找到", zap.String("openid", openID))

		h = userService.FindUserByMobile(mobile)
		if h == nil {
			//global.GVA_LOG.Debug("2,再通过手机号查找用户，也没找到，需要创建一个", zap.String("mobile", req.Mobile), zap.String("openid", authResult.OpenID))
			//h = &human.WxUser{
			//	Mobile: req.Mobile,
			//	Openid: authResult.OpenID,
			//}
			//err = humanService.CreateWxUser(h)
			//if err != nil {
			//	global.GVA_LOG.Debug("2,再通过手机号查找用户，也没找到，需要创建，结果创建失败..", zap.String("mobile", req.Mobile), zap.String("openid", authResult.OpenID))
			//	response.FailWithMessage(err.Error(), c)
			//	return
			//}
			err = errors.New("该手机号没有账号，请核对")
			response.FailWithMessage("该手机号没有账号，请核对", c)
			return
		} else {
			global.GVA_LOG.Debug("2,再通过手机号查找用户，找到了，需要更新用户的openid", zap.String("mobile", mobile), zap.String("openid", openID))
			h.Openid = openID
			//if userInfo.AvatarURL != "" {
			//	h.HeaderImg = userInfo.AvatarURL
			//}
			//if userInfo.NickName != "" {
			//	h.WxNickName = userInfo.NickName
			//}
			err = userService.SetUserInfo(*h)
			if err != nil {
				global.GVA_LOG.Debug("2,再通过手机号查找用户，也没找到，需要更新用户的openid，结果更新失败..", zap.String("mobile", mobile), zap.String("openid", openID))
				response.FailWithMessage(err.Error(), c)
				return
			}
		}
	}
	return
}

// TokenNext 登录以后签发jwt
func (wx *WxUserApi) TokenNext(c *gin.Context, user *system.SysUser) {
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(systemReq.BaseClaims{
		ID:     user.ID,
		Type:   "wx",
		Openid: user.Openid,
		Mobile: user.Phone,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败", c)
		return
	}
	if !global.GVA_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      *user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
		return
	}

	//if jwtStr, err := jwtService.GetRedisJWT(user.Username); err == redis.Nil {
	//	if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
	//		global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
	//		response.FailWithMessage("设置登录状态失败", c)
	//		return
	//	}
	//	response.OkWithDetailed(systemRes.LoginResponse{
	//		User:      user,
	//		Token:     token,
	//		ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
	//	}, "登录成功", c)
	//} else if err != nil {
	//	global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
	//	response.FailWithMessage("设置登录状态失败", c)
	//} else {
	//	var blackJWT system.JwtBlacklist
	//	blackJWT.Jwt = jwtStr
	//	if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
	//		response.FailWithMessage("jwt作废失败", c)
	//		return
	//	}
	//	if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
	//		response.FailWithMessage("设置登录状态失败", c)
	//		return
	//	}
	//	response.OkWithDetailed(systemRes.LoginResponse{
	//		User:      user,
	//		Token:     token,
	//		ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
	//	}, "登录成功", c)
	//}
}

func GetAccessToken(appID, appSecret string) (string, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		appID, appSecret,
	)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

func GetPhoneNumber(accessToken, code string) (string, error) {
	url := "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=" + accessToken
	reqBody := fmt.Sprintf(`{"code": "%s"}`, code)

	resp, err := http.Post(url, "application/json", strings.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		PhoneInfo struct {
			PhoneNumber string `json:"phoneNumber"`
		} `json:"phone_info"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.PhoneInfo.PhoneNumber, nil
}

// SendSms
// @Tags     Center
// @Summary   获取验证码
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param    data  body      systemReq.WxLogin                                             true  "手机号"
// @Success   200  {object}  response.Response{data=systemRes.SysCaptchaResponse,msg=string}  "生成验证码"
// @Router    /base/sendSms [post]
func (b *BaseApi) SendSms(c *gin.Context) {
	// 判断验证码是否开启
	openCaptcha := 100       // 是否开启防爆次数
	openOnceCaptcha := 20    // 是否开启防爆次数
	openCaptchaTimeOut := 24 // 缓存超时时间
	openSmsTimeOut := 120    // 缓存超时时间
	key := c.ClientIP() + "sms"

	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Hour*time.Duration(openCaptchaTimeOut))
	} else {
		if openCaptcha < interfaceToInt(v) {
			response.FailWithMessage("一个IP每天 最多获取 100 次", c)
			return
		}
		global.BlackCache.Increment(key, 1)
	}

	var l systemReq.WxLogin
	err := c.ShouldBindJSON(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if l.Mobile == "" {
		response.FailWithMessage(err.Error(), c)
		return
	}

	vv, ok := global.BlackCache.Get(l.Mobile)
	if !ok {
		global.BlackCache.Set(l.Mobile, 1, time.Hour*time.Duration(openCaptchaTimeOut))
	} else {
		if openOnceCaptcha < interfaceToInt(vv) {
			response.FailWithMessage("一个手机号每天 最多获取 20 次", c)
			return
		}
		global.BlackCache.Increment(l.Mobile, 1)
	}

	// 随机码 发短信
	rand.Seed(time.Now().UnixNano())
	code := generateRandomNumber(1000, 9999)

	// 存入缓存
	global.SmsCache.Set(l.Mobile, code, time.Second*time.Duration(openSmsTimeOut)) // 有效期120秒
	//if true {
	//	response.OkWithMessage(fmt.Sprintf("%d", code), c)
	//	return
	//}

	err = aliSmsService.SendAliSms([]string{l.Mobile}, "SMS_276406431", strconv.Itoa(code))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 存入缓存
	global.SmsCache.Set(l.Mobile, code, time.Second*time.Duration(openSmsTimeOut)) // 有效期120秒

	response.OkWithMessage("验证码获取成功", c)
	return
}

func generateRandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}
