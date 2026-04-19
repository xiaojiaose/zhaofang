package house

import (
	"time"

	"gorm.io/gorm"
)

type BatchUploadRecord struct {
	ID            uint           `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt     time.Time      // 创建时间
	UpdatedAt     time.Time      // 更新时间
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`                              // 软删除时间
	UserID        uint           `json:"userId" gorm:"index;comment:上传用户ID"`          // 上传用户ID
	FileName      string         `json:"fileName" gorm:"comment:文件名"`                 // 上传文件名
	BatchNo       string         `json:"batchNo" gorm:"index;comment:批次号"`            // 导入批次号
	Remark        string         `json:"remark" gorm:"type:text;comment:批量上传备注"`      // 批量上传备注
	TotalCount    int            `json:"totalCount" gorm:"comment:总条数"`               // Excel 总条数
	SuccessCount  int            `json:"successCount" gorm:"comment:成功条数"`            // 导入成功条数
	FailedCount   int            `json:"failedCount" gorm:"comment:失败条数"`             // 导入失败条数
	ResultSummary string         `json:"resultSummary" gorm:"type:text;comment:结果摘要"` // 导入结果摘要
}

func (BatchUploadRecord) TableName() string {
	return "house_batch_upload_records"
}
