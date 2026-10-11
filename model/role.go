package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name           string  `gorm:"type:varchar(20);not null;unique" json:"name"`
	NameDisplay    string  `gorm:"-" json:"nameDisplay,omitempty"`
	Keyword        string  `gorm:"type:varchar(20);not null;unique" json:"keyword"`
	Remark         string  `gorm:"type:varchar(100);comment:'Remark'" json:"remark"`
	RemarkDisplay  string  `gorm:"-" json:"remarkDisplay,omitempty"`
	Status         uint    `gorm:"type:tinyint(1);default:1;comment:'1 active, 2 disabled'" json:"status"`
	Sort           uint    `gorm:"type:int(3);default:999;comment:'Role order (higher values mean lower permissions; cannot view lower-numbered roles or edit users at the same level; 1 is super administrator)'" json:"sort"`
	Creator        string  `gorm:"type:varchar(20);" json:"creator"`
	CreatorDisplay string  `gorm:"-" json:"creatorDisplay,omitempty"`
	Users          []*User `gorm:"many2many:user_roles" json:"users"`
	Menus          []*Menu `gorm:"many2many:role_menus;" json:"menus"` // 角色菜单多对多关系
}
