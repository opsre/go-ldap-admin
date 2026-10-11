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

// Init initializes the connection
func InitLDAP() {
	if ldapInit {
		return
	}

	ldapInitOne.Do(func() {
		ldapInit = true
	})

	// Dial takes two parameters, network and address, and returns (*Conn, error)
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

	// Assign the global variable
	ldapPool = &LdapConnPool{
		conns:    make([]*ldap.Conn, 0),
		reqConns: make(map[uint64]chan *ldap.Conn),
		openConn: 0,
		maxOpen:  config.Conf.Ldap.MaxConn,
	}
	PutLADPConn(ldapConn)

	// Hide passwords
	showDsn := fmt.Sprintf(
		"%s:******@tcp(%s)",
		config.Conf.Ldap.AdminDN,
		config.Conf.Ldap.Url,
	)

	Log.Info("LDAP initialized, DSN: ", showDsn)
}

// GetLDAPConn gets an LDAP connection
func GetLDAPConn() (*ldap.Conn, error) {
	return ldapPool.GetConnection()
}

// PutLDAPConn returns an LDAP connection to the pool
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

// Get an LDAP connection
func (lcp *LdapConnPool) GetConnection() (*ldap.Conn, error) {
	lcp.mu.Lock()
	// Check whether the connection pool has a connection
	connNum := len(lcp.conns)
	if connNum > 0 {
		lcp.openConn++
		conn := lcp.conns[0]
		copy(lcp.conns, lcp.conns[1:])
		lcp.conns = lcp.conns[:connNum-1]

		lcp.mu.Unlock()
		// If the connection is closed, get a new one
		if conn.IsClosing() {
			return initLDAPConn()
		}
		return conn, nil
	}

	// If the pool is empty and the current count exceeds the maximum connection limit
	if lcp.maxOpen != 0 && lcp.openConn > lcp.maxOpen {
		// Create a wait queue
		req := make(chan *ldap.Conn, 1)
		reqKey := lcp.nextRequestKeyLocked()
		lcp.reqConns[reqKey] = req
		lcp.mu.Unlock()

		// Wait for a connection to be returned
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

	// First check whether a wait queue exists
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

// Get the next request token
func (lcp *LdapConnPool) nextRequestKeyLocked() uint64 {
	for {
		reqKey := rand.Uint64()
		if _, ok := lcp.reqConns[reqKey]; !ok {
			return reqKey
		}
	}
}

// Get an LDAP connection
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
