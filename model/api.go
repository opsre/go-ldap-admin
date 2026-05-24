package model

import "gorm.io/gorm"

type Api struct {
	gorm.Model
	Method          string `gorm:"type:varchar(20);comment:'请求方式'" json:"method"`
	Path            string `gorm:"type:varchar(100);comment:'访问路径'" json:"path"`
	Category        string `gorm:"type:varchar(50);comment:'所属类别'" json:"category"`
	CategoryDisplay string `gorm:"-" json:"categoryDisplay,omitempty"`
	Remark          string `gorm:"type:varchar(100);comment:'备注'" json:"remark"`
	RemarkDisplay   string `gorm:"-" json:"remarkDisplay,omitempty"`
	Creator         string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
	CreatorDisplay  string `gorm:"-" json:"creatorDisplay,omitempty"`
}
