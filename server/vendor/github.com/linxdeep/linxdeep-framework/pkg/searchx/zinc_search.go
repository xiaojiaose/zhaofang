//go:generate mockgen -source=./zinc_search.go -destination=./mock_zinc_search.go -package=searchx
package searchx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	client "github.com/zinclabs/sdk-go-zincsearch"
)

type Searcher interface {
	Add(ctx context.Context, index string, id string, data map[string]interface{}) error
	AddMulti(ctx context.Context, index string, data []map[string]interface{}) error
	Get(
		ctx context.Context,
		index string,
		ids []string,
		fields ...string,
	) (result []map[string]interface{}, err error)
	Search(
		ctx context.Context,
		index string,
		condition Condition,
		params QueryParams,
	) (result []map[string]interface{}, total int, err error)
	SearchAgg(
		ctx context.Context,
		index string,
		condition Condition,
		params QueryParams,
	) (result []map[string]interface{}, total int, agg map[string]map[string]int, err error)
	SearchV1Agg(
		ctx context.Context,
		index string,
		keyword string, aggField string) (result map[string]int, err error)
	SearchV1(
		ctx context.Context,
		index string,
		keyword string) (result []map[string]interface{}, total int, err error)
	Delete(ctx context.Context, index string, id string) error
}

var (
	_ Searcher = (*search)(nil)
)

type search struct {
	client *client.APIClient
	auth   *client.BasicAuth
}

func New(url, username, password string) Searcher {
	config := client.NewConfiguration()
	config.Servers = client.ServerConfigurations{
		client.ServerConfiguration{
			URL: url,
		},
	}
	_globalSearcher = &search{
		client: client.NewAPIClient(config),
		auth: &client.BasicAuth{
			UserName: username,
			Password: password,
		},
	}
	return _globalSearcher
}

// Add 添加或更新单条文档，只覆盖 data 中给定的字段
func (e *search) Add(ctx context.Context, index string, id string, data map[string]interface{}) error {
	if a, b, err := e.client.Document.IndexWithID(e.Context(ctx), index, id).Document(data).Execute(); err != nil {
		slog.LogAttrs(
			e.Context(ctx),
			slog.LevelWarn,
			err.Error(),
			slog.String("zincsearch", "add"),
			slog.String("operation", fmt.Sprintf("Failed to add data %v: a=%v, b=%v", data, a, b)),
		)
		return err
	}
	return nil
}

// AddMulti 批量添加或更新文档，只覆盖 data 中给定的字段
func (e *search) AddMulti(ctx context.Context, index string, data []map[string]interface{}) error {
	query := client.MetaJSONIngest{
		Index:   &index,
		Records: data,
	}

	if a, b, err := e.client.Document.Bulkv2(e.Context(ctx)).Query(query).Execute(); err != nil {
		slog.LogAttrs(
			e.Context(ctx),
			slog.LevelWarn,
			err.Error(),
			slog.String("zincsearch", "addMulti"),
			slog.String("operation", fmt.Sprintf("Failed to add data %v: a=%v, b=%v", data, a, b)),
		)
		return err
	}

	return nil
}

func (e *search) Get(
	ctx context.Context,
	index string,
	ids []string,
	fields ...string,
) (result []map[string]interface{}, err error) {
	idsQuery := *client.NewMetaIdsQuery()
	idsQuery.SetValues(ids)

	subQuery := *client.NewMetaQuery()
	subQuery.SetIds(idsQuery)

	query := *client.NewMetaZincQuery()
	filter := *client.NewMetaQuery()
	boolQuery := *client.NewMetaBoolQuery()
	boolQuery.SetMust([]client.MetaQuery{subQuery})
	filter.SetBool(boolQuery)
	query.SetQuery(filter)
	query.SetSource(fields)

	resp, _, err := e.client.Search.Search(e.Context(ctx), index).Query(query).Execute()
	if err != nil {
		slog.LogAttrs(e.Context(ctx), slog.LevelWarn, err.Error(), slog.String("zincsearch", "get"))
		return
	}

	for _, hit := range resp.Hits.Hits {
		result = append(result, hit.GetSource())
	}

	return
}

