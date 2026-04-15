package house

import (
	"time"

	"gorm.io/gorm"
)

type BatchUploadRecord struct {
	ID            uint `gorm:"primarykey" json:"ID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	UserID        uint           `json:"userId" gorm:"index;comment:上传用户ID"`
	FileName      string         `json:"fileName" gorm:"comment:文件名"`
	BatchNo       string         `json:"batchNo" gorm:"index;comment:批次号"`
	TotalCount    int            `json:"totalCount" gorm:"comment:总条数"`
	SuccessCount  int            `json:"successCount" gorm:"comment:成功条数"`
	FailedCount   int            `json:"failedCount" gorm:"comment:失败条数"`
	ResultSummary string         `json:"resultSummary" gorm:"type:text;comment:结果摘要"`
}

func (BatchUploadRecord) TableName() string {
	return "house_batch_upload_records"
}
