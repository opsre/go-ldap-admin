package model

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	Name           string  `gorm:"type:varchar(50);comment:'Menu name (English name for localization)'" json:"name"`
	NameDisplay    string  `gorm:"-" json:"nameDisplay,omitempty"`
	Title          string  `gorm:"type:varchar(50);comment:'Menu title (used when localization is unavailable)'" json:"title"`
	TitleDisplay   string  `gorm:"-" json:"titleDisplay,omitempty"`
	Icon           string  `gorm:"type:varchar(50);comment:'Menu icon'" json:"icon"`
	Path           string  `gorm:"type:varchar(100);comment:'Menu path'" json:"path"`
	Redirect       string  `gorm:"type:varchar(100);comment:'Redirect path'" json:"redirect"`
	Component      string  `gorm:"type:varchar(100);comment:'Frontend component path'" json:"component"`
	Sort           uint    `gorm:"type:int(3);default:999;comment:'Menu order (1-999)'" json:"sort"`
	Status         uint    `gorm:"type:tinyint(1);default:1;comment:'Menu status (active/disabled, active by default)'" json:"status"`
	Hidden         uint    `gorm:"type:tinyint(1);default:2;comment:'Hide menu in sidebar (1 hidden, 2 visible)'" json:"hidden"`
	NoCache        uint    `gorm:"type:tinyint(1);default:2;comment:'Whether menu is cached by <keep-alive> (1 no, 2 yes)'" json:"noCache"`
	AlwaysShow     uint    `gorm:"type:tinyint(1);default:2;comment:'Ignore previous rules and always show root route (1 yes, 2 no)'" json:"alwaysShow"`
	Breadcrumb     uint    `gorm:"type:tinyint(1);default:1;comment:'Breadcrumb visibility (visible/hidden, visible by default)'" json:"breadcrumb"`
	ActiveMenu     string  `gorm:"type:varchar(100);comment:'Route to highlight in sidebar on other routes'" json:"activeMenu"`
	ParentId       uint    `gorm:"default:0;comment:'Parent menu ID (0 indicates the root menu)'" json:"parentId"`
	Creator        string  `gorm:"type:varchar(20);comment:'Creator'" json:"creator"`
	CreatorDisplay string  `gorm:"-" json:"creatorDisplay,omitempty"`
	Children       []*Menu `gorm:"-" json:"children"`                  // 子菜单集合
	Roles          []*Role `gorm:"many2many:role_menus;" json:"roles"` // 角色菜单多对多关系
}
