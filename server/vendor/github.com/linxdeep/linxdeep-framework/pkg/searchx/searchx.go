package searchx

import "context"

var (
	_globalSearcher Searcher
)

func SetGlobal(searcher Searcher) {
	_globalSearcher = searcher
}

func Add(ctx context.Context, index string, id string, data map[string]interface{}) error {
	return _globalSearcher.Add(ctx, index, id, data)
}
func AddMulti(ctx context.Context, index string, data []map[string]interface{}) error {
	return _globalSearcher.AddMulti(ctx, index, data)
}
func Get(
	ctx context.Context,
	index string,
	ids []string,
	fields ...string,
) (result []map[string]interface{}, err error) {
	return _globalSearcher.Get(ctx, index, ids, fields...)
}
func Search(
	ctx context.Context,
	index string,
	condition Condition,
	params QueryParams,
) (result []map[string]interface{}, total int, err error) {
	return _globalSearcher.Search(ctx, index, condition, params)
}

func Delete(ctx context.Context, index string, id string) error {
	return _globalSearcher.Delete(ctx, index, id)
}
