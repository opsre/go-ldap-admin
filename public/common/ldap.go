package common

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/eryajf/go-ldap-admin/config"

	ldap "github.com/go-ldap/ldap/v3"
)

var ldapPool *LdapConnPool
var ldapInit = false
var ldapInitOne sync.Once

// Init 初始化连接
func InitLDAP() {
	if ldapInit {
		return
	}

	ldapInitOne.Do(func() {
		ldapInit = true
	})

	// Dail有两个参数 network,  address, 返回 (*Conn,  error)
	ldapConn, err := ldap.DialURL(config.Conf.Ldap.Url, ldap.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}))
	if err != nil {
		Log.Panicf("LDAP connection initialization error: %v", err)
		panic(fmt.Errorf("LDAP connection initialization error: %v", err))
	}
	err = ldapConn.Bind(config.Conf.Ldap.AdminDN, config.Conf.Ldap.AdminPass)
	if err != nil {
		Log.Panicf("admin account bind error: %v", err)
		panic(fmt.Errorf("admin account bind error: %v", err))
	}

	// 全局变量赋值
	ldapPool = &LdapConnPool{
		conns:    make([]*ldap.Conn, 0),
		reqConns: make(map[uint64]chan *ldap.Conn),
		openConn: 0,
		maxOpen:  config.Conf.Ldap.MaxConn,
	}
	PutLADPConn(ldapConn)

	// 隐藏密码
	showDsn := fmt.Sprintf(
		"%s:******@tcp(%s)",
		config.Conf.Ldap.AdminDN,
		config.Conf.Ldap.Url,
	)

	Log.Info("LDAP initialized, DSN: ", showDsn)
}

// GetLDAPConn 获取 LDAP 连接
func GetLDAPConn() (*ldap.Conn, error) {
	return ldapPool.GetConnection()
}

// PutLDAPConn 放回 LDAP 连接
func PutLADPConn(conn *ldap.Conn) {
	ldapPool.PutConnection(conn)
}

type LdapConnPool struct {
	mu       sync.Mutex
	conns    []*ldap.Conn
	reqConns map[uint64]chan *ldap.Conn
	openConn int
	maxOpen  int
}

// 获取一个 ladp Conn
func (lcp *LdapConnPool) GetConnection() (*ldap.Conn, error) {
	lcp.mu.Lock()
	// 判断当前连接池内是否存在连接
	connNum := len(lcp.conns)
	if connNum > 0 {
		lcp.openConn++
		conn := lcp.conns[0]
		copy(lcp.conns, lcp.conns[1:])
		lcp.conns = lcp.conns[:connNum-1]

		lcp.mu.Unlock()
		// 发现连接已经 close 重新获取连接
		if conn.IsClosing() {
			return initLDAPConn()
		}
		return conn, nil
	}

	// 当现有连接池为空时，并且当前超过最大连接限制
	if lcp.maxOpen != 0 && lcp.openConn > lcp.maxOpen {
		// 创建一个等待队列
		req := make(chan *ldap.Conn, 1)
		reqKey := lcp.nextRequestKeyLocked()
		lcp.reqConns[reqKey] = req
		lcp.mu.Unlock()

		// 等待请求归还
		return <-req, nil
	} else {
		lcp.openConn++
		lcp.mu.Unlock()
		return initLDAPConn()
	}
}

func (lcp *LdapConnPool) PutConnection(conn *ldap.Conn) {
	log.Println("Returned an LDAP connection")
	lcp.mu.Lock()
	defer lcp.mu.Unlock()

	// 先判断是否存在等待的队列
	if num := len(lcp.reqConns); num > 0 {
		var req chan *ldap.Conn
		var reqKey uint64
		for reqKey, req = range lcp.reqConns {
			break
		}
		delete(lcp.reqConns, reqKey)
		req <- conn
		return
	} else {
		lcp.openConn--
		if !conn.IsClosing() {
			lcp.conns = append(lcp.conns, conn)
		}
	}
}

// 获取下一个请求令牌
func (lcp *LdapConnPool) nextRequestKeyLocked() uint64 {
	for {
		reqKey := rand.Uint64()
		if _, ok := lcp.reqConns[reqKey]; !ok {
			return reqKey
		}
	}
}

// 获取 ladp 连接
func initLDAPConn() (*ldap.Conn, error) {
	ldap, err := ldap.DialURL(config.Conf.Ldap.Url, ldap.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}))
	if err != nil {
		return nil, err
	}
	err = ldap.Bind(config.Conf.Ldap.AdminDN, config.Conf.Ldap.AdminPass)
	if err != nil {
		return nil, err
	}
	return ldap, err
}
