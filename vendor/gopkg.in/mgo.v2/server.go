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
	"errors"
	"net"
	"sort"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// ---------------------------------------------------------------------------
// Mongo server encapsulation.

type mongoServer struct {
	sync.RWMutex
	Addr          string
	ResolvedAddr  string
	tcpaddr       *net.TCPAddr
	unusedSockets []*mongoSocket
	liveSockets   []*mongoSocket
	closed        bool
	abended       bool
	sync          chan bool
	dial          dialer
	pingValue     time.Duration
	pingIndex     int
	pingCount     uint32
	pingWindow    [6]time.Duration
	info          *mongoServerInfo
}

type dialer struct {
	old func(addr net.Addr) (net.Conn, error)
	new func(addr *ServerAddr) (net.Conn, error)
}

func (dial dialer) isSet() bool {
	return dial.old != nil || dial.new != nil
}

type mongoServerInfo struct {
	Master         bool
	Mongos         bool
	Tags           bson.D
	MaxWireVersion int
	SetName        string
}

var defaultServerInfo mongoServerInfo

func newServer(addr string, tcpaddr *net.TCPAddr, sync chan bool, dial dialer) *mongoServer {
	server := &mongoServer{
		Addr:         addr,
		ResolvedAddr: tcpaddr.String(),
		tcpaddr:      tcpaddr,
		sync:         sync,
		dial:         dial,
		info:         &defaultServerInfo,
		pingValue:    time.Hour, // Push it back before an actual ping.
	}
	go server.pinger(true)
	return server
}

var errPoolLimit = errors.New("per-server connection limit reached")
var errServerClosed = errors.New("server was closed")

// AcquireSocket returns a socket for communicating with the server.
// This will attempt to reuse an old connection, if one is available. Otherwise,
// it will establish a new one. The returned socket is owned by the call site,
// and will return to the cache when the socket has its Release method called
// the same number of times as AcquireSocket + Acquire were called for it.
// If the poolLimit argument is greater than zero and the number of sockets in
// use in this server is greater than the provided limit, errPoolLimit is
// returned.
func (server *mongoServer) AcquireSocket(poolLimit int, timeout time.Duration) (socket *mongoSocket, abended bool, err error) {
	for {
		server.Lock()
		abended = server.abended
		if server.closed {
			server.Unlock()
			return nil, abended, errServerClosed
		}
		n := len(server.unusedSockets)
		if poolLimit > 0 && len(server.liveSockets)-n >= poolLimit {
			server.Unlock()
			return nil, false, errPoolLimit
		}
		if n > 0 {
			socket = server.unusedSockets[n-1]
			server.unusedSockets[n-1] = nil // Help GC.
			server.unusedSockets = server.unusedSockets[:n-1]
			info := server.info
			server.Unlock()
			err = socket.InitialAcquire(info, timeout)
			if err != nil {
				continue
			}
		} else {
			server.Unlock()
			socket, err = server.Connect(timeout)
			if err == nil {
				server.Lock()
				// We've waited for the Connect, see if we got
				// closed in the meantime
				if server.closed {
					server.Unlock()
					socket.Release()
					socket.Close()
					return nil, abended, errServerClosed
				}
				server.liveSockets = append(server.liveSockets, socket)
				server.Unlock()
			}
		}
		return
	}
	panic("unreachable")
}

// Connect establishes a new connection to the server. This should
// generally be done through server.AcquireSocket().
func (server *mongoServer) Connect(timeout time.Duration) (*mongoSocket, error) {
	server.RLock()
	master := server.info.Master
	dial := server.dial
	server.RUnlock()

	logf("Establishing new connection to %s (timeout=%s)...", server.Addr, timeout)
	var conn net.Conn
	var err error
	switch {
	case !dial.isSet():
		// Cannot do this because it lacks timeout support. :-(
		//conn, err = net.DialTCP("tcp", nil, server.tcpaddr)
		conn, err = net.DialTimeout("tcp", server.ResolvedAddr, timeout)
		if tcpconn, ok := conn.(*net.TCPConn); ok {
			tcpconn.SetKeepAlive(true)
		} else if err == nil {
			panic("internal error: obtained TCP connection is not a *net.TCPConn!?")
		}
	case dial.old != nil:
		conn, err = dial.old(server.tcpaddr)
	case dial.new != nil:
		conn, err = dial.new(&ServerAddr{server.Addr, server.tcpaddr})
	default:
		panic("dialer is set, but both dial.old and dial.new are nil")
	}
	if err != nil {
		logf("Connection to %s failed: %v", server.Addr, err.Error())
		return nil, err
	}
	logf("Connection to %s established.", server.Addr)

	stats.conn(+1, master)
	return newSocket(server, conn, timeout), nil
}

// Close forces closing all sockets that are alive, whethe