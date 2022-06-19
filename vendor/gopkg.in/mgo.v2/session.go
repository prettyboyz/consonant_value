
// mgo - MongoDB driver for Go
//
// Copyright (c) 2010-2012 - Gustavo Niemeyer <gustavo@niemeyer.net>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package mgo

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Mode int

const (
	// Relevant documentation on read preference modes:
	//
	//     http://docs.mongodb.org/manual/reference/read-preference/
	//
	Primary            Mode = 2 // Default mode. All operations read from the current replica set primary.
	PrimaryPreferred   Mode = 3 // Read from the primary if available. Read from the secondary otherwise.
	Secondary          Mode = 4 // Read from one of the nearest secondary members of the replica set.
	SecondaryPreferred Mode = 5 // Read from one of the nearest secondaries if available. Read from primary otherwise.
	Nearest            Mode = 6 // Read from one of the nearest members, irrespective of it being primary or secondary.

	// Read preference modes are specific to mgo:
	Eventual  Mode = 0 // Same as Nearest, but may change servers between reads.
	Monotonic Mode = 1 // Same as SecondaryPreferred before first write. Same as Primary after first write.
	Strong    Mode = 2 // Same as Primary.
)

// mgo.v3: Drop Strong mode, suffix all modes with "Mode".

// When changing the Session type, check if newSession and copySession
// need to be updated too.

// Session represents a communication session with the database.
//
// All Session methods are concurrency-safe and may be called from multiple
// goroutines. In all session modes but Eventual, using the session from
// multiple goroutines will cause them to share the same underlying socket.
// See the documentation on Session.SetMode for more details.
type Session struct {
	m                sync.RWMutex
	cluster_         *mongoCluster
	slaveSocket      *mongoSocket
	masterSocket     *mongoSocket
	slaveOk          bool
	consistency      Mode
	queryConfig      query
	safeOp           *queryOp
	syncTimeout      time.Duration
	sockTimeout      time.Duration
	defaultdb        string
	sourcedb         string
	dialCred         *Credential
	creds            []Credential
	poolLimit        int
	bypassValidation bool
}

type Database struct {
	Session *Session
	Name    string
}

type Collection struct {
	Database *Database
	Name     string // "collection"
	FullName string // "db.collection"
}

type Query struct {
	m       sync.Mutex
	session *Session
	query   // Enables default settings in session.
}

type query struct {
	op       queryOp
	prefetch float64
	limit    int32
}

type getLastError struct {
	CmdName  int         "getLastError,omitempty"
	W        interface{} "w,omitempty"
	WTimeout int         "wtimeout,omitempty"
	FSync    bool        "fsync,omitempty"
	J        bool        "j,omitempty"
}

type Iter struct {
	m              sync.Mutex
	gotReply       sync.Cond
	session        *Session
	server         *mongoServer
	docData        queue
	err            error
	op             getMoreOp
	prefetch       float64
	limit          int32
	docsToReceive  int
	docsBeforeMore int
	timeout        time.Duration
	timedout       bool
	findCmd        bool
}

var (
	ErrNotFound = errors.New("not found")
	ErrCursor   = errors.New("invalid cursor")
)

const (
	defaultPrefetch  = 0.25
	maxUpsertRetries = 5
)

// Dial establishes a new session to the cluster identified by the given seed
// server(s). The session will enable communication with all of the servers in
// the cluster, so the seed servers are used only to find out about the cluster
// topology.
//
// Dial will timeout after 10 seconds if a server isn't reached. The returned
// session will timeout operations after one minute by default if servers
// aren't available. To customize the timeout, see DialWithTimeout,
// SetSyncTimeout, and SetSocketTimeout.
//
// This method is generally called just once for a given cluster.  Further
// sessions to the same cluster are then established using the New or Copy
// methods on the obtained session. This will make them share the underlying
// cluster, and manage the pool of connections appropriately.
//
// Once the session is not useful anymore, Close must be called to release the
// resources appropriately.
//
// The seed servers must be provided in the following format:
//
//     [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
//
// For example, it may be as simple as:
//
//     localhost
//
// Or more involved like:
//
//     mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
//
// If the port number is not provided for a server, it defaults to 27017.
//
// The username and password provided in the URL will be used to authenticate
// into the database named after the slash at the end of the host names, or
// into the "admin" database if none is provided.  The authentication information
// will persist in sessions obtained through the New method as well.
//
// The following connection options are supported after the question mark:
//
//     connect=direct
//
//         Disables the automatic replica set server discovery logic, and
//         forces the use of servers provided only (even if secondaries).
//         Note that to talk to a secondary the consistency requirements
//         must be relaxed to Monotonic or Eventual via SetMode.
//
//
//     connect=replicaSet
//
//  	   Discover replica sets automatically. Default connection behavior.
//
//
//     replicaSet=<setname>
//
//         If specified will prevent the obtained session from communicating
//         with any server which is not part of a replica set with the given name.
//         The default is to communicate with any server specified or discovered
//         via the servers contacted.
//
//
//     authSource=<db>
//
//         Informs the database used to establish credentials and privileges
//         with a MongoDB server. Defaults to the database name provided via
//         the URL path, and "admin" if that's unset.
//
//
//     authMechanism=<mechanism>
//
//        Defines the protocol for credential negotiation. Defaults to "MONGODB-CR",
//        which is the default username/password challenge-response mechanism.
//
//
//     gssapiServiceName=<name>
//
//        Defines the service name to use when authenticating with the GSSAPI
//        mechanism. Defaults to "mongodb".
//
//
//     maxPoolSize=<limit>
//
//        Defines the per-server socket pool limit. Defaults to 4096.
//        See Session.SetPoolLimit for details.
//
//
// Relevant documentation:
//
//     http://docs.mongodb.org/manual/reference/connection-string/
//
func Dial(url string) (*Session, error) {
	session, err := DialWithTimeout(url, 10*time.Second)
	if err == nil {
		session.SetSyncTimeout(1 * time.Minute)
		session.SetSocketTimeout(1 * time.Minute)
	}
	return session, err
}

