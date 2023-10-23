package main

/*commandSetEnv*/
// go env -w GO111MODULE=on
// go env -w GOPROXY=https://goproxy.cn,direct
// go env -w CGO_ENABLED=0
// go env -w GOOS=windows
// go env -w GOARCH=amd64
/*commandInitMod*/
// mkdir customerComplaints
// cd ./customerComplaints
// go mod init go.lwh.com/linweihao/customerComplaints
/*commandImportMod*/
// go get -u github.com/go-sql-driver/mysql
// go get -u github.com/jmoiron/sqlx
// go get -u github.com/xuri/excelize/v2
// go get -u golang.org/x/crypto/ssh
// go get -u github.com/go-redis/redis
// go get -u github.com/akavel/rsrc
/*commandBuild*/
// rsrc -manifest ./build/main.manifest -ico ./build/icon_go.ico -o rsrc.syso
// go fmt ./...
// go mod tidy
// GOOS=windows go build -o 数艺客诉导出.exe
/*commandRun*/
// go run main.go

import (
	"embed"
	"fmt"
	"go.lwh.com/linweihao/customerComplaints/config/env"
	"go.lwh.com/linweihao/customerComplaints/config/version"
	"go.lwh.com/linweihao/customerComplaints/exec/front/frontEnd"
	"go.lwh.com/linweihao/customerComplaints/utils/err"
	"go.lwh.com/linweihao/customerComplaints/utils/time"
)

//go:embed public
var FilesPublic embed.FS

func init() {
	fmt.Println("\n")
	msgVersion := version.GetMsgVersion()
	time.ShowTimeAndMsg(msgVersion)
	fmt.Println("\n")
	defer err.ThrowError()
	env.EnbaleCPU()
	env.SetFilesPublic(FilesPublic)
	return
}

func main() {
	doIt()
	return
}

func doIt() {
	frontEnd.ExecFrontEnd()
	return
}
