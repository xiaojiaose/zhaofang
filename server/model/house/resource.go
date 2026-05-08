package house

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
	"gorm.io/gorm"
	"time"
)

type Resource struct {
	ID              uint                 `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt       time.Time            // 创建时间
	UpdatedAt       time.Time            // 更新时间
	DeletedAt       gorm.DeletedAt       `gorm:"index" json:"-"`                                                     // 删除时间
	City            string               `json:"city"`                                                               // 所属城市
	Region          string               `json:"region"`                                                             // 行政区
	Districts       string               `json:"districts"`                                                            // 所属商圈
	DistrictIds     string               `json:"district_ids"`                                                         // 所属商圈ID串
	XiaoquId        uint                 `json:"xiaoqu_id" gorm:"index:idx_owner_xq_building_unit_house_room,priority:2;index:idx_owner_xq_building_unit_house,priority:2;index:idx_owner_xq,priority:2"` // 所属小区ID
	Xiaoqu          string               `json:"xiaoqu"`                                                             // 所属小区名称
	HouseType       string               `json:"house_type,omitempty"`                                               // 房屋类型
	RentType        string               `json:"rent_type,omitempty"`                                                 // 出租类型
	RoomCode        string               `json:"room_code" gorm:"index:idx_owner_xq_building_unit_house_room,priority:6"` // 房间号
	BuildingId      string               `json:"building_id,omitempty" gorm:"index:idx_owner_xq_building_unit_house_room,priority:3;index:idx_owner_xq_building_unit_house,priority:3"` // 楼栋ID
	UnitId          string               `json:"unit_id,omitempty" gorm:"index:idx_owner_xq_building_unit_house_room,priority:4;index:idx_owner_xq_building_unit_house,priority:4"`     // 单元ID
	HouseId         string               `json:"house_id,omitempty" gorm:"index:idx_owner_xq_building_unit_house_room,priority:5;index:idx_owner_xq_building_unit_house,priority:5"`   // 房号ID
	DoorNo          string               `json:"door_no,omitempty"`                                                   // 门牌号
	Floor           string               `json:"floor"`                                                               // 楼层
	RoomNumber      int                  `json:"room_number"`                                                         // 房间数量
	Area            string               `json:"area"`                                                                // 房源面积
	Price           int                  `json:"price"`                                                               // 房源价格
	CommissionPrice int                  `json:"commission_price"`                                                    // 返佣金额
	Feature         string               `json:"feature"`                                                             // 房源特色
	Remarks         string               `json:"remarks"`                                                             // 备注信息
	Attachments     common.AttachmentMap `json:"attachments" gorm:"TYPE:json"`                                        // 房源图片
	HasPic          bool                 `json:"hasPic"`                                                             // 是否有照片
	IsTeamHouse     bool                 `json:"is_team_house" gorm:"default:0;comment:当前状态"`                   // 是否团队房源
	Owner           uint                 `json:"owner" gorm:"index:idx_owner_xq_building_unit_house_room,priority:1;index:idx_owner_xq_building_unit_house,priority:1;index:idx_owner_xq,priority:1;index:idx_owner_status,priority:1"` // 房源归属用户ID
	Status          string               `json:"status" gorm:"index:idx_owner_status,priority:2"`                     // 状态
	ApprovalStatus  string               `json:"approval_status"`                                                    // 审批状态
	Phone           string               `json:"phone"`                                                              // 联系手机号
	UpdatedLastAt   time.Time            `json:"updated_last_at"`                                                   // 最后编辑时间
	Follow          int                  `json:"follow"`                                                             // 关注次数
	View            int                  `json:"view"`                                                               // 浏览次数
	Shared          int                  `json:"shared"`                                                             // 分享次数
	Click           int                  `json:"click"`                                                              // 电话获取次数
	//Saler        string `json:"saler"`        // 销售
	//Designer     string `json:"designer"`     // 设计师
	//LeaseEndDate string `json:"leaseEndDate"` // 截止日期
}

func (land Resource) TableName() string {
	return "house_resources"
}
