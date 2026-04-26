package response

type XiaoQuDictTreeResponse struct {
	XiaoquID    uint                     `json:"xiaoquId"`    // 小区ID
	CommunityID int                      `json:"communityId"` // 社区ID
	Buildings   []XiaoQuDictBuildingNode `json:"buildings"`   // 楼栋树
}

type XiaoQuDictBuildingNode struct {
	ID             uint                 `json:"id"`             // 楼栋记录ID
	BuildingOpenID string               `json:"buildingOpenId"` // 楼栋OpenID
	Name           string               `json:"name"`           // 楼栋名称
	Units          []XiaoQuDictUnitNode `json:"units"`          // 单元列表
}

type XiaoQuDictUnitNode struct {
	ID         uint                  `json:"id"`         // 单元记录ID
	UnitOpenID string                `json:"unitOpenId"` // 单元OpenID
	Name       string                `json:"name"`       // 单元名称
	Houses     []XiaoQuDictHouseNode `json:"houses"`     // 房号列表
}

type XiaoQuDictHouseNode struct {
	ID          uint   `json:"id"`          // 房号记录ID
	HouseOpenID string `json:"houseOpenId"` // 房号OpenID
	Name        string `json:"name"`        // 房号名称
}
