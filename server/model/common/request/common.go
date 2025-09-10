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
	Start time.Time `form:"start" binding:"required" time_format:"2006-01-02T15:04:05Z"`
	End   time.Time `form:"end" binding:"required" time_format:"2006-01-02T15:04:05Z"`
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
type SearchDistrict struct {
	AreaId uint `json:"areaId"` // 区域id
}

type Empty struct{}
type SearchXiaoqu struct {
	PageInfo
	CityId    string `json:"cityId"`    // 城市id
	AreaId    string `json:"areaId"`    // 区id
	Districts []int  `json:"districts"` //
}
type SearchResource struct {
	PageInfo
	XiaoquId uint   `json:"xiaoquId"` // 小区id
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type MySearchResource struct {
	PageInfo
	DoorNo   string `json:"doorNo"`   // 户室号
	Status   string `json:"status"`   // 出租类型： 已出租，已下架，待出租
	XiaoquId uint   `json:"xiaoquId"` // 小区id
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type SearchNameResource struct {
	PageInfo
	XiaoquId       uint   `json:"xiaoquId"` // 小区id
	OrderKey       string `json:"orderKey"` // 排序
	Desc           bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
	Phone          string `json:"phone"`
	ApprovalStatus string `json:"approvalStatus"` //  通过 未通过 待审批
}

type RentType struct {
	Name      string   `json:"name"` // 房型 整租、分整租、合租
	HouseType []string `json:"houseType"`
	Feature   []string `json:"feature"`
}

type SearchHouseResource struct {
	PageInfo
	XiaoquId       uint      `json:"xiaoquId"`       // 小区id
	RentType       string    `json:"rentType"`       // 出租类型
	DoorNo         string    `json:"doorNo"`         // 户室号
	OrderKey       string    `json:"orderKey"`       // 排序
	Desc           bool      `json:"desc"`           // 排序方式:升序false(默认)|降序true
	Phone          string    `json:"phone"`          // 手机号
	WxNo           string    `json:"wxNo"`           // 微信号
	ApprovalStatus string    `json:"approvalStatus"` //  审核状态： 通过 未通过 待审批
	HasPic         string    `json:"hasPic"`         // 是否有图片
	UpdatedAtLast  time.Time `json:"updatedAtLast"`  // 开始时间
	UpdatedAtStart time.Time `json:"updatedAtStart"` // 结束时间
}

type SearchOther struct {
	Phone          string    `json:"phone"`
	HasPic         string    `json:"hasPic"`         // 是否有图片
	UpdatedAtLast  time.Time `json:"updatedAtLast"`  // 开始时间
	UpdatedAtStart time.Time `json:"updatedAtStart"` // 结束时间
	RentType       string    `json:"rentType"`       // 出租类型
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
	XiaoquId  []int  `json:"xiaoquIds"` // 商圈ids
	HouseType string `json:"houseType"` // 1居室、2居室
	RentType  string `json:"rentType"`  // 整租、合租、分整租
	Price     int    `json:"price"`     // 价格 1580
	Feature   string `json:"feature"`   // 有无电梯

	Page     int `json:"page" form:"page"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页大小
}
