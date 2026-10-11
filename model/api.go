package model

import "gorm.io/gorm"

type Api struct {
	gorm.Model
	Method          string `gorm:"type:varchar(20);comment:'Request method'" json:"method"`
	Path            string `gorm:"type:varchar(100);comment:'Access path'" json:"path"`
	Category        string `gorm:"type:varchar(50);comment:'Category'" json:"category"`
	CategoryDisplay string `gorm:"-" json:"categoryDisplay,omitempty"`
	Remark          string `gorm:"type:varchar(100);comment:'Remark'" json:"remark"`
	RemarkDisplay   string `gorm:"-" json:"remarkDisplay,omitempty"`
	Creator         string `gorm:"type:varchar(20);comment:'Creator'" json:"creator"`
	CreatorDisplay  string `gorm:"-" json:"creatorDisplay,omitempty"`
}
