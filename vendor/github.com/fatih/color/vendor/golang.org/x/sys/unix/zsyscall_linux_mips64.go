// mksyscall.pl -tags linux,mips64 syscall_linux.go syscall_linux_mips64x.go
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

// +build linux,mips64

package unix

import (
	"syscall"
	"unsafe"
)

var _ syscall.Errno

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Linkat(olddirfd int, oldpath string, newdirfd int, newpath string, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(oldpath)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(newpath)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_LINKAT, uintptr(olddirfd), uintptr(unsafe.Pointer(_p0)), uintptr(newdirfd), uintptr(unsafe.Pointer(_p1)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func openat(dirfd int, path string, flags int, mode uint32) (fd int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall6(SYS_OPENAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(flags), uintptr(mode), 0, 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func ppoll(fds *PollFd, nfds int, timeout *Timespec, sigmask *Sigset_t) (n int, err error) {
	r0, _, e1 := Syscall6(SYS_PPOLL, uintptr(unsafe.Pointer(fds)), uintptr(nfds), uintptr(unsafe.Pointer(timeout)), uintptr(unsafe.Pointer(sigmask)), 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Readlinkat(dirfd int, path string, buf []byte) (n int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	var _p1 unsafe.Pointer
	if len(buf) > 0 {
		_p1 = unsafe.Pointer(&buf[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_READLINKAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(buf)), 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Symlinkat(o