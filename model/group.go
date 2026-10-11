package model

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	GroupName          string   `gorm:"type:varchar(128);comment:'Group name'" json:"groupName"`
	GroupNameDisplay   string   `gorm:"-" json:"groupNameDisplay,omitempty"`
	Remark             string   `gorm:"type:varchar(128);comment:'Group description in Chinese'" json:"remark"`
	RemarkDisplay      string   `gorm:"-" json:"remarkDisplay,omitempty"`
	Creator            string   `gorm:"type:varchar(20);comment:'Creator'" json:"creator"`
	CreatorDisplay     string   `gorm:"-" json:"creatorDisplay,omitempty"`
	GroupType          string   `gorm:"type:varchar(20);comment:'Group type: cn, ou'" json:"groupType"`
	Users              []*User  `gorm:"many2many:group_users" json:"users"`
	ParentId           uint     `gorm:"default:0;comment:'Parent group ID (0 indicates the root group)'" json:"parentId"`
	SourceDeptId       string   `gorm:"type:varchar(100);comment:'Department ID'" json:"sourceDeptId"`
	Source             string   `gorm:"type:varchar(20);comment:'Source: DingTalk, WeCom, LDAP, platform'" json:"source"`
	SourceDeptParentId string   `gorm:"type:varchar(100);comment:'Parent department ID'" json:"sourceDeptParentId"`
	SourceUserNum      int      `gorm:"default:0;comment:'Number of users in department from third-party data'" json:"source_user_num"`
	Children           []*Group `gorm:"-" json:"children"`
	GroupDN            string   `gorm:"type:varchar(255);not null;comment:'Group DN'" json:"groupDn"`                             // Group DN in LDAP
	SyncState          uint     `gorm:"type:tinyint(1);default:1;comment:'Sync status: 1 synced, 2 not synced'" json:"syncState"` // LDAP sync status
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
