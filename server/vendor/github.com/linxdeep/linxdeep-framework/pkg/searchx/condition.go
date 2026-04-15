package searchx

type Index struct {
	Name string `xml:",chardata"`
}

type Condition struct {
	Ors    []Condition `xml:"or"`    // should
	Nots   []Condition `xml:"not"`   // must_not
	Terms  []Term      `xml:"term"`  // must
	Ranges []Range     `xml:"range"` // filter
	Aggs   []Agg       `xml:"agg"`
}
type Agg struct {
	AggType string `xml:"aggType,attr"` //  目前仅支持 Terms
	Field   string `xml:"field,attr"`
}
type Term struct {
	Field string `xml:"field,attr"`
	Value string `xml:"value,attr"`
}

type Range struct {
	Field      string `xml:"field,attr"`
	GreatEqual string `xml:"gte,attr"`
	GreatThan  string `xml:"gt,attr"`
	LessEqual  string `xml:"lte,attr"`
	LessThan   string `xml:"lt,attr"`
}

type QueryParams struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Order   string `json:"order"`
	By      string `json:"by"`
	Asc     bool
	Fields  []string
	Filters map[string][]string
}
