# Copyright 2022 The CC Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY:	clean edit editor tmp x

clean:
	rm -f log-* cpu.test mem.test *.out
	go clean

edit:
	@touch log
	@if [ -f "Session.vim" ]; then gvim -S & else gvim -p Makefile *.c *.go & fi

editor:
	ccgo3 -o vfs.go vfs.c -I../testdata/sqlite-amalgamation-3380500 \
		-lmodernc.org/sqlite/lib -pkgname embedvfs -nocapi \
		-export-externs X -D SQLITE_OS_UNIX \
		-hide=embedFullPathname \
		-hide=embedOpen \
		-hide=embedRead \
		-hide=embedAccess \
		-hide=embedFileSize \
		-hide=embedClose
	go build -v
	gofmt -l -s -w *.go
	go test -o /dev/null -c

