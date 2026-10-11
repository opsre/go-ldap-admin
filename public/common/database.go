package common

import (
	"fmt"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库对象
var DB *gorm.DB

// 初始化数据库
func InitDB() {
	switch config.Conf.Database.Driver {
	case "mysql":
		DB = ConnMysql()
	case "sqlite3":
		DB = ConnSqlite()
	}
	dbAutoMigrate()
}

// 自动迁移表结构
func dbAutoMigrate() {
	_ = DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Group{},
		&model.Menu{},
		&model.Api{},
		&model.OperationLog{},
		&model.FieldRelation{},
	)
}

func ConnSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.Conf.Database.Source), &gorm.Config{
		// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		Log.Panicf("failed to connect sqlite3: %v", err)
	}
	dbObj, err := db.DB()
	if err != nil {
		Log.Panicf("failed to get sqlite3 obj: %v", err)
	}
	// 参见： https://github.com/glebarez/sqlite/issues/52
	dbObj.SetMaxOpenConns(1)
	return db
}

func ConnMysql() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		config.Conf.Mysql.Username,
		config.Conf.Mysql.Password,
		config.Conf.Mysql.Host,
		config.Conf.Mysql.Port,
		config.Conf.Mysql.Database,
		config.Conf.Mysql.Charset,
		config.Conf.Mysql.Collation,
		config.Conf.Mysql.Query,
	)
	// 隐藏密码
	showDsn := fmt.Sprintf(
		"%s:******@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		config.Conf.Mysql.Username,
		config.Conf.Mysql.Host,
		config.Conf.Mysql.Port,
		config.Conf.Mysql.Database,
		config.Conf.Mysql.Charset,
		config.Conf.Mysql.Collation,
		config.Conf.Mysql.Query,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		Log.Panicf("MySQL database initialization error: %v", err)
	}
	// 开启mysql日志
	if config.Conf.Mysql.LogMode {
		db.Debug()
	}
	Log.Infof("MySQL database initialized, DSN: %s", showDsn)
	return db
}
