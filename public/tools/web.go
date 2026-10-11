package tools

type PageOption struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

var defaultOptions *PageOption

func init() {
	// By default, get 10 records from page 1
	defaultOptions = &PageOption{
		PageNum:  0,
		PageSize: 10,
	}
}

// NewPageOption creates pagination options
func NewPageOption(pageNum, pageSize int) *PageOption {
	if pageSize <= 0 || pageSize > 1000 || pageNum < 0 {
		return defaultOptions
	}

	pNum := (pageNum - 1) * pageSize
	return &PageOption{
		PageNum:  pNum,
		PageSize: pageSize,
	}
}
