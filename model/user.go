package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username      string  `gorm:"type:varchar(50);not null;unique;comment:'Username'" json:"username"`                           // 用户名
	Password      string  `gorm:"size:255;not null;comment:'User password'" json:"password"`                                     // 用户密码
	Nickname      string  `gorm:"type:varchar(50);comment:'Chinese name'" json:"nickname"`                                       // 昵称
	GivenName     string  `gorm:"type:varchar(50);comment:'Nickname'" json:"givenName"`                                          // 花名，如果有的话，没有的话用昵称占位
	Mail          string  `gorm:"type:varchar(100);comment:'Email'" json:"mail"`                                                 // 邮箱
	JobNumber     string  `gorm:"type:varchar(20);comment:'Employee number'" json:"jobNumber"`                                   // 工号
	Mobile        string  `gorm:"type:varchar(15);not null;unique;comment:'Mobile number'" json:"mobile"`                        // 手机号
	Avatar        string  `gorm:"type:varchar(255);comment:'Avatar'" json:"avatar"`                                              // 头像
	PostalAddress string  `gorm:"type:varchar(255);comment:'Address'" json:"postalAddress"`                                      // 地址
	Departments   string  `gorm:"type:varchar(512);comment:'Department'" json:"departments"`                                     // 部门
	Position      string  `gorm:"type:varchar(128);comment:'Position'" json:"position"`                                          //  职位
	Introduction  string  `gorm:"type:varchar(255);comment:'Personal introduction'" json:"introduction"`                         // 个人简介
	Status        uint    `gorm:"type:tinyint(1);default:1;comment:'Status: 1 active, 2 departed'" json:"status"`                // 状态
	Creator       string  `gorm:"type:varchar(20);;comment:'Creator'" json:"creator"`                                            // 创建者
	Source        string  `gorm:"type:varchar(50);comment:'User source: DingTalk, WeCom, Feishu, LDAP, platform'" json:"source"` // 来源
	DepartmentId  string  `gorm:"type:varchar(100);not null;comment:'Department ID'" json:"departmentId"`                        // 部门id
	Roles         []*Role `gorm:"many2many:user_roles" json:"roles"`                                                             // 角色
	SourceUserId  string  `gorm:"type:varchar(100);not null;comment:'Third-party user ID'" json:"sourceUserId"`                  // 第三方用户id
	SourceUnionId string  `gorm:"type:varchar(100);not null;comment:'Unique third-party union ID'" json:"sourceUnionId"`         // 第三方唯一unionId
	UserDN        string  `gorm:"type:varchar(255);not null;comment:'User DN'" json:"userDn"`                                    // 用户在ldap的dn
	SyncState     uint    `gorm:"type:tinyint(1);default:1;comment:'Sync status: 1 synced, 2 not synced'" json:"syncState"`      // 数据到ldap的同步状态
}

func (u *User) SetUserName(userName string) {
	u.Username = userName
}

func (u *User) SetNickName(nickName string) {
	u.Nickname = nickName
}

func (u *User) SetGivenName(givenName string) {
	u.GivenName = givenName
}

func (u *User) SetMail(mail string) {
	u.Mail = mail
}

func (u *User) SetJobNumber(jobNum string) {
	u.JobNumber = jobNum
}

func (u *User) SetMobile(mobile string) {
	u.Mobile = mobile
}

func (u *User) SetAvatar(avatar string) {
	u.Avatar = avatar
}

func (u *User) SetPostalAddress(address string) {
	u.PostalAddress = address
}

func (u *User) SetPosition(position string) {
	u.Position = position
}

func (u *User) SetIntroduction(desc string) {
	u.Introduction = desc
}

func (u *User) SetSourceUserId(sourceUserId string) {
	u.SourceUserId = sourceUserId
}

func (u *User) SetSourceUnionId(sourceUnionId string) {
	u.SourceUnionId = sourceUnionId
}
