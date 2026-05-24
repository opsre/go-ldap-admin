package logic

import (
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
)

type MenuLogic struct{}

// Add 添加数据
func (l MenuLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.MenuAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if isql.Menu.Exist(tools.H{"name": r.Name}) {
		return nil, tools.NewMySqlI18nError("menu.name_exists", nil)

	}

	// 获取当前用户
	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("legacy.common.current_user_failed", nil)
	}

	menu := model.Menu{
		Name:       r.Name,
		Title:      r.Title,
		Icon:       r.Icon,
		Path:       r.Path,
		Redirect:   r.Redirect,
		Component:  r.Component,
		Sort:       r.Sort,
		Status:     r.Status,
		Hidden:     r.Hidden,
		NoCache:    r.NoCache,
		AlwaysShow: r.AlwaysShow,
		Breadcrumb: r.Breadcrumb,
		ActiveMenu: r.ActiveMenu,
		ParentId:   r.ParentId,
		Creator:    ctxUser.Username,
	}

	err = isql.Menu.Add(&menu)
	if err != nil {
		return nil, tools.NewMySqlI18nError("menu.create_failed", i18n.Args{"error": err.Error()})
	}

	return nil, nil
}

// Update 更新数据
func (l MenuLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.MenuUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": int(r.ID)}
	if !isql.Menu.Exist(filter) {
		return nil, tools.NewMySqlI18nError("menu.record_not_found", nil)
	}

	// 获取当前登陆用户
	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("legacy.common.current_user_failed", nil)
	}

	oldData := new(model.Menu)
	err = isql.Menu.Find(filter, oldData)
	if err != nil {
		return nil, tools.NewMySqlI18nError("menu.record_get_failed", i18n.Args{"error": err.Error()})
	}

	menu := model.Menu{
		Model:      oldData.Model,
		Name:       r.Name,
		Title:      r.Title,
		Icon:       r.Icon,
		Path:       r.Path,
		Redirect:   r.Redirect,
		Component:  r.Component,
		Sort:       r.Sort,
		Status:     r.Status,
		Hidden:     r.Hidden,
		NoCache:    r.NoCache,
		AlwaysShow: r.AlwaysShow,
		Breadcrumb: r.Breadcrumb,
		ActiveMenu: r.ActiveMenu,
		ParentId:   r.ParentId,
		Creator:    ctxUser.Username,
	}

	err = isql.Menu.Update(&menu)
	if err != nil {
		return nil, tools.NewMySqlI18nError("menu.update_failed", i18n.Args{"error": err.Error()})
	}

	return nil, nil
}

// Delete 删除数据
func (l MenuLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.MenuDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.MenuIds {
		filter := tools.H{"id": int(id)}
		if !isql.Menu.Exist(filter) {
			return nil, tools.NewMySqlI18nError("menu.record_not_found", nil)
		}
	}

	// 删除接口
	err := isql.Menu.Delete(r.MenuIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("menu.delete_failed", i18n.Args{"error": err.Error()})
	}
	return nil, nil
}

// GetTree 数据树
func (l MenuLogic) GetTree(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.MenuGetTreeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	menus, err := isql.Menu.List()
	if err != nil {
		return nil, tools.NewMySqlI18nError("legacy.common.resource_list_failed", i18n.Args{"error": err.Error()})
	}

	tree := isql.GenMenuTree(0, menus)
	localizeMenus(c, tree)

	return tree, nil
}

// GetAccessTree 获取用户菜单树
func (l MenuLogic) GetAccessTree(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.MenuGetAccessTreeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	// 校验
	filter := tools.H{"id": r.ID}
	if !isql.User.Exist(filter) {
		return nil, tools.NewValidatorI18nError("legacy.user.not_found", nil)
	}
	user := new(model.User)
	err := isql.User.Find(filter, user)
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.mysql_query_failed", i18n.Args{"error": err.Error()})
	}
	var roleIds []uint
	for _, role := range user.Roles {
		roleIds = append(roleIds, role.ID)
	}
	menus, err := isql.Menu.ListUserMenus(roleIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("legacy.common.resource_list_failed", i18n.Args{"error": err.Error()})
	}

	tree := isql.GenMenuTree(0, menus)
	localizeMenus(c, tree)

	return tree, nil
}
