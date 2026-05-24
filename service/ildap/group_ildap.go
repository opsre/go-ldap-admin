package ildap

import (
	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"

	ldap "github.com/go-ldap/ldap/v3"
)

type GroupService struct{}

// Add 添加资源
func (x GroupService) Add(g *model.Group) error { //organizationalUnit
	if g.Remark == "" {
		g.Remark = g.GroupName
	}
	add := ldap.NewAddRequest(g.GroupDN, nil)
	if g.GroupType == "ou" {
		add.Attribute("objectClass", []string{"organizationalUnit", "top"}) // 如果定义了 groupOfNAmes，那么必须指定member，否则报错如下：object class 'groupOfNames' requires attribute 'member'
	}
	if g.GroupType == "cn" {
		add.Attribute("objectClass", []string{"groupOfUniqueNames", "top"})
		add.Attribute("uniqueMember", []string{config.Conf.Ldap.AdminDN}) // 所以这里创建组的时候，默认将admin加入其中，以免创建时没有人员而报上边的错误
	}
	add.Attribute(g.GroupType, []string{g.GroupName})
	add.Attribute("description", []string{g.Remark})

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	return conn.Add(add)
}

// UpdateGroup 更新一个分组
func (x GroupService) Update(oldGroup, newGroup *model.Group) error {
	modify1 := ldap.NewModifyRequest(oldGroup.GroupDN, nil)
	modify1.Replace("description", []string{newGroup.Remark})

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	err = conn.Modify(modify1)
	if err != nil {
		return err
	}
	// 如果配置文件允许修改分组名称，且分组名称发生了变化，那么执行修改分组名称
	if config.Conf.Ldap.GroupNameModify && newGroup.GroupName != oldGroup.GroupName {
		modify2 := ldap.NewModifyDNRequest(oldGroup.GroupDN, newGroup.GroupDN, true, "")
		err := conn.ModifyDN(modify2)
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete 删除资源
func (x GroupService) Delete(gdn string) error {
	del := ldap.NewDelRequest(gdn, nil)

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	return conn.Del(del)
}

// AddUserToGroup 添加用户到分组
func (x GroupService) AddUserToGroup(dn, udn string) error {
	//判断dn是否以ou开头
	if dn[:3] == "ou=" {
		return tools.NewLdapI18nError("group.ou_cannot_add_user", nil)
	}
	newmr := ldap.NewModifyRequest(dn, nil)
	newmr.Add("uniqueMember", []string{udn})

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	return conn.Modify(newmr)
}

// DelUserFromGroup 将用户从分组删除
func (x GroupService) RemoveUserFromGroup(gdn, udn string) error {
	newmr := ldap.NewModifyRequest(gdn, nil)
	newmr.Delete("uniqueMember", []string{udn})

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	return conn.Modify(newmr)
}

// DelUserFromGroup 将用户从分组删除
func (x GroupService) ListGroupDN() (groups []*model.Group, err error) {
	// Construct query request
	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,                                     // This is basedn, we will start searching from this node.
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, // Here several parameters are respectively scope, derefAliases, sizeLimit, timeLimit,  typesOnly
		"(|(objectClass=organizationalUnit)(objectClass=groupOfUniqueNames))", // This is Filter for LDAP query
		[]string{"DN"}, // Here are the attributes returned by the query, provided as an array. If empty, all attributes are returned
		nil,
	)

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return groups, err
	}
	var sr *ldap.SearchResult
	// Search through ldap built-in search
	sr, err = conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	if len(sr.Entries) > 0 {
		for _, v := range sr.Entries {
			groups = append(groups, &model.Group{
				GroupDN: v.DN,
			})
		}
	}
	return
}
