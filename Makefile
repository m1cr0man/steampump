.PHONY: default
default: generate build

.PHONY: builddeps
builddeps:
	go get github.com/akavel/rsrc
	go get github.com/markbates/pkger/cmd/pkger

.PHONY: generate
generate:
	cat steampump.ico | ./embed.sh main Icon --compress > cmd/steampump/icon.go
	rsrc -manifest steampump.exe.manifest -ico steampump.ico -o cmd/steampump/rsrc.syso
	pkger -o pkg/server

.PHONY: build
build:
	go build -ldflags -H=windowsgui -o steampump.exe cmd/steampump/main.go
