package search

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
	"github.com/linxdeep/linxdeep-framework/pkg/searchx"
)

type ResourceSearch struct {
	engine searchx.Searcher
	index  string
}

const IndexDevices = "house"

func NewResource(searcher searchx.Searcher) *ResourceSearch {
	return &ResourceSearch{
		engine: searcher,
		index:  IndexDevices + "-resource",
	}
}

func (d ResourceSearch) Add(ctx context.Context, device search.ResourceBuffer) error {
	if err := d.engine.Add(ctx, d.index, device.ID, device.ToData()); err != nil {
		return err
	}
	return nil
}

func (d ResourceSearch) Search(
	ctx context.Context,
	condition searchx.Condition,
	params searchx.QueryParams,
) (result []*search.ResourceBuffer, total int, err error) {
	var data []map[string]interface{}
	if data, total, err = d.engine.Search(ctx, d.index, condition, params); err != nil {
		return
	}
	result = make([]*search.ResourceBuffer, 0, len(data))
	for _, item := range data {
		entity := search.FromDeviceES(item)
		result = append(result, entity)
	}
	return
}

func (d ResourceSearch) SearchAgg(
	ctx context.Context,
	condition searchx.Condition,
	params searchx.QueryParams,
) (result []*search.ResourceBuffer, total int, agg map[string]map[string]int, err error) {
	var data []map[string]interface{}
	//var agg map[string]map[string]int
	if data, total, agg, err = d.engine.SearchAgg(ctx, d.index, condition, params); err != nil {
		return
	}
	result = make([]*search.ResourceBuffer, 0, len(data))
	for _, item := range data {
		entity := search.FromDeviceES(item)
		result = append(result, entity)
	}
	return
}
