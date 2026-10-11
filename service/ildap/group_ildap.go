package ildap

import (
	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"

	ldap "github.com/go-ldap/ldap/v3"
)

type GroupService struct{}

// Add a resource
func (x GroupService) Add(g *model.Group) error { //organizationalUnit
	if g.Remark == "" {
		g.Remark = g.GroupName
	}
	add := ldap.NewAddRequest(g.GroupDN, nil)
	if g.GroupType == "ou" {
		add.Attribute("objectClass", []string{"organizationalUnit", "top"}) // If groupOfNames is defined, member must be specified, or this error occurs: object class 'groupOfNames' requires attribute 'member'
	}
	if g.GroupType == "cn" {
		add.Attribute("objectClass", []string{"groupOfUniqueNames", "top"})
		add.Attribute("uniqueMember", []string{config.Conf.Ldap.AdminDN}) // Therefore, add admin to the group by default to avoid the error above when creating an empty group
	}
	add.Attribute(g.GroupType, []string{g.GroupName})
	add.Attribute("description", []string{g.Remark})

	// Get an LDAP connection
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	return conn.Add(add)
}

// UpdateGroup updates a group
func (x GroupService) Update(oldGroup, newGroup *model.Group) error {
	modify1 := ldap.NewModifyRequest(oldGroup.GroupDN, nil)
	modify1.Replace("description", []string{newGroup.Remark})

	// Get an LDAP connection
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	err = conn.Modify(modify1)
	if err != nil {
		return err
	}
	// Rename the group if allowed by configuration and its name has changed
	if config.Conf.Ldap.GroupNameModify && newGroup.GroupName != oldGroup.GroupName {
		modify2 := ldap.NewModifyDNRequest(oldGroup.GroupDN, newGroup.GroupDN, true, "")
		err := conn.ModifyDN(modify2)
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete deletes a resource
func (x GroupService) Delete(gdn string) error {
	del := ldap.NewDelRequest(gdn, nil)

	// Get an LDAP connection
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	return conn.Del(del)
}

// Add a user to a group
func (x GroupService) AddUserToGroup(dn, udn string) error {
	//Check whether the DN starts with ou
	if dn[:3] == "ou=" {
		return tools.NewLdapI18nError("group.ou_cannot_add_user", nil)
	}
	newmr := ldap.NewModifyRequest(dn, nil)
	newmr.Add("uniqueMember", []string{udn})

	// Get an LDAP connection
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	return conn.Modify(newmr)
}

// Remove a user from a group
func (x GroupService) RemoveUserFromGroup(gdn, udn string) error {
	newmr := ldap.NewModifyRequest(gdn, nil)
	newmr.Delete("uniqueMember", []string{udn})

	// Get an LDAP connection
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	return conn.Modify(newmr)
}

// Remove a user from a group
func (x GroupService) ListGroupDN() (groups []*model.Group, err error) {
	// Construct query request
	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,                                     // This is basedn, we will start searching from this node.
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, // Here several parameters are respectively scope, derefAliases, sizeLimit, timeLimit,  typesOnly
		"(|(objectClass=organizationalUnit)(objectClass=groupOfUniqueNames))", // This is Filter for LDAP query
		[]string{"DN"}, // Here are the attributes returned by the query, provided as an array. If empty, all attributes are returned
		nil,
	)

	// Get an LDAP connection
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
