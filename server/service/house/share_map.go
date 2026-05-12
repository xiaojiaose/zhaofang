package house

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	response2 "github.com/flipped-aurora/gin-vue-admin/server/model/house/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/linxdeep/linxdeep-framework/pkg/searchx"
)

func (service *ResourceService) SharedMapAgg(userID uint, filter *request.ResourceSearch) (list []response2.SharedMapXiaoqu, err error) {
	condition := searchx.Condition{
		Terms: []searchx.Term{
			{Field: "owner", Value: strconv.Itoa(int(userID))},
			{Field: "status", Value: "待出租"},
		},
		Aggs: []searchx.Agg{{Field: "xiaoqu_id"}},
		Ors:  []searchx.Condition{{Terms: []searchx.Term{}}},
		Nots: []searchx.Condition{{Terms: []searchx.Term{}}},
	}
	if filter != nil {
		service.applyFilterToCondition(&condition, filter)
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

func (service *ResourceService) SharedMapList(userID, xiaoquID uint, info request.PageInfo, filter *request.ResourceSearch) (resources []*house.Resource, total int64, err error) {
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
	if filter != nil {
		service.applyFilterToCondition(&condition, filter)
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

func (service *ResourceService) applyFilterToCondition(condition *searchx.Condition, filter *request.ResourceSearch) {
	if filter == nil {
		return
	}
	if len(filter.Feature) > 0 {
		for _, f := range strings.Split(filter.Feature, ",") {
			condition.Terms = append(condition.Terms, searchx.Term{Field: "feature", Value: "*" + f + "*"})
		}
	}
	if len(filter.HouseType) > 0 {
		for _, v := range []string{"1居", "2居", "3居", "4居+", "开间", "主卧", "次卧", "暗间"} {
			if !strings.Contains(filter.HouseType, v) {
				condition.Nots[0].Terms = append(condition.Nots[0].Terms, searchx.Term{Field: "house_type", Value: v})
			}
		}
	}
	if len(filter.RentType) > 0 {
		condition.Terms = append(condition.Terms, searchx.Term{Field: "rent_type", Value: filter.RentType + "*"})
	}
	if filter.Price > 0 {
		priceOption := getPriceRange(filter.Price)
		if len(priceOption) == 2 {
			condition.Ranges = append(condition.Ranges, searchx.Range{Field: "price", GreatEqual: fmt.Sprintf("%d", priceOption[0]), LessEqual: fmt.Sprintf("%d", priceOption[1])})
		}
	}
	if len(filter.HouseSource) > 0 {
		applyHouseSourceCondition(condition, filter.HouseSource, true)
	}
	if len(filter.XiaoquId) > 0 {
		for _, i := range filter.XiaoquId {
			condition.Ors[0].Terms = append(condition.Ors[0].Terms, searchx.Term{Field: "xiaoqu_id", Value: strconv.Itoa(i)})
		}
	}
}

var priceRangeMap = map[int][]int{
	1: {0, 500},
	2: {500, 1000},
	3: {1000, 1500},
	4: {1500, 2000},
	5: {2000, 2500},
	6: {2500, 3000},
	7: {3000, 100000},
}

func getPriceRange(price int) []int {
	if v, ok := priceRangeMap[price]; ok {
		return v
	}
	return nil
}

func applyHouseSourceCondition(condition *searchx.Condition, houseSource string, allowTeam bool) {
	for _, item := range strings.Split(houseSource, ",") {
		switch strings.TrimSpace(item) {
		case "commission", "有返佣":
			condition.Ranges = append(condition.Ranges, searchx.Range{Field: "commission_price", GreatEqual: "1"})
		case "landlord", "房东房源":
			condition.Terms = append(condition.Terms, searchx.Term{Field: "rent_type", Value: "房东房源"})
		case "team", "团队房源":
			if allowTeam {
				condition.Terms = append(condition.Terms, searchx.Term{Field: "is_team_house", Value: "1"})
			}
		}
	}
}
