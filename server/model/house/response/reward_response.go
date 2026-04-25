package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/house"

type RewardRecentContact struct {
	ResourceID          uint   `json:"resourceId"`          // 房源ID
	Xiaoqu              string `json:"xiaoqu"`              // 小区名称
	DoorNo              string `json:"doorNo"`              // 户室号
	PublisherUserID     uint   `json:"publisherUserId"`     // 发布人ID
	PublisherPhone      string `json:"publisherPhone"`      // 发布人手机号
	PublisherHeaderImg  string `json:"publisherHeaderImg"`  // 发布人头像
	PublisherWxNo       string `json:"publisherWxNo"`       // 发布人微信号
	PublisherWxNickName string `json:"publisherWxNickName"` // 发布人微信昵称
	Status              string `json:"status"`              // 房源状态/审核状态
	LastContactAt       string `json:"lastContactAt"`       // 最近联系时间
}

type RewardApplicationResponse struct {
	house.RewardApplication
	Xiaoqu              string `json:"xiaoqu"`              // 小区名称
	DoorNo              string `json:"doorNo"`              // 户室号
	ApplyUserHeaderImg  string `json:"applyUserHeaderImg"`  // 申请人头像
	ApplyUserWxNickName string `json:"applyUserWxNickName"` // 申请人微信昵称
}
