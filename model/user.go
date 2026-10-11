package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username      string  `gorm:"type:varchar(50);not null;unique;comment:'Username'" json:"username"`                           // Username
	Password      string  `gorm:"size:255;not null;comment:'User password'" json:"password"`                                     // User password
	Nickname      string  `gorm:"type:varchar(50);comment:'Chinese name'" json:"nickname"`                                       // Nickname
	GivenName     string  `gorm:"type:varchar(50);comment:'Nickname'" json:"givenName"`                                          // Preferred name, if available; otherwise use the nickname
	Mail          string  `gorm:"type:varchar(100);comment:'Email'" json:"mail"`                                                 // Email
	JobNumber     string  `gorm:"type:varchar(20);comment:'Employee number'" json:"jobNumber"`                                   // Employee number
	Mobile        string  `gorm:"type:varchar(15);not null;unique;comment:'Mobile number'" json:"mobile"`                        // Mobile number
	Avatar        string  `gorm:"type:varchar(255);comment:'Avatar'" json:"avatar"`                                              // Avatar
	PostalAddress string  `gorm:"type:varchar(255);comment:'Address'" json:"postalAddress"`                                      // Address
	Departments   string  `gorm:"type:varchar(512);comment:'Department'" json:"departments"`                                     // Department
	Position      string  `gorm:"type:varchar(128);comment:'Position'" json:"position"`                                          //  Position
	Introduction  string  `gorm:"type:varchar(255);comment:'Personal introduction'" json:"introduction"`                         // Personal introduction
	Status        uint    `gorm:"type:tinyint(1);default:1;comment:'Status: 1 active, 2 departed'" json:"status"`                // Status
	Creator       string  `gorm:"type:varchar(20);;comment:'Creator'" json:"creator"`                                            // Creator
	Source        string  `gorm:"type:varchar(50);comment:'User source: DingTalk, WeCom, Feishu, LDAP, platform'" json:"source"` // Source
	DepartmentId  string  `gorm:"type:varchar(100);not null;comment:'Department ID'" json:"departmentId"`                        // Department ID
	Roles         []*Role `gorm:"many2many:user_roles" json:"roles"`                                                             // Role
	SourceUserId  string  `gorm:"type:varchar(100);not null;comment:'Third-party user ID'" json:"sourceUserId"`                  // Third-party user ID
	SourceUnionId string  `gorm:"type:varchar(100);not null;comment:'Unique third-party union ID'" json:"sourceUnionId"`         // Unique third-party union ID
	UserDN        string  `gorm:"type:varchar(255);not null;comment:'User DN'" json:"userDn"`                                    // User DN in LDAP
	SyncState     uint    `gorm:"type:tinyint(1);default:1;comment:'Sync status: 1 synced, 2 not synced'" json:"syncState"`      // LDAP sync status
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