// DialWithTimeout works like Dial, but uses timeout as the amount of time to
// wait for a server to respond when first connecting and also on follow up
// operations in the session. If timeout is zero, the call may block
// forever waiting for a connection to be made.
//
// See SetSyncTimeout for customizing the timeout for the session.
func DialWithTimeout(url string, timeout time.Duration) (*Session, error) {
	info, err := ParseURL(url)
	if err != nil {
		return nil, err
	}
	info.Timeout = timeout
	return DialWithInfo(info)
}

// ParseURL parses a MongoDB URL as accepted by the Dial function and returns
// a value suitable for providing into DialWithInfo.
//
// See Dial for more details on the format of url.
func ParseURL(url string) (*DialInfo, error) {
	uinfo, err := extractURL(url)
	if err != nil {
		return nil, err
	}
	direct := false
	mechanism := ""
	service := ""
	source := ""
	setName := ""
	poolLimit := 0
	for k, v := range uinfo.options {
		switch k {
		case "authSource":
			source = v
		case "authMechanism":
			mechanism = v
		case "gssapiServiceName":
			service = v
		case "replicaSet":
			setName = v
		case "maxPoolSize":
			poolLimit, err = strconv.Atoi(v)
			if err != nil {
				return nil, errors.New("bad value for maxPoolSize: " + v)
			}
		case "connect":
			if v == "direct" {
				direct = true
				break
			}
			if v == "replicaSet" {
				break
			}
			fallthrough
		default:
			return nil, errors.New("unsupported connection URL option: " + k + "=" + v)
		}
	}
	info := DialInfo{
		Addrs:          uinfo.addrs,
		Direct:         direct,
		Database:       uinfo.db,
		Username:       uinfo.user,
		Password:       uinfo.pass,
		Mechanism:      mechanism,
		Service:        service,
		Source:         source,
		PoolLimit:      poolLimit,
		ReplicaSetName: setName,
	}
	return &info, nil
}

// DialInfo holds options for establishing a session with a MongoDB cluster.
// To use a URL, see the Dial function.
type DialInfo struct {
	// Addrs holds the addresses for the seed servers.
	Addrs []string

	// Direct informs whether to establish connections only with the
	// specified seed servers, or to obtain information for the whole
	// cluster and establish connections with further servers too.
	Direct bool

	// Timeout is the amount of time to wait for a server to respond when
	// first connecting and on follow up operations in the session. If
	// timeout is zero, the call may block forever waiting for a connection
	// to be established. Timeout does not affect logic in DialServer.
	Timeout time.Duration

	// FailFast will cause connection and query attempts to fail faster when
	// the server is unavailable, instead of retrying until the configured
	// timeout period. Note that an unavailable server may silently drop
	// packets instead of rejecting them, in which case it's impossible to
	// distinguish it from a slow server, so the timeout stays relevant.
	FailFast bool

	// Database is the default database name used when the Session.DB method
	// is called with an empty name, and is also used during the initial
	// authentication if Source is unset.
	Database string

	// ReplicaSetName, if specified, will prevent the obtained session from
	// communicating with any server which is not part of a replica set
	// with the given name. The default is to communicate with any server
	// specified or discovered via the servers contacted.
	ReplicaSetName string

	// Source is the database used to establish credentials and privileges
	// with a MongoDB server. Defaults to the value of Database, if that is
	// set, or "admin" otherwise.
	Source string

	// Service defines the service name to use when authenticating with the GSSAPI
	// mechanism. Defaults to "mongodb".
	Service string

	// ServiceHost defines which hostname to use when authenticating
	// with the GSSAPI mechanism. If not specified, defaults to the MongoDB
	// server's address.
	ServiceHost string

	// Mechanism defines the protocol for credential negotiation.
	// Defaults to "MONGODB-CR".
	Mechanism string

	// Username and Password inform the credentials for the initial authentication
	// done on the database defined by the Source field. See Session.Login.
	Username string
	Password string

	// PoolLimit defines the per-server socket pool limit. Defaults to 4096.
	// See Session.SetPoolLimit for details.
	PoolLimit int

	// DialServer optionally specifies the dial function for establishing
	// connections with the MongoDB servers.
	DialServer func(addr *ServerAddr) (net.Conn, error)

	// WARNING: This field is obsolete. See DialServer above.
	Dial func(addr net.Addr) (net.Conn, error)
}

// mgo.v3: Drop DialInfo.Dial.

// ServerAddr represents the address for establishing a connection to an
// individual MongoDB server.
type ServerAddr struct {
	str string
	tcp *net.TCPAddr
}

// String returns the address that was provided for the server before resolution.
func (addr *ServerAddr) String() string {
	return addr.str
}

// TCPAddr returns the resolved TCP address for the server.
func (addr *ServerAddr) TCPAddr() *net.TCPAddr {
	return addr.tcp
}

// DialWithInfo establishes a new session to the cluster identified by info.
func DialWithInfo(info *DialInfo) (*Session, error) {
	addrs := make([]string, len(info.Addrs))
	for i, addr := range info.Addrs {
		p := strings.LastIndexAny(addr, "]:")
		if p == -1 || addr[p] != ':' {
			// XXX This is untested. The test suite doesn't use the standard port.
			addr += ":27017"
		}
		addrs[i] = addr
	}
	cluster := newCluster(addrs, info.Direct, info.FailFast, dialer{info.Dial, info.DialServer}, info.ReplicaSetName)
	session := newSession(Eventual, cluster, info.Timeout)
	session.defaultdb = info.Database
	if session.defaultdb == "" {
		session.defaultdb = "test"
	}
	session.sourcedb = info.Source
	if session.sourcedb == "" {
		session.sourcedb = info.Database
		if session.sourcedb == "" {
			session.sourcedb = "admin"
		}
	}
	if info.Username != "" {
		source := session.sourcedb
		if info.Source == "" &&
			(info.Mechanism == "GSSAPI" || info.Mechanism == "PLAIN" || info.Mechanism == "MONGODB-X509") {
			source = "$external"
		}
		session.dialCred = &Credential{
			Username:    info.Username,
			Password:    info.Password,
			Mechanism:   info.Mechanism,
			Service:     info.Service,
			ServiceHost: info.ServiceHost,
			Source:      source,
		}
		session.creds = []Credential{*session.dialCred}
	}
	if info.PoolLimit > 0 {
		session.poolLimit = info.PoolLimit
	}
	cluster.Release()

	// People get confused when we return a session that is not actually
	// established to any servers yet (e.g. what if url was wrong). So,
	// ping the server to ensure there's someone there, and abort if it
	// fails.
	if err := session.Ping(); err != nil {
		session.Close()
		return nil, err
	}
	session.SetMode(Strong, true)
	return session, nil
}