// Search 查询文档
func (e *search) Search(
	ctx context.Context,
	index string,
	condition Condition,
	params QueryParams,
) (result []map[string]interface{}, total int, err error) {
	query := *client.NewMetaZincQuery()
	filter := *client.NewMetaQuery()
	boolQuery := *client.NewMetaBoolQuery()

	buildFlatQueries := func(condition Condition) (result []client.MetaQuery) {
		for _, term := range condition.Terms {
			if term.Field == "*" {
				subQuery := *client.NewMetaQuery()
				multiMatch := *client.NewMetaMultiMatchQuery()
				multiMatch.SetQuery(term.Value)
				multiMatch.SetFields([]string{"family", "name", "status"}) // 使用 * 表示所有字段
				multiMatch.SetType("most_fields")                          // 也可以使用 "most_fields" 或 "cross_fields"
				multiMatch.SetMinimumShouldMatch(1)
				multiMatch.SetOperator("OR")
				multiMatch.SetAnalyzer("ik_smart")

				subQuery.SetMultiMatch(multiMatch)
				result = append(result, subQuery)
			} else if term.Value != "" {
				subQuery := *client.NewMetaQuery()
				if strings.Contains(term.Value, "*") {
					wildcardQuery := *client.NewMetaWildcardQuery()
					wildcardQuery.SetValue(term.Value)
					// 指定字段的模糊查询
					subQuery.SetWildcard(map[string]client.MetaWildcardQuery{term.Field: wildcardQuery})
				} else {
					matchQuery := *client.NewMetaMatchQuery()
					matchQuery.SetQuery(term.Value)

					subQuery.SetMatch(map[string]client.MetaMatchQuery{term.Field: matchQuery})
				}

				result = append(result, subQuery)
			}
		}

		return
	}

	if len(condition.Terms) > 0 {
		boolQuery.SetMust(buildFlatQueries(condition))
	}

	if len(condition.Nots) > 0 {
		queries := make([]client.MetaQuery, 0, len(condition.Nots))
		for _, not := range condition.Nots {
			queries = append(queries, buildFlatQueries(not)...)
		}

		boolQuery.SetMustNot(queries)
	}

	if len(condition.Ors) > 0 {
		queries := make([]client.MetaQuery, 0, len(condition.Ors))
		for _, or := range condition.Ors {
			queries = append(queries, buildFlatQueries(or)...)
		}

		boolQuery.SetShould(queries)
		boolQuery.SetMinimumShouldMatch(1)
	}

	if len(condition.Ranges) > 0 {
		queries := make([]client.MetaQuery, 0, len(condition.Ranges))
		for _, r := range condition.Ranges {
			subQueryFilterQuery := *client.NewMetaRangeQuery()
			subQueryFilterQuery.SetGte(r.GreatEqual) // >=
			subQueryFilterQuery.SetLte(r.LessEqual)  // <=
			a := *client.NewMetaQuery()
			a.SetRange(map[string]client.MetaRangeQuery{
				r.Field: subQueryFilterQuery,
			})
			queries = append(queries, a)
		}
		boolQuery.SetFilter(queries)
	}

	filter.SetBool(boolQuery)
	query.SetQuery(filter)

	var sort string
	if params.Asc {
		sort = params.By
	} else {
		sort = fmt.Sprintf("-%s", params.By)
	}

	query.SetSort([]string{sort}) // "-@timestamp"

	query.SetFrom(int32(params.Size * params.Page))
	query.SetSize(int32(params.Size))

	query.SetSource(append(params.Fields, "id"))

	resp, httpResp, err := e.client.Search.Search(e.Context(ctx), index).Query(query).Execute()
	if err != nil {
		var apiError *client.GenericOpenAPIError
		if errors.As(err, &apiError) {
			body := string(apiError.Body())
			if strings.Contains(body, "index") && strings.Contains(body, "does not exists") {
				err = nil
				return
			}
		}

		slog.LogAttrs(
			e.Context(ctx),
			slog.LevelWarn,
			err.Error(),
			slog.String("zincsearch", "search"),
			slog.String("index", index),
			slog.String("operation", fmt.Sprintf("failed fetch data : resp=%v, httpResp=%v", resp, httpResp)),
		)

		return
	}

	for _, hit := range resp.Hits.Hits {
		result = append(result, hit.GetSource())
	}

	total = int(*resp.Hits.Total.Value)
	return
}

