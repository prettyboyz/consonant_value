// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Solaris system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and wrap
// it in our own nicer implementation, either here or in
// syscall_solaris.go or syscall_unix.go.

package unix

import (
	"sync/atomic"
	"syscall"
	"unsafe"
)

// Implemented in runtime/syscall_solaris.go.
type syscallFunc uintptr

func rawSysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)
func sysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

type SockaddrDatalink struct {
	Family uint16
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [244]int8
	raw    RawSockaddrDatalink
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

// ParseDirent parses up to max directory entries in buf,
// appending the names to names.  It returns the number
// bytes consumed from buf, the number of entries added
// to names, and the new names slice.
func ParseDirent(buf []byte, max int, names []string) (consumed int, count int, newnames []string) {
	origlen := len(buf)
	for max != 0 && len(buf) > 0 {
		dirent := (*Dirent)(unsafe.Pointer(&buf[0]))
		if dirent.Reclen == 0 {
			buf = nil
			break
		}
		buf = buf[dirent.Reclen:]
		if dirent.Ino == 0 { // File absent in directory.
			continue
		}
		bytes := (*[10000]byte)(unsafe.Pointer(&dirent.Name[0]))
		var name = string(bytes[0:clen(bytes[:])])
		if name == "." || name == ".." { // Useless names
			continue
		}
		max--
		count++
		names = append(names, name)
	}
	return origlen - len(buf), count, names
}

//sysnb	pipe(p *[2]_C_int) (n int, err error)

func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	n, err := pipe(&pp)
	if n != 0 {
		return err
	}
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return nil
}

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet4, nil
}

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet6, nil
}

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, _Socklen, error) {
	name := sa.Name
	n := len(name)
	if n >= len(sa.raw.Path) {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = int8(name[i])
	}
	// length is family (uint16), name, NUL.
	sl := _Socklen(2)
	if n > 0 {
		sl += _Socklen(n) + 1
	}
	if sa.raw.Path[0] == '@' {
		sa.raw.Path[0] = 0
		// Don't count trailing NUL for abstract address.
		sl--
	}

	return unsafe.Pointer(&sa.raw), sl, nil
}

//sys	getsockname(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error) = libsocket.getsockname

func Getsockname(fd int) (sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil {
		return
	}
	return anyToSockaddr(&rsa)
}

const ImplementsGetwd = true

//sys	Getcwd(buf []byte) (n int, err error)

func Getwd() (wd string, err error) {
	var buf [PathMax]byte
	// Getcwd will return an error if it failed for any reason.
	_, err = Getcwd(buf[0:])
	if err != nil {
		return "", err
	}
	n := clen(buf[:])
	if n < 1 {
		return "", EINVAL
	}
	return string(buf[:n]), nil
}

/*
 * Wrapped
 */

//sysnb	getgroups(ngid int, gid *_Gid_t) (n int, err error)
//sysnb	setgroups(ngid int, gid *_Gid_t) (err error)

func Getgroups() (gids []int, err error) {
	n, err := getgroups(0, nil)
	// Check for error and sanity check group count.  Newer versions of
	// Solaris allow up to 1024 (NGROUPS_MAX).
	if n < 0 || n > 1024 {
		if err != nil {
			return nil, err
		}
		return nil, EINVAL
	} else if n == 0 {
		return nil, nil
	}

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if n == -1 {
		return nil, err
	}
	gids = make([]int, n)
	for i, v := range a[0:n] {
		gids[i] = int(v)
	}
	return
}

func Setgroups(gids []int) (err error) {
	if len(gids) == 0 {
		return setgroups(0, nil)
	}

	a := make([]_Gid_t, len(gids))
	for i, v := range gids {
		a[i] = _Gid_t(v)
	}
	return setgroups(len(a), &a[0])
}

func ReadDirent(fd int, buf []byte) (n int, err error) {
	// Final argument is (basep *uintptr) and the syscall doesn't take nil.
	// TODO(rsc): Can we use a single global basep for all calls?
	return Getdents(fd, buf, new(uintptr))
}

// Wait status is 7 bits at bottom, either 0 (exited),
// 0x7F (stopped), or a signal number that caused an exit.
// The 0x80 bit is whether there was a core dump.
// An extra number (exit code, signal causing a stop)
// is in the high bits.

type WaitStatus uint32

const (
	mask  = 0x7F
	core  = 0x80
	shift = 8

	exited  = 0
	stopped = 0x7F
)

func (w WaitStatus) Exited() bool { return w&mask == exited }

func (w WaitStatus) ExitStatus() int {
	if w&mask != exited {
		return -1
	}
	return int(w >> shift)
}

func (w WaitStatus) Signaled() bool { return w&mask != stopped && w&mask != 0 }

func (w WaitStatus) Signal() syscall.Signal {
	sig := syscall.Signal(w & mask)
	if sig == stopped || sig == 0 {
		return -1
	}
	return sig
}

func (w WaitStatus) CoreDump() bool { return w.Signaled() && w&core != 0 }

func (w WaitStatus) Stopped() bool { return w&mask == stopped && syscall.Signal(w>>shift) != SIGSTOP }

func (w WaitStatus) Continued() bool { return w&mask == stopped && syscall.Signal(w>>shift) == SIGSTOP }

func (w WaitStatus) StopSignal() syscall.Signal {
	if !w.Stopped() {
		return -1
	}
	return syscall.Signal(w>>shift) & 0xFF
}

func (w WaitStatus) TrapCause() int { return -1 }

//sys	wait4(pid int32, statusp *_C_int, options int, rusage *Rusage) (wpid int32, err error)

func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (int, error) {
	var status _C_int
	rpid, err := wait4(int32(pid), &status, options, rusage)
	wpid := int(rpid)
	if wpid == -1 {
		return wpid, err
	}
	if wstatus != nil {
		*wstatus = WaitStatus(status)
	}
	return wpid, nil
}

//sys	gethostname(buf []byte) (n int, err error)

func Gethostname() (name string, err error) {
	var buf [MaxHostNameLen]byte
	n, err := gethostname(buf[:])
	if n != 0 {
		return "", err
	}
	n = clen(buf[:])
	if n < 1 {
		return "", EFAULT
	}
	return string(buf[:n]), nil
}

//sys	utimes(path string, times *[2]Timeval) (err error)

func Utimes(path string, tv []Timeval) (err error) {
	if tv == nil {
		return utimes(path, nil)
	}
	if len(tv) != 2 {
		return EINVAL
	}
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}

//sys	utimensat(fd int, path string, times *[2]Timespec, flag int) (err error)

func UtimesNano(path string, ts []Timespec) error {
	if ts == nil {
		return utimensat(AT_FDCWD, path, nil, 0)
	}
	if len(ts) != 2 {
		return EINVAL
	}
	return utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
}

func UtimesNanoAt(dirfd int, path string, ts []Timespec, flags int) error {
	if ts == nil {
		return utimensat(dirfd, path, nil, flags)
	}
	if len(ts) != 2 {
		return EINVAL
	}
	return utimensat(dirfd, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), flags)
}

//sys	fcntl(fd int, cmd int, arg int) (val int, err error)

// FcntlFlock performs a fcntl syscall for the F_GETLK, F_SETLK or F_SETLKW command.
func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error {
	_, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procfcntl)), 3, uintptr(fd), uintptr(cmd), uintptr(unsafe.Pointer(lk)), 0, 0, 0)
	if e1 != 0 {
		return e1
	