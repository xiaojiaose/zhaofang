package response

type LandlordContactViewResponse struct {
	ID              uint   `json:"id"`              // 记录ID
	ResourceID      uint   `json:"resourceId"`      // 房源ID
	Xiaoqu          string `json:"xiaoqu"`          // 小区
	DoorNo          string `json:"doorNo"`          // 门牌号
	PublisherUserID uint   `json:"publisherUserId"` // 发布人ID
	PublisherName   string `json:"publisherName"`   // 发布人名称
	PublisherPhone  string `json:"publisherPhone"`  // 发布人手机号
	ResourceAt      string `json:"resourceAt"`      // 房源录入时间
	ViewerUserID    uint   `json:"viewerUserId"`    // 查看人ID
	ViewerName      string `json:"viewerName"`      // 查看人名称
	ViewerPhone     string `json:"viewerPhone"`     // 查看人手机号
	ViewAt          string `json:"viewAt"`          // 查看时间
	Status          string `json:"status"`          // 当前状态
	LastOperateAt   string `json:"lastOperateAt"`   // 最后操作时间
}
