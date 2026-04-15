package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/house"

type RewardRecentContact struct {
	ResourceID          uint   `json:"resourceId"`
	Xiaoqu              string `json:"xiaoqu"`
	DoorNo              string `json:"doorNo"`
	PublisherUserID     uint   `json:"publisherUserId"`
	PublisherPhone      string `json:"publisherPhone"`
	PublisherWxNo       string `json:"publisherWxNo"`
	PublisherWxNickName string `json:"publisherWxNickName"`
	LastContactAt       string `json:"lastContactAt"`
}

type RewardApplicationResponse struct {
	house.RewardApplication
	Xiaoqu string `json:"xiaoqu"`
	DoorNo string `json:"doorNo"`
}
