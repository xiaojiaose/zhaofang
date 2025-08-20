package response

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"` // 总数
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

type ReturnList struct {
	List interface{} `json:"list"`
}

type ReturnValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Area struct {
	Name      string      `json:"name"` // 行政区name
	Id        int         `json:"id"`   // 行政区id
	Sort      int         `json:"sort"` // 排序
	Districts []Districts // 行政区下面的商圈
}

type Districts struct {
	Name string `json:"name"` // 商圈name
	Id   uint   `json:"id"`   // 商圈id
}

type XiaoquRsp struct {
	Id   int    `json:"id"`   // 小区id
	Name string `json:"name"` // 小区名
	Num  int    `json:"num"`  // 房源数量
}
