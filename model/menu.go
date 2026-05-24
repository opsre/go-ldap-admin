package model

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	Name           string  `gorm:"type:varchar(50);comment:'菜单名称(英文名, 可用于国际化)'" json:"name"`
	NameDisplay    string  `gorm:"-" json:"nameDisplay,omitempty"`
	Title          string  `gorm:"type:varchar(50);comment:'菜单标题(无法国际化时使用)'" json:"title"`
	TitleDisplay   string  `gorm:"-" json:"titleDisplay,omitempty"`
	Icon           string  `gorm:"type:varchar(50);comment:'菜单图标'" json:"icon"`
	Path           string  `gorm:"type:varchar(100);comment:'菜单访问路径'" json:"path"`
	Redirect       string  `gorm:"type:varchar(100);comment:'重定向路径'" json:"redirect"`
	Component      string  `gorm:"type:varchar(100);comment:'前端组件路径'" json:"component"`
	Sort           uint    `gorm:"type:int(3);default:999;comment:'菜单顺序(1-999)'" json:"sort"`
	Status         uint    `gorm:"type:tinyint(1);default:1;comment:'菜单状态(正常/禁用, 默认正常)'" json:"status"`
	Hidden         uint    `gorm:"type:tinyint(1);default:2;comment:'菜单在侧边栏隐藏(1隐藏，2显示)'" json:"hidden"`
	NoCache        uint    `gorm:"type:tinyint(1);default:2;comment:'菜单是否被 <keep-alive> 缓存(1不缓存，2缓存)'" json:"noCache"`
	AlwaysShow     uint    `gorm:"type:tinyint(1);default:2;comment:'忽略之前定义的规则，一直显示根路由(1忽略，2不忽略)'" json:"alwaysShow"`
	Breadcrumb     uint    `gorm:"type:tinyint(1);default:1;comment:'面包屑可见性(可见/隐藏, 默认可见)'" json:"breadcrumb"`
	ActiveMenu     string  `gorm:"type:varchar(100);comment:'在其它路由时，想在侧边栏高亮的路由'" json:"activeMenu"`
	ParentId       uint    `gorm:"default:0;comment:'父菜单编号(编号为0时表示根菜单)'" json:"parentId"`
	Creator        string  `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
	CreatorDisplay string  `gorm:"-" json:"creatorDisplay,omitempty"`
	Children       []*Menu `gorm:"-" json:"children"`                  // 子菜单集合
	Roles          []*Role `gorm:"many2many:role_menus;" json:"roles"` // 角色菜单多对多关系
}
