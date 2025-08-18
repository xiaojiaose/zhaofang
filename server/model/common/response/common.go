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
	Name      string `json:"name"`
	Id        int    `json:"id"`
	Sort      int    `json:"sort"`
	Districts []Districts
}

type Districts struct {
	Name string `json:"name"`
	Id   uint   `json:"id"`
}
