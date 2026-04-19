package house

import (
	"context"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	response2 "github.com/flipped-aurora/gin-vue-admin/server/model/house/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/linxdeep/linxdeep-framework/pkg/searchx"
)

func (service *ResourceService) SharedMapAgg(userID uint) (list []response2.SharedMapXiaoqu, err error) {
	// 个人房源地图先按发布人做 ES 聚合，再把聚合结果映射到小区坐标。
	// 这样分享页不需要先拉整批房源再自己分组，前端也能直接按小区点位渲染。
	condition := searchx.Condition{
		Terms: []searchx.Term{
			{Field: "owner", Value: strconv.Itoa(int(userID))},
			{Field: "status", Value: "待出租"},
		},
		Aggs: []searchx.Agg{{Field: "xiaoqu_id"}},
		Ors:  []searchx.Condition{{Terms: []searchx.Term{}}},
		Nots: []searchx.Condition{{Terms: []searchx.Term{}}},
	}
	_, _, agg, err := global.Gva_ResourceSearch.SearchAgg(context.Background(), condition, searchx.QueryParams{
		Fields: []string{"xiaoqu_id"},
		Size:   0,
	})
	if err != nil {
		return
	}
	for id, num := range agg["xiaoqu_id"] {
		xID, parseErr := strconv.Atoi(id)
		if parseErr != nil {
			continue
		}
		var xq system.XiaoQu
		qErr := global.GVA_DB.Where("id = ?", xID).First(&xq).Error
		if qErr != nil {
			continue
		}
		list = append(list, response2.SharedMapXiaoqu{
			XiaoquId:  uint(xID),
			Name:      xq.Name,
			Latitude:  xq.Latitude,
			Longitude: xq.Longitude,
			Num:       num,
		})
	}
	return
}

func (service *ResourceService) SharedMapList(userID, xiaoquID uint, info request.PageInfo) (resources []*house.Resource, total int64, err error) {
	// 点击某个小区点位后，仍然通过 ES 按 owner + xiaoqu_id 过滤，
	// 再回表取完整房源数据，避免直接全表扫描。
	condition := searchx.Condition{
		Terms: []searchx.Term{
			{Field: "owner", Value: strconv.Itoa(int(userID))},
			{Field: "status", Value: "待出租"},
		},
		Ors:  []searchx.Condition{{Terms: []searchx.Term{}}},
		Nots: []searchx.Condition{{Terms: []searchx.Term{}}},
	}
	if xiaoquID != 0 {
		condition.Terms = append(condition.Terms, searchx.Term{Field: "xiaoqu_id", Value: strconv.Itoa(int(xiaoquID))})
	}
	houseList, totalCount, err := global.Gva_ResourceSearch.Search(context.Background(), condition, searchx.QueryParams{
		Fields: []string{"house_id", "xiaoqu_id", "xiaoqu", "_id"},
		Size:   info.PageSize,
		Page:   info.Page - 1,
		Asc:    false,
		By:     "last_update",
	})
	if err != nil {
		return nil, 0, err
	}
	var ids []uint
	for _, buffer := range houseList {
		id, parseErr := strconv.ParseUint(buffer.HouseId, 10, 0)
		if parseErr != nil {
			continue
		}
		ids = append(ids, uint(id))
	}
	if len(ids) == 0 {
		return []*house.Resource{}, int64(totalCount), nil
	}
	resources, err = service.GetListByIds(ids)
	return resources, int64(totalCount), err
}
