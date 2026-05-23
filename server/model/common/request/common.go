package request

import (
	"gorm.io/gorm"
	"time"
)

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`   // 关键字
}

func (r *PageInfo) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.Page <= 0 {
			r.Page = 1
		}
		switch {
		case r.PageSize > 100:
			r.PageSize = 100
		case r.PageSize <= 0:
			r.PageSize = 10
		}
		offset := (r.Page - 1) * r.PageSize
		return db.Offset(offset).Limit(r.PageSize)
	}
}

// GetById Find by id structure
type GetById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

type GetByIdStr struct {
	ID string `json:"id" form:"id"` // 主键ID
}

type GetStatis struct {
	Start string `form:"start" binding:"required"` // 开始时间
	End   string `form:"end" binding:"required"`   // 结束时间
	Phone string `form:"phone"`                    // 手机号筛选
}

type VisitReq struct {
	Phone string `json:"phone" from:"phone"` // 手机号
	WxNo  string `json:"wxNo" from:"wxNo"`   // 微信号
	PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type VisitResponse struct {
	Phone string    `json:"phone" from:"phone"` // 手机号
	WxNo  string    `json:"wxNo" from:"wxNo"`   // 微信号
	Date  time.Time `json:"date" from:"date"`   // 访问时间
}

func (r *GetById) Uint() uint {
	return uint(r.ID)
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}

// GetAuthorityId Get role by id structure
type GetAuthorityId struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
}

type SearchArea struct {
	CityId string `json:"cityId"` // 城市id
}

type WxProfileSync struct {
	WxNickName string `json:"wxNickName"` // 微信昵称
	HeaderImg  string `json:"headerImg"`  // 微信头像
	WxNo       string `json:"wxNo"`       // 微信号
}

type SearchDistrict struct {
	AreaId uint `json:"areaId"` // 区域id
}

type Empty struct{}
type SearchXiaoqu struct {
	PageInfo
	CityId    string `json:"cityId"`    // 城市id
	AreaId    string `json:"areaId"`    // 区id
	Districts []int  `json:"districts"` //
	Able      bool   `json:"able"`
}
type SearchResource struct {
	PageInfo
	XiaoquId uint   `json:"xiaoquId"` // 小区ID
	OrderKey string `json:"orderKey"` // 排序字段
	Desc     bool   `json:"desc"`     // 排序方式
}

type MySearchResource struct {
	PageInfo
	DoorNo   string `json:"doorNo"`   // 户室号
	Status   string `json:"status"`   // 房源状态
	XiaoquId uint   `json:"xiaoquId"` // 小区ID
	OrderKey string `json:"orderKey"` // 排序字段
	Desc     bool   `json:"desc"`     // 排序方式
}

type SearchNameResource struct {
	PageInfo
	XiaoquId       uint      `json:"xiaoquId"`       // 小区ID
	OrderKey       string    `json:"orderKey"`       // 排序字段
	Desc           bool      `json:"desc"`           // 排序方式
	Phone          string    `json:"phone"`          // 手机号
	ApprovalStatus string    `json:"approvalStatus"` // 审批状态
	Status         string    `json:"status"`         // 房源状态：待出租、已下架
	UpdatedAtStart time.Time `json:"updatedAtStart"` // 最后编辑时间开始
	UpdatedAtLast  time.Time `json:"updatedAtLast"`  // 最后编辑时间结束
}

type RentType struct {
	Name      string   `json:"name"`      // 房型名称
	HouseType []string `json:"houseType"` // 适用房型
	Feature   []string `json:"feature"`   // 可选亮点
}

type SearchHouseResource struct {
	PageInfo
	XiaoquId       uint      `json:"xiaoquId"`       // 小区ID
	RentType       string    `json:"rentType"`       // 出租类型
	DoorNo         string    `json:"doorNo"`         // 户室号
	OrderKey       string    `json:"orderKey"`       // 排序字段
	Desc           bool      `json:"desc"`           // 排序方式
	Phone          string    `json:"phone"`          // 手机号
	WxNo           string    `json:"wxNo"`           // 微信号
	ApprovalStatus string    `json:"approvalStatus"` // 审核状态
	HasPic         string    `json:"hasPic"`         // 是否有图片
	UpdatedAtLast  time.Time `json:"updatedAtLast"`  // 开始时间
	UpdatedAtStart time.Time `json:"updatedAtStart"` // 结束时间
}

type SearchOther struct {
	Phone         string    `json:"phone"`          // 手机号
	HasPic        string    `json:"hasPic"`         // 是否有图片
	UpdatedAtStart time.Time `json:"updatedAtStart"` // 开始时间
	UpdatedAtLast  time.Time `json:"updatedAtLast"`  // 结束时间
	RentType      string    `json:"rentType"`       // 出租类型
	IsTeamHouse   string    `json:"isTeamHouse"`    // 是否团队房源
	HouseType     string    `json:"houseType"`      // 房屋类型
	HasCommission string    `json:"hasCommission"`  // 是否有返佣
	Status        string    `json:"status"`         // 房源状态：待出租、已下架
}

type HouseStateReq struct {
	Ids   []uint // 房源id
	State int    // 1 上架（通过），2下架（不通过）
}

type FavoriteSearch struct {
	PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type ResourceSearch struct {
	XiaoquId    []int  `json:"xiaoquIds"`   // 商圈ids
	HouseType   string `json:"houseType"`   // 1居室、2居室
	RentType    string `json:"rentType"`    // 整租、合租、分整租
	Price       int    `json:"price"`       // 价格 1580
	Feature     string `json:"feature"`     // 有无电梯
	HouseSource string `json:"houseSource"` // commission \ landlord \ team
	Page        int    `json:"page"`        // 页码
	PageSize    int    `json:"pageSize"`    // 每页大小
}

type ResourceShareSearch struct {
	Token string `json:"token"` // 分享token
	ResourceSearch
	Page     int `json:"page"`         // 页码
	PageSize int `json:"pageSize"`     // 每页大小
}

type ResourceShareListSearch struct {
	Token    string `json:"token"`    // 分享token
	XiaoquID uint   `json:"xiaoquId"` // 小区ID
	ResourceSearch
	Page     int `json:"page"`         // 页码
	PageSize int `json:"pageSize"`     // 每页大小
}

type RewardApplicationCreate struct {
	ResourceID uint   `json:"resourceId"` // 房源ID
	Remark     string `json:"remark"`     // 申请备注
}

type RewardApplicationSearch struct {
	PageInfo
	PublisherConfirmStatus string `json:"publisherConfirmStatus"` // 发布人确认状态
	AuditStatus            string `json:"auditStatus"`            // 后台审核状态
}

type RewardApplicationAction struct {
	ID     uint   `json:"id"`     // 申请单ID
	Action string `json:"action"` // 操作类型
}

type ContactQuotaCreate struct {
	UserPhone string `json:"userPhone"` // 经纪人手机号
	Amount    int    `json:"amount"`    // 增加次数
	Remark    string `json:"remark"`    // 备注
}

type ContactQuotaSearch struct {
	PageInfo
	UserPhone string `json:"userPhone"` // 经纪人手机号
}

type LandlordContactViewSearch struct {
	PageInfo
	Status         string `json:"status"`         // 当前状态
	PublisherPhone string `json:"publisherPhone"` // 发布人手机号
	ViewerPhone    string `json:"viewerPhone"`    // 查看人手机号
}

type LandlordContactViewAction struct {
	ID     uint   `json:"id"`     // 记录ID
	Action string `json:"action"` // 操作类型 pending  processing approved paid rejected   待审核 审核中 已通过，未付款 已付款 未通过，已拒绝
}

type ResourceShareCreate struct {
	ExpireDays int `json:"expireDays"` // 有效天数
}