func (e *search) SearchAgg(
	ctx context.Context,
	index string,
	condition Condition,
	params QueryParams,
) (result []map[string]interface{}, total int, agg map[string]map[string]int, err error) {
	query := *client.NewMetaZincQuery()
	filter := *client.NewMetaQuery()
	boolQuery := *client.NewMetaBoolQuery()

	buildFlatQueries := func(condition Condition) (result []client.MetaQuery) {
		for _, term := range condition.Terms {
			if term.Field == "*" {
				subQuery := *client.NewMetaQuery()
				multiMatch := *client.NewMetaMultiMatchQuery()
				multiMatch.SetQuery(term.Value)
				multiMatch.SetFields([]string{"family", "name", "status"}) // 使用 * 表示所有字段
				multiMatch.SetType("most_fields")                          // 也可以使用 "most_fields" 或 "cross_fields"
				multiMatch.SetMinimumShouldMatch(1)
				multiMatch.SetOperator("OR")
				multiMatch.SetAnalyzer("ik_smart")

				subQuery.SetMultiMatch(multiMatch)
				result = append(result, subQuery)
			} else if term.Value != "" {
				subQuery := *client.NewMetaQuery()
				if strings.Contains(term.Value, "*") {
					wildcardQuery := *client.NewMetaWildcardQuery()
					wildcardQuery.SetValue(term.Value)
					// 指定字段的模糊查询
					subQuery.SetWildcard(map[string]client.MetaWildcardQuery{term.Field: wildcardQuery})
				} else {
					matchQuery := *client.NewMetaMatchQuery()
					matchQuery.SetQuery(term.Value)

					subQuery.SetMatch(map[string]client.MetaMatchQuery{term.Field: matchQuery})
				}

				result = append(result, subQuery)
			}
		}

		return
	}

	if len(condition.Terms) > 0 {
		boolQuery.SetMust(buildFlatQueries(condition))
	}

	if len(condition.Nots) > 0 {
		queries := make([]client.MetaQuery, 0, len(condition.Nots))
		for _, not := range condition.Nots {
			queries = append(queries, buildFlatQueries(not)...)
		}

		boolQuery.SetMustNot(queries)
	}

	if len(condition.Ors) > 0 {
		queries := make([]client.MetaQuery, 0, len(condition.Ors))
		for _, or := range condition.Ors {
			queries = append(queries, buildFlatQueries(or)...)
		}

		boolQuery.SetShould(queries)
		boolQuery.SetMinimumShouldMatch(1)
	}

	if len(condition.Ranges) > 0 {
		queries := make([]client.MetaQuery, 0, len(condition.Ranges))
		for _, r := range condition.Ranges {
			subQueryFilterQuery := *client.NewMetaRangeQuery()
			subQueryFilterQuery.SetGte(r.GreatEqual) // >=
			subQueryFilterQuery.SetLte(r.LessEqual)  // <=
			a := *client.NewMetaQuery()
			a.SetRange(map[string]client.MetaRangeQuery{
				r.Field: subQueryFilterQuery,
			})
			queries = append(queries, a)
		}
		boolQuery.SetFilter(queries)
	}

	filter.SetBool(boolQuery)
	query.SetQuery(filter)

	var sort string
	if params.Asc {
		sort = params.By
	} else {
		sort = fmt.Sprintf("-%s", params.By)
	}

	query.SetSort([]string{sort}) // "-@timestamp"

	query.SetFrom(int32(params.Size * params.Page))
	query.SetSize(int32(params.Size))

	query.SetSource(append(params.Fields, "id"))

	if len(condition.Aggs) > 0 {
		aggs := make(map[string]client.MetaAggregations)
		for _, ag := range condition.Aggs {
			aggs[ag.Field] = client.MetaAggregations{
				Terms: &client.MetaAggregationsTerms{ // Terms	按字段分组统计 Avg	计算平均值 Sum	计算总和 Min	计算最小值 DateHistogram	按日期区间分组
					Field: StringPtr(ag.Field), // 按 brand 字段分组
					Size:  IntPtr(30),          // 返回前 10 个分组
				},
			}
		}

		query.SetAggs(aggs)
	}

	resp, httpResp, err := e.client.Search.Search(e.Context(ctx), index).Query(query).Execute()
	if err != nil && !strings.Contains(err.Error(), "cannot unmarshal array into Go struct field MetaAggregationResponse.aggregations.buckets") {
		slog.LogAttrs(
			e.Context(ctx),
			slog.LevelWarn,
			err.Error(),
			slog.String("zincsearch", "search"),
			slog.String("operation", fmt.Sprintf("failed fetch data : resp=%v, httpResp=%v", resp, httpResp)),
		)
		var apiError *client.GenericOpenAPIError
		if errors.As(err, &apiError) {
			body := string(apiError.Body())
			if strings.Contains(body, "index") && strings.Contains(body, "does not exists") {
				err = nil
			}
		}
		return
	}

	for _, hit := range resp.Hits.Hits {
		result = append(result, hit.GetSource())
	}

	total = int(*resp.Hits.Total.Value)

	if len(*resp.Aggregations) > 0 {
		var body []byte
		body, err = io.ReadAll(httpResp.Body)
		if err != nil {
			slog.LogAttrs(
				e.Context(ctx),
				slog.LevelWarn,
				err.Error(),
				slog.String("zincsearch", "search"),
				slog.String("operation", fmt.Sprintf("failed to read request body : resp=%v, httpResp=%v", resp, httpResp)),
			)
			return
		}
		var r AggregationResponse
		err = json.Unmarshal(body, &r)
		if err != nil {
			slog.LogAttrs(
				e.Context(ctx),
				slog.LevelWarn,
				err.Error(),
				slog.String("zincsearch", "search"),
				slog.String("operation", fmt.Sprintf("unmarshal body : resp=%v, httpResp=%v", resp, httpResp)),
			)
			return
		}
		agg = make(map[string]map[string]int, len(r.Aggregations))
		for cate, v := range r.Aggregations {
			agg[cate] = make(map[string]int, len(v.Buckets))
			for _, bucket := range v.Buckets {
				agg[cate][fmt.Sprintf("%v", bucket.Key)] = bucket.DocCount
			}
		}
		fmt.Println(r)
	}
	return
}

