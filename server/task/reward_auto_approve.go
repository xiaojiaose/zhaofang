package task

import (
	"time"

	"gorm.io/gorm"
)

func RewardAutoApprove(db *gorm.DB) error {
	cutoff := time.Now().Add(-48 * time.Hour).UnixMilli()
	return db.Table("house_reward_applications").
		Where("publisher_confirm_status = ? AND audit_status IN ?", "已确认", []string{"待审核", "审核中"}).
		Where("last_operated_at_unix_milli < ?", cutoff).
		Update("audit_status", "审核通过待发放").Error
}
