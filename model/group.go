package model

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	GroupName          string   `gorm:"type:varchar(128);comment:'分组名称'" json:"groupName"`
	GroupNameDisplay   string   `gorm:"-" json:"groupNameDisplay,omitempty"`
	Remark             string   `gorm:"type:varchar(128);comment:'分组中文说明'" json:"remark"`
	RemarkDisplay      string   `gorm:"-" json:"remarkDisplay,omitempty"`
	Creator            string   `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
	CreatorDisplay     string   `gorm:"-" json:"creatorDisplay,omitempty"`
	GroupType          string   `gorm:"type:varchar(20);comment:'分组类型：cn、ou'" json:"groupType"`
	Users              []*User  `gorm:"many2many:group_users" json:"users"`
	ParentId           uint     `gorm:"default:0;comment:'父组编号(编号为0时表示根组)'" json:"parentId"`
	SourceDeptId       string   `gorm:"type:varchar(100);comment:'部门编号'" json:"sourceDeptId"`
	Source             string   `gorm:"type:varchar(20);comment:'来源：dingTalk、weCom、ldap、platform'" json:"source"`
	SourceDeptParentId string   `gorm:"type:varchar(100);comment:'父部门编号'" json:"sourceDeptParentId"`
	SourceUserNum      int      `gorm:"default:0;comment:'部门下的用户数量，从第三方获取的数据'" json:"source_user_num"`
	Children           []*Group `gorm:"-" json:"children"`
	GroupDN            string   `gorm:"type:varchar(255);not null;comment:'分组dn'" json:"groupDn"`             // 分组在ldap的dn
	SyncState          uint     `gorm:"type:tinyint(1);default:1;comment:'同步状态:1已同步, 2未同步'" json:"syncState"` // 数据到ldap的同步状态
}

func (g *Group) SetGroupName(groupName string) {
	g.GroupName = groupName
}

func (g *Group) SetRemark(remark string) {
	g.Remark = remark
}

func (g *Group) SetSourceDeptId(sourceDeptId string) {
	g.SourceDeptId = sourceDeptId
}

func (g *Group) SetSourceDeptParentId(sourceDeptParentId string) {
	g.SourceDeptParentId = sourceDeptParentId
}