// SearchV1 适用于 简单搜索  全文搜索 不指定字段 不进行 各种 and  or  not 等操作
func (e *search) SearchV1(
	ctx context.Context,
	index string,
	keyword string) (result []map[string]interface{}, total int, err error) {

	//  换小写， zincsearch 对大小搜索 不支持，保存的时候 索引按小写保存的
	keyword = strings.ToLower(keyword)

	con := client.NewV1ZincQuery()
	// fuzzy 模糊 不进行分词，match 用分词器进行分析索引 “-” 是分隔符， 例如 ： 设备id= f2d305de-1b81-483d-9fc4-84954e5d5929，systemSN = "123456"
	// match 会对 输入的 keyword 为 "f2d305-123456" 会分成 "f2d305" 和 "123456"，分别和 设备id 和 systemSN 进行匹配，目标找到 "123456" （注意不是部分，是整个索引），
	// fuzzy 会对 输入的 keyword 为 "f2d305-123456" 不会进行分词，
	con.SetSearchType("term") // term 精确匹配 不分词
	con.SetQuery(
		client.V1QueryParams{
			Term: &keyword,
		})
	con.SetFrom(0)
	con.SetMaxResults(10)
	resp, httpResp, err := e.client.Search.SearchV1(e.Context(ctx), index).Query(*con).Execute()
	if err != nil {
		slog.LogAttrs(
			e.Context(ctx),
			slog.LevelWarn,
			err.Error(),
			slog.String("zincsearch", "searchV1"),
			slog.String("operation", fmt.Sprintf("failed fetch data : resp=%v, httpResp=%v", resp, httpResp)),
		)
		return
	}

	for _, hit := range resp.Hits.Hits {
		result = append(result, hit.GetSource())
	}

	total = int(*resp.Hits.Total.Value)
	return
}

