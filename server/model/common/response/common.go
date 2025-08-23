package response

import (
	"time"
)

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"` // 总数
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

type ReturnList struct {
	List interface{} `json:"list"`
}

type ReturnValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Area struct {
	Name      string      `json:"name"` // 行政区name
	Id        int         `json:"id"`   // 行政区id
	Sort      int         `json:"sort"` // 排序
	Districts []Districts // 行政区下面的商圈
}

type Districts struct {
	Name string `json:"name"` // 商圈name
	Id   uint   `json:"id"`   // 商圈id
}

type XiaoquRsp struct {
	Id   int    `json:"id"`   // 小区id
	Name string `json:"name"` // 小区名
	Num  int    `json:"num"`  // 房源数量
}

type Resource struct {
	ID             uint   `json:"id"`                   // 主键ID
	XiaoquId       uint   `json:"xiaoqu_id"  `          // 所属小区id
	Xiaoqu         string `json:"xiaoqu"`               // 所属小区名字
	HouseType      string `json:"house_type,omitempty"` // 户型
	RentType       string `json:"rent_type,omitempty"`  // 出租类型
	DoorNo         string `json:"door_no,omitempty"`    //  门牌号
	Floor          string `json:"floor"`                // 楼层
	RoomNumber     int    `json:"room_number" `         // 房间数量
	Area           string `json:"area" `                // 房源面积
	Price          int    `json:"price"`                // 房源价格
	Feature        string `json:"feature" `             // 房源特色
	Remarks        string `json:"remarks"`              // 备注信息
	Owner          uint   `json:"owner"`                // 业主
	Status         string `json:"status"`               // 状态 已出租，已下架，待出租
	ApprovalStatus string `json:"approval_status"`      // 审批状态： 通过 未通过 待审批

	UpdatedLastAt time.Time `json:"updated_last_at" ` // 最后编辑时间
	Follow        int       `json:"follow"`           // 关注次数
	View          int       `json:"view"`             // 浏览次数
	Click         int       `json:"click"`            // 电话获取次数
	Phone         string    `json:"phone"`            // 销售
	WxID          string    `json:"wx_id"`            // 微信号
	WxNickName    string    `json:"wx_nick_name"`     // 微信昵称
}