func isOptSep(c rune) bool {
	return c == ';' || c == '&'
}

type urlInfo struct {
	addrs   []string
	user    string
	pass    string
	db      string
	options map[string]string
}

func extractURL(s string) (*urlInfo, error) {
	if strings.HasPrefix(s, "mongodb://") {
		s = s[10:]
	}
	info := &urlInfo{options: make(map[string]string)}
	if c := strings.Index(s, "?"); c != -1 {
		for _, pair := range strings.FieldsFunc(s[c+1:], isOptSep) {
			l := strings.SplitN(pair, "=", 2)
			if len(l) != 2 || l[0] == "" || l[1] == "" {
				return nil, errors.New("connection option must be key=value: " + pair)
			}
			info.options[l[0]] = l[1]
		}
		s = s[:c]
	}
	if c := strings.Index(s, "@"); c != -1 {
		pair := strings.SplitN(s[:c], ":", 2)
		if len(pair) > 2 || pair[0] == "" {
			return nil, errors.New("credentials must be provided as user:pass@host")
		}
		var err error
		info.user, err = url.QueryUnescape(pair[0])
		if err != nil {
			return nil, fmt.Errorf("cannot unescape username in URL: %q", pair[0])
		}
		if len(pair) > 1 {
			info.pass, err = url.QueryUnescape(pair[1])
			if err != nil {
				return nil, fmt.Errorf("cannot unescape password in URL")
			}
		}
		s = s[c+1:]
	}
	if c := strings.Index(s, "/"); c != -1 {
		info.db = s[c+1:]
		s = s[:c]
	}
	info.addrs = strings.Split(s, ",")
	return info, nil
}

func newSession(consistency Mode, cluster *mongoCluster, timeout time.Duration) (session *Session) {
	cluster.Acquire()
	session = &Session{
		cluster_:    cluster,
		syncTimeout: timeout,
		sockTimeout: timeout,
		poolLimit:   4096,
	}
	debugf("New session %p on cluster %p", session, cluster)
	session.SetMode(consistency, true)
	session.SetSafe(&Safe{})
	session.queryConfig.prefetch = defaultPrefetch
	return session
}

func copySession(session *Session, keepCreds bool) (s *Session) {
	cluster := session.cluster()
	cluster.Acquire()
	if session.masterSocket != nil {
		session.masterSocket.Acquire()
	}
	if session.slaveSocket != nil {
		session.slaveSocket.Acquire()
	}
	var creds []Credential
	if keepCreds {
		creds = make([]Credential, len(session.creds))
		copy(creds, session.creds)
	} else if session.dialCred != nil {
		creds = []Credential{*session.dialCred}
	}
	scopy := *session
	scopy.m = sync.RWMutex{}
	scopy.creds = creds
	s = &scopy
	debugf("New session %p on cluster %p (copy from %p)", s, cluster, session)
	return s
}

// LiveServers returns a list of server addresses which are
// currently known to be alive.
func (s *Session) LiveServers() (addrs []string) {
	s.m.RLock()
	addrs = s.cluster().LiveServers()
	s.m.RUnlock()
	return addrs
}

// DB returns a value representing the named database. If name
// is empty, the database name provided in the dialed URL is
// used instead. If that is also empty, "test" is used as a
// fallback in a way equivalent to the mongo shell.
//
// Creating this value is a very lightweight operation, and
// involves no network communication.
func (s *Session) DB(name string) *Database {
	if name == "" {
		name = s.defaultdb
	}
	return &Database{s, name}
}

// C returns a value representing the named collection.
//
// Creating this value is a very lightweight operation, and
// involves no network communication.
func (db *Database) C(name string) *Collection {
	return &Collection{db, name, db.Name + "." + name}
}

// With returns a copy of db that uses session s.
func (db *Database) With(s *Session) *Database {
	newdb := *db
	newdb.Session = s
	return &newdb
}

// With returns a copy of c that uses session s.
func (c *Collection) With(s *Session) *Collection {
	newdb := *c.Database
	newdb.Session = s
	newc := *c
	newc.Database = &newdb
	return &newc
}

// GridFS returns a GridFS value representing collections in db that
// follow the standard GridFS specification.
// The provided prefix (sometimes known as root) will determine which
// collections to use, and is usually set to "fs" when there is a
// single GridFS in the database.
//
// See the GridFS Create, Open, and OpenId methods for more details.
//
// Relevant documentation:
//
//     http://www.mongodb.org/display/DOCS/GridFS
//     http://www.mongodb.org/display/DOCS/GridFS+Tools
//     http://www.mongodb.org/display/DOCS/GridFS+Specification
//
func (db *Database) GridFS(prefix string) *GridFS {
	return newGridFS(db, prefix)
}

// Run issues the provided command on the db database and unmarshals
// its result in the respective argument. The cmd argument may be either
// a string with the command name itself, in which case an empty document of
// the form bson.M{cmd: 1} will be used, or it may be a full command document.
//
// Note that MongoDB considers the first marshalled key as the command
// name, so when providing a command with options, it's important to
// use an ordering-preserving document, such as a struct value or an
// instance of bson.D.  For instance:
//
//     db.Run(bson.D{{"create", "mycollection"}, {"size", 1024}})
//
// For privilleged commands typically run on the "admin" database, see
// the Run method in the Session type.
//
// Relevant documentation:
//
//     http://www.mongodb.org/display/DOCS/Commands
//     http://www.mongodb.org/display/DOCS/List+of+Database+CommandSkips
//
func (db *Database) Run(cmd interface{}, result interface{}) error {
	socket, err := db.Session.acquireSocket(true)
	if err != nil {
		return err
	}
	defer socket.Release()

	// This is an optimized form of db.C("$cmd").Find(cmd).One(result).
	return db.run(socket, cmd, result)
}