func (e *search) SearchV1Agg(
	ctx context.Context,
	index string,
	keyword string, aggField string) (result map[string]int, err error) {

	if len(aggField) == 0 {
		err = errors.New("aggField is empty")
		return
	}

	//  换小写
	keyword = strings.ToLower(keyword)

	con := client.NewV1ZincQuery()
	con.SetSearchType("match") // fuzzy 模糊，match 用分词器进行分析索引 “-” 是分隔符
	con.SetQuery(
		client.V1QueryParams{
			Term: &keyword,
		})

	aggType := "terms"
	agg := make(map[string]client.V1AggregationParams)
	agg[aggField] = client.V1AggregationParams{
		AggType: &aggType,
		Field:   &aggField,
		//Size:    10,
	}
	con.SetAggs(agg)

	resp, httpResp, err := e.client.Search.SearchV1(e.Context(ctx), index).Query(*con).Execute()
	if err != nil && !strings.Contains(err.Error(), "cannot unmarshal array into Go struct field MetaAggregationResponse.aggregations.buckets") {
		slog.LogAttrs(
			e.Context(ctx),
			slog.LevelWarn,
			err.Error(),
			slog.String("zincsearch", "search"),
			slog.String("operation", fmt.Sprintf("failed fetch data : resp=%v, httpResp=%v", resp, httpResp)),
		)
		return
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		slog.LogAttrs(
			e.Context(ctx),
			slog.LevelWarn,
			err.Error(),
			slog.String("zincsearch", "searchV1"),
			slog.String("operation", fmt.Sprintf("Failed to read request body : resp=%v, httpResp=%v", resp, httpResp)),
		)
		return
	}
	var r AggregationResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		slog.LogAttrs(
			e.Context(ctx),
			slog.LevelWarn,
			err.Error(),
			slog.String("zincsearch", "searchV1"),
			slog.String("operation", fmt.Sprintf("unmarshal body : resp=%v, httpResp=%v", resp, httpResp)),
		)
		return
	}
	result = make(map[string]int)
	for _, v := range r.Aggregations {
		for _, bucket := range v.Buckets {
			result[fmt.Sprintf("%v", bucket.Key)] = bucket.DocCount
		}
	}

	fmt.Println(r)
	return
}

// Delete 删除指定索引中的文档
func (e *search) Delete(ctx context.Context, index string, id string) error {
	if a, b, err := e.client.Document.Delete(e.Context(ctx), index, id).Execute(); err != nil {
		slog.LogAttrs(
			e.Context(ctx),
			slog.LevelWarn,
			err.Error(),
			slog.String("zincsearch", "delete"),
			slog.String("operation", fmt.Sprintf("failed fetch data : a=%v, b=%v", a, b)),
		)
		return err
	}

	return nil
}

func (e *search) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, client.ContextBasicAuth, *e.auth)
}

type Bucket struct {
	DocCount    int         `json:"doc_count"`
	Key         interface{} `json:"key"`
	KeyAsString string      `json:"key_as_string,omitempty"`
}

type AggregationResponse struct {
	Aggregations map[string]struct {
		Buckets []Bucket `json:"buckets"`
	} `json:"aggregations"`
}

func StringPtr(v string) *string {
	return &v
}
func IntPtr(v int32) *int32 {
	return &v
}
