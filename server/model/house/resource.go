package house

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
	"gorm.io/gorm"
	"time"
)

type Resource struct {
	ID          uint                 `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt   time.Time            // 创建时间
	UpdatedAt   time.Time            // 更新时间
	DeletedAt   gorm.DeletedAt       `gorm:"index" json:"-"`               // 删除时间
	City        string               `json:"city" `                        // 所属城市
	Districts   string               `json:"districts" `                   // 所属商圈s
	DistrictIds string               `json:"district_ids" `                // 所属商圈s
	XiaoquId    uint                 `json:"xiaoqu_id"  `                  // 所属小区id
	Xiaoqu      string               `json:"xiaoqu"`                       // 所属小区名字
	HouseType   string               `json:"house_type,omitempty"`         // 户型
	RentType    string               `json:"rent_type,omitempty"`          // 出租类型
	DoorNo      string               `json:"door_no,omitempty"`            //  门牌号
	Floor       string               `json:"floor"`                        // 楼层
	RoomNumber  int                  `json:"room_number" `                 // 房间数量
	Area        string               `json:"area" `                        // 房源面积
	Price       int                  `json:"price"`                        // 房源价格
	Feature     string               `json:"feature" `                     // 房源特色
	Remarks     string               `json:"remarks"`                      // 备注信息
	Attachments common.AttachmentMap `json:"attachments" gorm:"TYPE:json"` // 公寓照片
	Owner       uint                 `json:"owner"`                        // 业主
	Status      string               `json:"status"`                       // 状态 已出租，已下架，待出租
	//Saler        string `json:"saler"`        // 销售
	//Designer     string `json:"designer"`     // 设计师
	//LeaseEndDate string `json:"leaseEndDate"` // 截止日期
}

func (land Resource) TableName() string {
	return "house_resources"
}
