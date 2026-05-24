package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name           string  `gorm:"type:varchar(20);not null;unique" json:"name"`
	NameDisplay    string  `gorm:"-" json:"nameDisplay,omitempty"`
	Keyword        string  `gorm:"type:varchar(20);not null;unique" json:"keyword"`
	Remark         string  `gorm:"type:varchar(100);comment:'备注'" json:"remark"`
	RemarkDisplay  string  `gorm:"-" json:"remarkDisplay,omitempty"`
	Status         uint    `gorm:"type:tinyint(1);default:1;comment:'1正常, 2禁用'" json:"status"`
	Sort           uint    `gorm:"type:int(3);default:999;comment:'角色排序(排序越大权限越低, 不能查看比自己序号小的角色, 不能编辑同序号用户权限, 排序为1表示超级管理员)'" json:"sort"`
	Creator        string  `gorm:"type:varchar(20);" json:"creator"`
	CreatorDisplay string  `gorm:"-" json:"creatorDisplay,omitempty"`
	Users          []*User `gorm:"many2many:user_roles" json:"users"`
	Menus          []*Menu `gorm:"many2many:role_menus;" json:"menus"` // 角色菜单多对多关系
}
