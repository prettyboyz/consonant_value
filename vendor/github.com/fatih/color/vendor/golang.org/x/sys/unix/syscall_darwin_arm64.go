// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build arm64,darwin

package unix

import (
	"syscall"
	"unsafe"
)

func Getpagesize() int { return 16384 }

func TimespecToNsec(ts Timespec) int64 { return int64(ts.Sec)*1e9 + int64(ts.Nsec) }

func NsecToTimespec(nsec int64) (ts Timespec) {
	ts.Sec = nsec / 1e9
	ts.Nsec = nsec % 1e9
	return
}

func NsecToTimeval(nsec int64) (tv Timeval) {
	nsec += 999 // round up to microsecond
	tv.Usec = int32(nsec % 1e9 / 1e3)
	tv.Sec = int64(nsec / 1e9)
	return
}

//sysnb	gettimeofday(tp *Timeval) (sec int64, usec int32, err error)
func Gettimeofday(tv *Timeval) (err error) {
	// The tv passed to gettimeofday must be non-nil
	// but is otherwise unused.  The answers come back
	// in the two registers.
	sec, usec, err := gettimeofday(tv