// Credential holds details to authenticate with a MongoDB server.
type Credential struct {
	// Username and Password hold the basic details for authentication.
	// Password is optional with some authentication mechanisms.
	Username string
	Password string

	// Source is the database used to establish credentials and privileges
	// with a MongoDB server. Defaults to the default database provided
	// during dial, or "admin" if that was unset.
	Source string

	// Service defines the service name to use when authenticating with the GSSAPI
	// mechanism. Defaults to "mongodb".
	Service string

	// ServiceHost defines which hostname to use when authenticating
	// with the GSSAPI mechanism. If not specified, defaults to the MongoDB
	// server's address.
	ServiceHost string

	// Mechanism defines the protocol for credential negotiation.
	// Defaults to "MONGODB-CR".
	Mechanism string
}

// Login authenticates with MongoDB using the provided credential.  The
// authentication is valid for the whole session and will stay valid until
// Logout is explicitly called for the same database, or the session is
// closed.
func (db *Database) Login(user, pass string) error {
	return db.Session.Login(&Credential{Username: user, Password: pass, Source: db.Name})
}

// Login authenticates with MongoDB using the provided credential.  The
// authentication is valid for the whole session and will stay valid until
// Logout is explicitly called for the same database, or the session is
// closed.
func (s *Session) Login(cred *Credential) error {
	socket, err := s.acquireSocket(true)
	if err != nil {
		return err
	}
	defer socket.Release()

	credCopy := *cred
	if cred.Source == "" {
		if cred.Mechanism == "GSSAPI" {
			credCopy.Source = "$external"
		} else {
			credCopy.Source = s.sourcedb
		}
	}
	err = socket.Login(credCopy)
	if err != nil {
		return err
	}

	s.m.Lock()
	s.creds = append(s.creds, credCopy)
	s.m.Unlock()
	return nil
}

func (s *Session) socketLogin(socket *mongoSocket) error {
	for _, cred := range s.creds {
		if err := socket.Login(cred); err != nil {
			return err
		}
	}
	return nil
}

// Logout removes any established authentication credentials for the database.
func (db *Database) Logout() {
	session := db.Session
	dbname := db.Name
	session.m.Lock()
	found := false
	for i, cred := range session.creds {
		if cred.Source == dbname {
			copy(session.creds[i:], session.creds[i+1:])
			session.creds = session.creds[:len(session.creds)-1]
			found = true
			break
		}
	}
	if found {
		if session.masterSocket != nil {
			session.masterSocket.Logout(dbname)
		}
		if session.slaveSocket != nil {
			session.slaveSocket.Logout(dbname)
		}
	}
	session.m.Unlock()
}

// LogoutAll removes all established authentication credentials for the session.
func (s *Session) LogoutAll() {
	s.m.Lock()
	for _, cred := range s.creds {
		if s.masterSocket != nil {
			s.masterSocket.Logout(cred.Source)
		}
		if s.slaveSocket != nil {
			s.slaveSocket.Logout(cred.Source)
		}
	}
	s.creds = s.creds[0:0]
	s.m.Unlock()
}

// User represents a MongoDB user.
//
// Relevant documentation:
//
//     http://docs.mongodb.org/manual/reference/privilege-documents/
//     http://docs.mongodb.org/manual/reference/user-privileges/
//
type User struct {
	// Username is how the user identifies itself to the system.
	Username string `bson:"user"`

	// Password is the plaintext password for the user. If set,
	// the UpsertUser method will hash it into PasswordHash and
	// unset it before the user is added to the database.
	Password string `bson:",omitempty"`

	// PasswordHash is the MD5 hash of Username+":mongo:"+Password.
	PasswordHash string `bson:"pwd,omitempty"`

	// CustomData holds arbitrary data admins decide to associate
	// with this user, such as the full name or employee id.
	CustomData interface{} `bson:"customData,omitempty"`

	// Roles indicates the set of roles the user will be provided.
	// See the Role constants.
	Roles []Role `bson:"roles"`

	// OtherDBRoles allows assigning roles in other databases from
	// user documents inserted in the admin database. This field
	// only works in the admin database.
	OtherDBRoles map[string][]Role `bson:"otherDBRoles,omitempty"`

	// UserSource indicates where to look for this user's credentials.
	// It may be set to a database name, or to "$external" for
	// consulting an external resource such as Kerberos. UserSource
	// must not be set if Password or PasswordHash are present.
	//
	// WARNING: This setting was only ever supported in MongoDB 2.4,
	// and is now obsolete.
	UserSource string `bson:"userSource,omitempty"`
}

type Role string

const (
	// Relevant documentation:
	//
	//     http://docs.mongodb.org/manual/reference/user-privileges/
	//
	RoleRoot         Role = "root"
	RoleRead         Role = "read"
	RoleReadAny      Role = "readAnyDatabase"
	RoleReadWrite    Role = "readWrite"
	RoleReadWriteAny Role = "readWriteAnyDatabase"
	RoleDBAdmin      Role = "dbAdmin"
	RoleDBAdminAny   Role = "dbAdminAnyDatabase"
	RoleUserAdmin    Role = "userAdmin"
	RoleUserAdminAny Role = "userAdminAnyDatabase"
	RoleClusterAdmin Role = "clusterAdmin"
)

