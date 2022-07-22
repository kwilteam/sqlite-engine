// Copyright 2022 The Sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package embedvfs

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"sync"
	"sync/atomic"
	"unsafe"

	"modernc.org/libc"
	"modernc.org/sqlite/lib"
)

var (
	// Client code must initialize FS before using the embed VFS.
	FS embed.FS

	fToken uintptr

	objectMu sync.Mutex
	objects  = map[uintptr]interface{}{}
)

func token() uintptr { return atomic.AddUintptr(&fToken, 1) }

func addObject(o interface{}) uintptr {
	t := token()
	objectMu.Lock()
	objects[t] = o
	objectMu.Unlock()
	return t
}

func getObject(t uintptr) interface{} {
	objectMu.Lock()
	o := objects[t]
	if o == nil {
		panic("internal error")
	}

	objectMu.Unlock()
	return o
}

func removeObject(t uintptr) {
	objectMu.Lock()
	if _, ok := objects[t]; !ok {
		panic("internal error")
	}

	delete(objects, t)
	objectMu.Unlock()
}

var embedio = sqlite3_io_methods{
	iVersion: 1, // iVersion
}

func init() {
	*(*func(*libc.TLS, uintptr) int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&embedio)) + 8)) = embedClose
	*(*func(*libc.TLS, uintptr, uintptr, int32, sqlite_int64) int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&embedio)) + 16)) = embedRead
	*(*func(*libc.TLS, uintptr, uintptr) int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&embedio)) + 48)) = embedFileSize
	*(*func(*libc.TLS, uintptr, int32) int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&embedio)) + 56)) = embedLock
	*(*func(*libc.TLS, uintptr, int32) int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&embedio)) + 64)) = embedUnlock
	*(*func(*libc.TLS, uintptr, uintptr) int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&embedio)) + 72)) = embedCheckReservedLock
	*(*func(*libc.TLS, uintptr, int32, uintptr) int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&embedio)) + 80)) = embedFileControl
	*(*func(*libc.TLS, uintptr) int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&embedio)) + 88)) = embedSectorSize
	*(*func(*libc.TLS, uintptr) int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&embedio)) + 96)) = embedDeviceCharacteristics
}

func embedFullPathname(tls *libc.TLS, pVfs uintptr, zPath uintptr, nPathOut int32, zPathOut uintptr) int32 {
	out := libc.GoString(zPath)
	for i := 0; i < len(out) && i < int(nPathOut); i++ {
		*(*byte)(unsafe.Pointer(zPathOut)) = out[i]
		zPathOut++
	}
	return sqlite3.SQLITE_OK
}

func embedOpen(tls *libc.TLS, pVfs uintptr, zName uintptr, pFile uintptr, flags int32, pOutFlags uintptr) int32 {
	if zName == 0 {
		return sqlite3.SQLITE_IOERR
	}

	if flags&sqlite3.SQLITE_OPEN_MAIN_JOURNAL != 0 {
		return sqlite3.SQLITE_NOMEM
	}

	p := pFile
	*(*EmbedFile)(unsafe.Pointer(p)) = EmbedFile{}
	f, err := FS.Open(libc.GoString(zName))
	if err != nil {
		panic(err.Error())
		return sqlite3.SQLITE_CANTOPEN
	}

	h := addObject(f)
	(*EmbedFile)(unsafe.Pointer(p)).fsFile = h
	if pOutFlags != 0 {
		*(*int32)(unsafe.Pointer(pOutFlags)) = int32(os.O_RDONLY)
	}
	(*EmbedFile)(unsafe.Pointer(p)).base.pMethods = uintptr(unsafe.Pointer(&embedio))
	return sqlite3.SQLITE_OK
}

func embedRead(tls *libc.TLS, pFile uintptr, zBuf uintptr, iAmt int32, iOfst sqlite_int64) int32 {
	p := pFile
	f := getObject((*EmbedFile)(unsafe.Pointer(p)).fsFile).(fs.File)
	seeker, ok := f.(io.Seeker)
	if !ok {
		return sqlite3.SQLITE_IOERR_READ
	}

	if n, err := seeker.Seek(iOfst, io.SeekStart); err != nil || n != iOfst {
		return sqlite3.SQLITE_IOERR_READ
	}

	b := unsafe.Slice((*byte)(unsafe.Pointer(zBuf)), iAmt)
	n, err := f.Read(b)
	if n == int(iAmt) {
		return sqlite3.SQLITE_OK
	}

	if n < int(iAmt) && err == nil {
		b := b[n:]
		for i := range b {
			b[i] = 0
		}
		return sqlite3.SQLITE_IOERR_SHORT_READ
	}

	return sqlite3.SQLITE_IOERR_READ
}

func embedAccess(tls *libc.TLS, pVfs uintptr, zPath uintptr, flags int32, pResOut uintptr) int32 {
	if flags == sqlite3.SQLITE_ACCESS_READWRITE {
		*(*int32)(unsafe.Pointer(pResOut)) = 0
		return sqlite3.SQLITE_OK
	}

	fn := libc.GoString(zPath)
	if _, err := fs.Stat(FS, fn); err != nil {
		*(*int32)(unsafe.Pointer(pResOut)) = 0
		return sqlite3.SQLITE_OK
	}

	*(*int32)(unsafe.Pointer(pResOut)) = 1
	return sqlite3.SQLITE_OK
}

func embedFileSize(tls *libc.TLS, pFile uintptr, pSize uintptr) int32 {
	p := pFile
	f := getObject((*EmbedFile)(unsafe.Pointer(p)).fsFile).(fs.File)
	fi, err := f.Stat()
	if err != nil {
		return sqlite3.SQLITE_IOERR_FSTAT
	}

	*(*sqlite_int64)(unsafe.Pointer(pSize)) = fi.Size()
	return sqlite3.SQLITE_OK
}

func embedClose(tls *libc.TLS, pFile uintptr) int32 {
	p := pFile
	h := (*EmbedFile)(unsafe.Pointer(p)).fsFile
	f := getObject(h).(fs.File)
	f.Close()
	removeObject(h)
	return sqlite3.SQLITE_OK
}
