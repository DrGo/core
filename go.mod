module github.com/drgo/core

go 1.21

require (
	github.com/dannyvankooten/extemplate v0.0.0-20221206123735-ea3f2b2b17ac
	github.com/fsnotify/fsnotify v1.6.0
	github.com/matryer/is v1.4.1
	golang.org/x/crypto v0.14.0
)

require (
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/drgo/mdson v0.0.0 => ../mdson/