// UpsertUser updates the authentication credentials and the roles for
// a MongoDB user within the db database. If the named user doesn't exist
// it will be created.
//
// This method should only be used from MongoDB 2.4 and on. For older
// MongoDB releases, use the obsolete AddUser method instead.
//
// Relevant documentation:
//
//     http://docs.mongodb.org/manual/reference/user-privileges/
//     http://docs.mongodb.org/manual/reference/privilege-documents/
//
func (db *Database) UpsertUser(user *User) error {
	if user.Username == "" {
		return fmt.Errorf("user has no Username")
	}
	if (user.Password != "" || user.PasswordHash != "") && user.UserSource != "" {
		return fmt.Errorf("user has both Password/PasswordHash and UserSource set")
	}
	if len(user.OtherDBRoles) > 0 && db.Name != "admin" && db.Name != "$external" {
		return fmt.Errorf("user with OtherDBRoles is only supported in the admin or $external databases")
	}

	// Attempt to run this using 2.6+ commands.
	rundb := db
	if user.UserSource != "" {
		// Compatibility logic for the userSource field of MongoDB <= 2.4.X
		rundb = db.Session.DB(user.UserSource)
	}
	err := rundb.runUserCmd("updateUser", user)
	// retry with createUser when isAuthError in order to enable the "localhost exception"
	if isNotFound(err) || isAuthError(err) {
		return rundb.runUserCmd("createUser", user)
	}
	if !isNoCmd(err) {
		return err
	}

	// Command does not exist. Fallback to pre-2.6 behavior.
	var set, unset bson.D
	if user.Password != "" {
		psum := md5.New()
		psum.Write([]byte(user.Username + ":mongo:" + user.Password))
		set = append(set, bson.DocElem{"pwd", hex.EncodeToString(psum.Sum(nil))})
		unset = append(unset, bson.DocElem{"userSource", 1})
	} else if user.PasswordHash != "" {
		set = append(set, bson.DocElem{"pwd", user.PasswordHash})
		unset = append(unset, bson.DocElem{"userSource", 1})
	}
	if user.UserSource != "" {
		set = append(set, bson.DocElem{"userSource", user.UserSource})
		unset = append(unset, bson.DocElem{"pwd", 1})
	}
	if user.Roles != nil || user.OtherDBRoles != nil {
		set = append(set, bson.DocElem{"roles", user.Roles})
		if len(user.OtherDBRoles) > 0 {
			set = append(set, bson.DocElem{"otherDBRoles", user.OtherDBRoles})
		} else {
			unset = append(unset, bson.DocElem{"otherDBRoles", 1})
		}
	}
	users := db.C("system.users")
	err = users.Update(bson.D{{"user", user.Username}}, bson.D{{"$unset", unset}, {"$set", set}})
	if err == ErrNotFound {
		set = append(set, bson.DocElem{"user", user.Username})
		if user.Roles == nil && user.OtherDBRoles == nil {
			// Roles must be sent, as it's the way MongoDB distinguishes
			// old-style documents from new-style documents in pre-2.6.
			set = append(set, bson.DocElem{"roles", user.Roles})
		}
		err = users.Insert(set)
	}
	return err
}

func isNoCmd(err error) bool {
	e, ok := err.(*QueryError)
	return ok && (e.Code == 59 || e.Code == 13390 || strings.HasPrefix(e.Message, "no such cmd:"))
}

func isNotFound(err error) bool {
	e, ok := err.(*QueryError)
	return ok && e.Code == 11
}

func isAuthError(err error) bool {
	e, ok := err.(*QueryError)
	return ok && e.Code == 13
}

func (db *Database) runUserCmd(cmdName string, user *User) error {
	cmd := make(bson.D, 0, 16)
	cmd = append(cmd, bson.DocElem{cmdName, user.Username})
	if user.Password != "" {
		cmd = append(cmd, bson.DocElem{"pwd", user.Password})
	}
	var roles []interface{}
	for _, role := range user.Roles {
		roles = append(roles, role)
	}
	for db, dbroles := range user.OtherDBRoles {
		for _, role := range dbroles {
			roles = append(roles, bson.D{{"role", role}, {"db", db}})
		}
	}
	if roles != nil || user.Roles != nil || cmdName == "createUser" {
		cmd = append(cmd, bson.DocElem{"roles", roles})
	}
	err := db.Run(cmd, nil)
	if !isNoCmd(err) && user.UserSource != "" && (user.UserSource != "$external" || db.Name != "$external") {
		return fmt.Errorf("MongoDB 2.6+ does not support the UserSource setting")
	}
	return err
}

// AddUser creates or updates the authentication credentials of user within
// the db database.
//
// WARNING: This method is obsolete and should only be used with MongoDB 2.2
// or earlier. For MongoDB 2.4 and on, use UpsertUser instead.
func (db *Database) AddUser(username, password string, readOnly bool) error {
	// Try to emulate the old behavior on 2.6+
	user := &User{Username: username, Password: password}
	if db.Name == "admin" {
		if readOnly {
			user.Roles = []Role{RoleReadAny}
		} else {
			user.Roles = []Role{RoleReadWriteAny}
		}
	} else {
		if readOnly {
			user.Roles = []Role{RoleRead}
		} else {
			user.Roles = []Role{RoleReadWrite}
		}
	}
	err := db.runUserCmd("updateUser", user)
	if isNotFound(err) {
		return db.runUserCmd("createUser", user)
	}
	if !isNoCmd(err) {
		return err
	}

	// Command doesn't exist. Fallback to pre-2.6 behavior.
	psum := md5.New()
	psum.Write([]byte(username + ":mongo:" + password))
	digest := hex.EncodeToString(psum.Sum(nil))
	c := db.C("system.users")
	_, err = c.Upsert(bson.M{"user": username}, bson.M{"$set": bson.M{"user": username, "pwd": digest, "readOnly": readOnly}})
	return err
}

// RemoveUser removes the authentication credentials of user from the database.
func (db *Database) RemoveUser(user string) error {
	err := db.Run(bson.D{{"dropUser", user}}, nil)
	if isNoCmd(err) {
		users := db.C("system.users")
		return users.Remove(bson.M{"user": user})
	}
	if isNotFound(err) {
		return ErrNotFound
	}
	return err
}

type indexSpec struct {
	Name, NS         string
	Key              bson.D
	Unique           bool    ",omitempty"
	DropDups         bool    "dropDups,omitempty"
	Background       bool    ",omitempty"
	Sparse           bool    ",omitempty"
	Bits             int     ",omitempty"
	Min, Max         float64 ",omitempty"
	BucketSize       float64 "bucketSize,omitempty"
	ExpireAfter      int     "expireAfterSeconds,omitempty"
	Weights          bson.D  ",omitempty"
	DefaultLanguage  string  "default_language,omitempty"
	LanguageOverride string  "language_override,omitempty"
	TextIndexVersion int     "textIndexVersion,omitempty"

	Collation *Collation "collation,omitempty"
}

