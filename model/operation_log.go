package model

import (
	"gorm.io/gorm"
)

type OperationLog struct {
	gorm.Model
	Username   string `gorm:"type:varchar(20);comment:'Login username'" json:"username"`
	Ip         string `gorm:"type:varchar(20);comment:'IP address'" json:"ip"`
	IpLocation string `gorm:"type:varchar(20);comment:'IP location'" json:"ipLocation"`
	Method     string `gorm:"type:varchar(20);comment:'Request method'" json:"method"`
	Path       string `gorm:"type:varchar(100);comment:'Access path'" json:"path"`
	Remark     string `gorm:"type:varchar(100);comment:'Remark'" json:"remark"`
	Status     int    `gorm:"type:int(4);comment:'Response status code'" json:"status"`
	StartTime  string `gorm:"type:varchar(2048);comment:'Start time'" json:"startTime"`
	TimeCost   int64  `gorm:"type:int(6);comment:'Request duration (ms)'" json:"timeCost"`
	UserAgent  string `gorm:"type:varchar(2048);comment:'Browser identifier'" json:"userAgent"`
}
