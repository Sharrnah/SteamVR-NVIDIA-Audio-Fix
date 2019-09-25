set GOARCH=amd64
set GOOS=windows
go get github.com/akavel/rsrc
go get ./...
rsrc -manifest steamvraudiofix.exe.manifest -o steamvraudiofix.syso
go build -v -o steamvraudiofix.exe
