package common

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// 全局CasbinEnforcer
var CasbinEnforcer *casbin.Enforcer

// 初始化casbin策略管理器
func InitCasbinEnforcer() {
	e, err := mysqlCasbin()
	if err != nil {
		Log.Panicf("Casbin initialization failed：%v", err)
		panic(fmt.Sprintf("Casbin initialization failed：%v", err))
	}

	CasbinEnforcer = e
	Log.Info("Casbin initialized")
}

var casbinModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`

func mysqlCasbin() (*casbin.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		return nil, err
	}
	m, err := model.NewModelFromString(casbinModel)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return e, nil
}
