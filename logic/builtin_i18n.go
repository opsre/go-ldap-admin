package logic

import (
	"strings"

	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/i18n"

	"github.com/gin-gonic/gin"
)

func localizeMenu(c *gin.Context, menu *model.Menu) {
	if menu == nil {
		return
	}
	menu.NameDisplay = localizeBuiltin(c, "builtin.menu."+menu.Name+".name", menu.Name)
	menu.TitleDisplay = localizeBuiltin(c, "builtin.menu."+menu.Name+".title", menu.Title)
	menu.CreatorDisplay = localizeBuiltinValue(c, "creator", menu.Creator)
	for _, child := range menu.Children {
		localizeMenu(c, child)
	}
}

func localizeMenus(c *gin.Context, menus []*model.Menu) {
	for _, menu := range menus {
		localizeMenu(c, menu)
	}
}

func localizeApi(c *gin.Context, api *model.Api) {
	if api == nil {
		return
	}
	api.CategoryDisplay = localizeBuiltinValue(c, "api.category", api.Category)
	api.RemarkDisplay = localizeBuiltin(c, "builtin.api.remark."+apiRemarkKey(api.Method, api.Path), api.Remark)
	api.CreatorDisplay = localizeBuiltinValue(c, "creator", api.Creator)
}

func localizeApis(c *gin.Context, apis []*model.Api) {
	for _, api := range apis {
		localizeApi(c, api)
	}
}

func localizeRole(c *gin.Context, role *model.Role) {
	if role == nil {
		return
	}
	role.NameDisplay = localizeBuiltin(c, "builtin.role."+role.Keyword+".name", role.Name)
	role.RemarkDisplay = localizeBuiltin(c, "builtin.role."+role.Keyword+".remark", role.Remark)
	role.CreatorDisplay = localizeBuiltinValue(c, "creator", role.Creator)
	for _, menu := range role.Menus {
		localizeMenu(c, menu)
	}
}

func localizeRoles(c *gin.Context, roles []*model.Role) {
	for _, role := range roles {
		localizeRole(c, role)
	}
}

func localizeGroup(c *gin.Context, group *model.Group) {
	if group == nil {
		return
	}
	group.GroupNameDisplay = localizeBuiltin(c, "builtin.group.name."+group.Source+"."+group.GroupName, group.GroupName)
	group.RemarkDisplay = localizeBuiltin(c, "builtin.group.remark."+group.Source+"."+group.GroupName, group.Remark)
	group.CreatorDisplay = localizeBuiltinValue(c, "creator", group.Creator)
	for _, child := range group.Children {
		localizeGroup(c, child)
	}
}

func localizeGroups(c *gin.Context, groups []*model.Group) {
	for _, group := range groups {
		localizeGroup(c, group)
	}
}

func localizeBuiltinValue(c *gin.Context, group string, value string) string {
	if value == "" {
		return ""
	}
	return localizeBuiltin(c, "builtin."+group+"."+value, value)
}

func localizeBuiltin(c *gin.Context, key string, fallback string) string {
	if key == "" {
		return fallback
	}
	msg := i18n.TC(c, key, nil)
	if msg == key {
		return fallback
	}
	return msg
}

func apiRemarkKey(method string, path string) string {
	key := strings.ToLower(strings.TrimSpace(method)) + "." + strings.Trim(strings.TrimSpace(path), "/")
	replacer := strings.NewReplacer("/", "_", "-", "_")
	key = replacer.Replace(key)
	return key
}
