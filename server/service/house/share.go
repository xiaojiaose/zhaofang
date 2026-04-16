package house

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
)

type ShareService struct{}

func (service *ShareService) Create(userID uint, expireDays int) (*house.ShareToken, error) {
	// 分享链接默认 7 天有效，前端不传时走默认值。
	if expireDays <= 0 {
		expireDays = 7
	}
	token := randomToken(24)
	entity := &house.ShareToken{
		UserID:       userID,
		Token:        token,
		Status:       "active",
		ExpireAtUnix: time.Now().Add(time.Duration(expireDays) * 24 * time.Hour).Unix(),
	}
	return entity, global.GVA_DB.Create(entity).Error
}

func (service *ShareService) GetByToken(token string) (*house.ShareToken, error) {
	// 分享页是未登录可访问的，所以这里只校验 token 状态和过期时间。
	var entity house.ShareToken
	if err := global.GVA_DB.Where("token = ? AND status = ?", token, "active").First(&entity).Error; err != nil {
		return nil, err
	}
	if entity.ExpireAtUnix < time.Now().Unix() {
		return nil, errors.New("分享链接已失效")
	}
	return &entity, nil
}

func randomToken(n int) string {
	// 正常情况下走随机字节；
	// 极端情况下随机源不可用时，退化成时间戳字符串，避免整个分享流程失败。
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(buf)
}