type Index struct {
	Key        []string // Index key fields; prefix name with dash (-) for descending order
	Unique     bool     // Prevent two documents from having the same index key
	DropDups   bool     // Drop documents with the same index key as a previously indexed one
	Background bool     // Build index in background and return immediately
	Sparse     bool     // Only index documents containing the Key fields

	// If ExpireAfter is defined the server will periodically delete
	// documents with indexed time.Time older than the provided delta.
	ExpireAfter time.Duration

	// Name holds the stored index name. On creation if this field is unset it is
	// computed by EnsureIndex based on the index key.
	Name string

	// Properties for spatial indexes.
	//
	// Min and Max were improperly typed as int when they should have been
	// floats.  To preserve backwards compatibility they are still typed as
	// int and the following two fields enable reading and writing the same
	// fields as float numbers. In mgo.v3, these fields will be dropped and
	// Min/Max will become floats.
	Min, Max   int
	Minf, Maxf float64
	BucketSize float64
	Bits       int

	// Properties for text indexes.
	DefaultLanguage  string
	LanguageOverride string

	// Weights defines the significance of provided fields relative to other
	// fields in a text index. The score for a given word in a document is derived
	// from the weighted sum of the frequency for each of the indexed fields in
	// that document. The default field weight is 1.
	Weights map[string]int

	// Collation defines the collation to use for the index.
	Collation *Collation
}

type Collation struct {

	// Locale defines the collation locale.
	Locale string `bson:"locale"`

	// CaseLevel defines whether to turn case sensitivity on at strength 1 or 2.
	CaseLevel bool `bson:"caseLevel,omitempty"`

	// CaseFirst may be set to "upper" or "lower" to define whether
	// to have uppercase or lowercase items first. Default is "off".
	CaseFirst string `bson:"caseFirst,omitempty"`

	// Strength defines the priority of comparison properties, as follows:
	//
	//   1 (primary)    - Strongest level, denote difference between base characters
	//   2 (secondary)  - Accents in characters are considered secondary differences
	//   3 (tertiary)   - Upper and lower case differences in characters are
	//                    distinguished at the tertiary level
	//   4 (quaternary) - When punctuation is ignored at level 1-3, an additional
	//                    level can be used to distinguish words with and without
	//                    punctuation. Should only be used if ignoring punctuation
	//                    is required or when processing Japanese text.
	//   5 (identical)  - When all other levels are equal, the identical level is
	//                    used as a tiebreaker. The Unicode code point values of
	//                    the NFD form of each string are compared at this level,
	//                    just in case there is no difference at levels 1-4
	//
	// Strength defaults to 3.
	Strength int `bson:"strength,omitempty"`

	// NumericOrdering defines whether to order numbers based on numerical
	// order and not collation order.
	NumericOrdering bool `bson:"numericOrdering,omitempty"`

	// Alternate controls whether spaces and punctuation are considered base characters.
	// May be set to "non-ignorable" (spaces and punctuation considered base characters)
	// or "shifted" (spaces and punctuation not considered base characters, and only
	// distinguished at strength > 3). Defaults to "non-ignorable".
	Alternate string `bson:"alternate,omitempty"`

	// Backwards defines whether to have secondary differences considered in reverse order,
	// as done in the French language.
	Backwards bool `bson:"backwards,omitempty"`
}

// mgo.v3: Drop Minf and Maxf and transform Min and Max to floats.
// mgo.v3: Drop DropDups as it's unsupported past 2.8.

type indexKeyInfo struct {
	name    string
	key     bson.D
	weights bson.D
}

func parseIndexKey(key []string) (*indexKeyInfo, error) {
	var keyInfo indexKeyInfo
	isText := false
	var order interface{}
	for _, field := range key {
		raw := field
		if keyInfo.name != "" {
			keyInfo.name += "_"
		}
		var kind string
		if field != "" {
			if field[0] == '$' {
				if c := strings.Index(field, ":"); c > 1 && c < len(field)-1 {
					kind = field[1:c]
					field = field[c+1:]
					keyInfo.name += field + "_" + kind
				} else {
					field = "\x00"
				}
			}
			switch field[0] {
			case 0:
				// Logic above failed. Reset and error.
				field = ""
			case '@':
				order = "2d"
				field = field[1:]
				// The shell used to render this field as key_ instead of key_2d,
				// and mgo followed suit. This has been fixed in recent server
				// releases, and mgo followed as well.
				keyInfo.name += field + "_2d"
			case '-':
				order = -1
				field = field[1:]
				keyInfo.name += field + "_-1"
			case '+':
				field = field[1:]
				fallthrough
			default:
				if kind == "" {
					order = 1
					keyInfo.name += field + "_1"
				} else {
					order = kind
				}
			}
		}
		if field == "" || kind != "" && order != kind {
			return nil, fmt.Errorf(`invalid index key: want "[$<kind>:][-]<field name>", got %q`, raw)
		}
		if kind == "text" {
			if !isText {
				keyInfo.key = append(keyInfo.key, bson.DocElem{"_fts", "text"}, bson.DocElem{"_ftsx", 1})
				isText = true
			}
			keyInfo.weights = append(keyInfo.weights, bson.DocElem{field, 1})
		} else {
			keyInfo.key = append(keyInfo.key, bson.DocElem{field, order})
		}
	}
	if keyInfo.name == "" {
		return nil, errors.New("invalid index key: no fields provided")
	}
	return &keyInfo, nil
}

// EnsureIndexKey ensures an index with the given key exists, creating it
// if necessary.
//
// This example:
//
//     err := collection.EnsureIndexKey("a", "b")
//
// Is equivalent to:
//
//     err := collection.EnsureIndex(mgo.Index{Key: []string{"a", "b"}})
//
// See the EnsureIndex method for more details.
func (c *Collection) EnsureIndexKey(key ...string) error {
	return c.EnsureIndex(Index{Key: key})
}

