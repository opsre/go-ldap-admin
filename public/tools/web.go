package tools

type PageOption struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

var defaultOptions *PageOption

func init() {
	// 默认取 第 1 页的 10 条记录
	defaultOptions = &PageOption{
		PageNum:  0,
		PageSize: 10,
	}
}

// NewPageOption 创建一个分页参数
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
