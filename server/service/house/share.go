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
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(buf)
}