// EnsureIndex ensures an index with the given key exists, creating it with
// the provided parameters if necessary. EnsureIndex does not modify a previously
// existent index with a matching key. The old index must be dropped first instead.
//
// Once EnsureIndex returns successfully, following requests for the same index
// will not contact the server unless Collection.DropIndex is used to drop the
// same index, or Session.ResetIndexCache is called.
//
// For example:
//
//     index := Index{
//         Key: []string{"lastname", "firstname"},
//         Unique: true,
//         DropDups: true,
//         Background: true, // See notes.
//         Sparse: true,
//     }
//     err := collection.EnsureIndex(index)
//
// The Key value determines which fields compose the index. The index ordering
// will be ascending by default.  To obtain an index with a descending order,
// the field name should be prefixed by a dash (e.g. []string{"-time"}). It can
// also be optionally prefixed by an index kind, as in "$text:summary" or
// "$2d:-point". The key string format is:
//
//     [$<kind>:][-]<field name>
//
// If the Unique field is true, the index must necessarily contain only a single
// document per Key.  With DropDups set to true, documents with the same key
// as a previously indexed one will be dropped rather than an error returned.
//
// If Background is true, other connections will be allowed to proceed using
// the collection without the index while it's being built. Note that the
// session executing EnsureIndex will be blocked for as long as it takes for
// the index to be built.
//
// If Sparse is true, only documents containing the provided Key fields will be
// included in the index.  When using a sparse index for sorting, only indexed
// documents will be returned.
//
// If ExpireAfter is non-zero, the server will periodically scan the collection
// and remove documents containing an indexed time.Time field with a value
// older than ExpireAfter. See the documentation for details:
//
//     http://docs.mongodb.org/manual/tutorial/expire-data
//
// Other kinds of indexes are also supported through that API. Here is an example:
//
//     index := Index{
//         Key: []string{"$2d:loc"},
//         Bits: 26,
//     }
//     err := collection.EnsureIndex(index)
//
// The example above requests the creation of a "2d" index for the "loc" field.
//
// The 2D index bounds may be changed using the Min and Max attributes of the
// Index value.  The default bound setting of (-180, 180) is suitable for
// latitude/longitude pairs.
//
// The Bits parameter sets the precision of the 2D geohash values.  If not
// provided, 26 bits are used, which is roughly equivalent to 1 foot of
// precision for the default (-180, 180) index bounds.
//
// Relevant documentation:
//
//     http://www.mongodb.org/display/DOCS/Indexes
//     http://www.mongodb.org/display/DOCS/Indexing+Advice+and+FAQ
//     http://www.mongodb.org/display/DOCS/Indexing+as+a+Background+Operation
//     http://www.mongodb.org/display/DOCS/Geospatial+Indexing
//     http://www.mongodb.org/display/DOCS/Multikeys
//
func (c *Collection) EnsureIndex(index Index) error {
	keyInfo, err := parseIndexKey(index.Key)
	if err != nil {
		return err
	}

	session := c.Database.Session
	cacheKey := c.FullName + "\x00" + keyInfo.name
	if session.cluster().HasCachedIndex(cacheKey) {
		return nil
	}

	spec := indexSpec{
		Name:             keyInfo.name,
		NS:               c.FullName,
		Key:              keyInfo.key,
		Unique:           index.Unique,
		DropDups:         index.DropDups,
		Background:       index.Background,
		Sparse:           index.Sparse,
		Bits:             index.Bits,
		Min:              index.Minf,
		Max:              index.Maxf,
		BucketSize:       index.BucketSize,
		ExpireAfter:      int(index.ExpireAfter / time.Second),
		Weights:          keyInfo.weights,
		DefaultLanguage:  index.DefaultLanguage,
		LanguageOverride: index.LanguageOverride,
		Collation:        index.Collation,
	}

	if spec.Min == 0 && spec.Max == 0 {
		spec.Min = float64(index.Min)
		spec.Max = float64(index.Max)
	}

	if index.Name != "" {
		spec.Name = index.Name
	}

NextField:
	for name, weight := range index.Weights {
		for i, elem := range spec.Weights {
			if elem.Name == name {
				spec.Weights[i].Value = weight
				continue NextField
			}
		}
		panic("weight provided for field that is not part of index key: " + name)
	}

	cloned := session.Clone()
	defer cloned.Close()
	cloned.SetMode(Strong, false)
	cloned.EnsureSafe(&Safe{})
	db := c.Database.With(cloned)

	// Try with a command first.
	err = db.Run(bson.D{{"createIndexes", c.Name}, {"indexes", []indexSpec{spec}}}, nil)
	if isNoCmd(err) {
		// Command not yet supported. Insert into the indexes collection instead.
		err = db.C("system.indexes").Insert(&spec)
	}
	if err == nil {
		session.cluster().CacheIndex(cacheKey, true)
	}
	return err
}

// DropIndex drops the index with the provided key from the c collection.
//
// See EnsureIndex for details on the accepted key variants.
//
// For example:
//
//     err1 := collection.DropIndex("firstField", "-secondField")
//     err2 := collection.DropIndex("customIndexName")
//
func (c *Collection) DropIndex(key ...string) error {
	keyInfo, err := parseIndexKey(key)
	if err != nil {
		return err
	}

	session := c.Database.Session
	cacheKey := c.FullName + "\x00" + keyInfo.name
	session.cluster().CacheIndex(cacheKey, false)

	session = session.Clone()
	defer session.Close()
	session.SetMode(Strong, false)

	db := c.Database.With(session)
	result := struct {
		ErrMsg string
		Ok     bool
	}{}
	err = db.Run(bson.D{{"dropIndexes", c.Name}, {"index", keyInfo.name}}, &result)
	if err != nil {
		return err
	}
	if !result.Ok {
		return errors.New(result.ErrMsg)
	}
	return nil
}

