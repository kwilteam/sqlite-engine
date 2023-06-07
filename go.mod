module modernc.org/sqlite

go 1.18

require (
	github.com/google/pprof v0.0.0-20221118152302-e6195bd50e26
	github.com/klauspost/cpuid/v2 v2.2.3
	github.com/mattn/go-sqlite3 v1.14.16
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab
	modernc.org/ccgo/v3 v3.16.13
	modernc.org/libc v1.22.5
	modernc.org/mathutil v1.5.0
	modernc.org/tcl v1.15.2
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/tools v0.0.0-20201124115921-2c860bdd6e78 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	lukechampine.com/uint128 v1.2.0 // indirect
	modernc.org/cc/v3 v3.40.0 // indirect
	modernc.org/httpfs v1.0.6 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/opt v0.1.3 // indirect
	modernc.org/strutil v1.1.3 // indirect
	modernc.org/token v1.0.1 // indirect
	modernc.org/z v1.7.3 // indirect
)

replace modernc.org/sqlite => ./

retract [v1.16.0, v1.17.2] // https://gitlab.com/cznic/sqlite/-/issues/100

retract v1.19.0 // module source tree too large (max size is 524288000 bytes)

retract v1.20.1 // https://gitlab.com/cznic/sqlite/-/issues/123
