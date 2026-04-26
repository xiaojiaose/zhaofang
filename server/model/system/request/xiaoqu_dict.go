package request

type XiaoQuDictUpsertRequest struct {
	XiaoquID    uint                   `json:"xiaoquId"`    // 小区ID
	BuildingOps []XiaoQuDictBuildingOp `json:"buildingOps"` // 楼栋新增/修改列表
	UnitOps     []XiaoQuDictUnitOp     `json:"unitOps"`     // 单元新增/修改列表
	HouseOps    []XiaoQuDictHouseOp    `json:"houseOps"`    // 房号新增/修改列表
}

type XiaoQuDictDeleteRequest struct {
	XiaoquID    uint   `json:"xiaoquId"`    // 小区ID
	BuildingIDs []uint `json:"buildingIds"` // 要删除的楼栋ID列表
	UnitIDs     []uint `json:"unitIds"`     // 要删除的单元ID列表
	HouseIDs    []uint `json:"houseIds"`    // 要删除的房号ID列表
}

type XiaoQuDictBuildingOp struct {
	ID             uint   `json:"id"`             // 楼栋记录ID（更新必填）
	BuildingOpenID string `json:"buildingOpenId"` // 楼栋OpenID（新增可选）
	Name           string `json:"name"`           // 楼栋名称（必填）
}

type XiaoQuDictUnitOp struct {
	ID             uint   `json:"id"`             // 单元记录ID（更新必填）
	BuildingOpenID string `json:"buildingOpenId"` // 所属楼栋OpenID（新增必填）
	UnitOpenID     string `json:"unitOpenId"`     // 单元OpenID（新增可选）
	Name           string `json:"name"`           // 单元名称（必填）
}

type XiaoQuDictHouseOp struct {
	ID          uint   `json:"id"`          // 房号记录ID（更新必填）
	UnitOpenID  string `json:"unitOpenId"`  // 所属单元OpenID（新增必填）
	HouseOpenID string `json:"houseOpenId"` // 房号OpenID（新增可选）
	Name        string `json:"name"`        // 房号名称（必填）
}
