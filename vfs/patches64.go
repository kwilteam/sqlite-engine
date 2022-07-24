// Copyright 2022 The Sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !arm && !386
// +build !arm,!386

package vfs

import (
	"modernc.org/libc"
	sqlite3 "modernc.org/sqlite/lib"
)

func vfsFullPathname(tls *libc.TLS, pVfs uintptr, zPath uintptr, nPathOut int32, zPathOut uintptr) int32 {
	libc.Xstrncpy(tls, zPathOut, zPath, uint64(nPathOut))
	return sqlite3.SQLITE_OK
}