// DropIndexName removes the index with the provided index name.
//
// For example:
//
//     err := collection.DropIndex("customIndexName")
//
func (c *Collection) DropIndexName(name string) error {
	session := c.Database.Session

	session = session.Clone()
	defer session.Close()
	session.SetMode(Strong, false)

	c = c.With(session)

	indexes, err := c.Indexes()
	if err != nil {
		return err
	}

	var index Index
	for _, idx := range indexes {
		if idx.Name == name {
			index = idx
			break
		}
	}

	if index.Name != "" {
		keyInfo, err := parseIndexKey(index.Key)
		if err != nil {
			return err
		}

		cacheKey := c.FullName + "\x00" + keyInfo.name
		session.cluster().CacheIndex(cacheKey, false)
	}

	result := struct {
		ErrMsg string
		Ok     bool
	}{}
	err = c.Database.Run(bson.D{{"dropIndexes", c.Name}, {"index", name}}, &result)
	if err != nil {
		return err
	}
	if !result.Ok {
		return errors.New(result.ErrMsg)
	}
	return nil
}

// nonEventual returns a clone of session and ensures it is not Eventual.
// This guarantees that the server that is used for queries may be reused
// afterwards when a cursor is received.
func (session *Session) nonEventual() *Session {
	cloned := session.Clone()
	if cloned.consistency == Eventual {
		cloned.SetMode(Monotonic, false)
	}
	return cloned
}

// Indexes returns a list of all indexes for the collection.
//
// For example, this snippet would drop all available indexes:
//
//   indexes, err := collection.Indexes()
//   if err != nil {
//       return err
//   }
//   for _, index := range indexes {
//       err = collection.DropIndex(index.Key...)
//       if err != nil {
//           return err
//       }
//   }
//
// See the EnsureIndex method for more details on indexes.
func (c *Collection) Indexes() (indexes []Index, err error) {
	cloned := c.Database.Session.nonEventual()
	defer cloned.Close()

	batchSize := int(cloned.queryConfig.op.limit)

	// Try with a command.
	var result struct {
		Indexes []bson.Raw
		Cursor  cursorData
	}
	var iter *Iter
	err = c.Database.With(cloned).Run(bson.D{{"listIndexes", c.Name}, {"cursor", bson.D{{"batchSize", batchSize}}}}, &result)
	if err == nil {
		firstBatch := result.Indexes
		if firstBatch == nil {
			firstBatch = result.Cursor.FirstBatch
		}
		ns := strings.SplitN(result.Cursor.NS, ".", 2)
		if len(ns) < 2 {
			iter = c.With(cloned).NewIter(nil, firstBatch, result.Cursor.Id, nil)
		} else {
			iter = cloned.DB(ns[0]).C(ns[1]).NewIter(nil, firstBatch, result.Cursor.Id, nil)
		}
	} else if isNoCmd(err) {
		// Command not yet supported. Query the database instead.
		iter = c.Database.C("system.indexes").Find(bson.M{"ns": c.FullName}).Iter()
	} else {
		return nil, err
	}

	var spec indexSpec
	for iter.Next(&spec) {
		indexes = append(indexes, indexFromSpec(spec))
	}
	if err = iter.Close(); err != nil {
		return nil, err
	}
	sort.Sort(indexSlice(indexes))
	return indexes, nil
}

func indexFromSpec(spec indexSpec) Index {
	index := Index{
		Name:             spec.Name,
		Key:              simpleIndexKey(spec.Key),
		Unique:           spec.Unique,
		DropDups:         spec.DropDups,
		Background:       spec.Background,
		Sparse:           spec.Sparse,
		Minf:             spec.Min,
		Maxf:             spec.Max,
		Bits:             spec.Bits,
		BucketSize:       spec.BucketSize,
		DefaultLanguage:  spec.DefaultLanguage,
		LanguageOverride: spec.LanguageOverride,
		ExpireAfter:      time.Duration(spec.ExpireAfter) * time.Second,
		Collation:        spec.Collation,
	}
	if float64(int(spec.Min)) == spec.Min && float64(int(spec.Max)) == spec.Max {
		index.Min = int(spec.Min)
		index.Max = int(spec.Max)
	}
	if spec.TextIndexVersion > 0 {
		index.Key = make([]string, len(spec.Weights))
		index.Weights = make(map[string]int)
		for i, elem := range spec.Weights {
			index.Key[i] = "$text:" + elem.Name
			if w, ok := elem.Value.(int); ok {
				index.Weights[elem.Name] = w
			}
		}
	}
	return index
}

type indexSlice []Index

func (idxs indexSlice) Len() int           { return len(idxs) }
func (idxs indexSlice) Less(i, j int) bool { return idxs[i].Name < idxs[j].Name }
func (idxs indexSlice) Swap(i, j int)      { idxs[i], idxs[j] = idxs[j], idxs[i] }

func simpleIndexKey(realKey bson.D) (key []string) {
	for i := range realKey {
		field := realKey[i].Name
		vi, ok := realKey[i].Value.(int)
		if !ok {
			vf, _ := realKey[i].Value.(float64)
			vi = int(vf)
		}
		if vi == 1 {
			key = append(key, field)
			continue
		}
		if vi == -1 {
			key = append(key, "-"+field)
			continue
		}
		if vs, ok := realKey[i].Value.(string); ok {
			key = append(key, "$"+vs+":"+field)
			continue
		}
		panic("Got unknown index key type for field " + field)
	}
	return
}

// ResetIndexCache() clears the cache of previously ensured indexes.
// Following requests to EnsureIndex will contact the server.
func (s *Session) ResetIndexCache() {
	s.cluster().ResetIndexCache()
}

// New creates a new session with the same parameters as the original
// session, including consistency, batch size, prefetching, safety mode,
// etc. The returned session will use sockets from the pool, so there's
// a chance that writes just performed in another session may not yet
// be visible.
//
// Login information from the original session will not be copied over
// into the new session unless it was provided through the initial URL
// for the Dial function.
//
// See the Copy and Clone methods.
//
func (s *Session) New() *Session {
	s.m.Lock()
	scopy := copySession(s, false)
	s.m.Unlock()
	scopy.Refresh()
	return scopy
}

// Copy works just like New, but preserves the exact authentication
// information from the original session.
func (s *Session) Copy() *Session {
	s.m.Lock()
	scopy := copySession(s, true)
	s.m.Unlock()
	scopy.Refresh()
	return scopy
}
