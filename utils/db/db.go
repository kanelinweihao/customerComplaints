package db

import (
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.lwh.com/linweihao/customerComplaints/config/env"
	"go.lwh.com/linweihao/customerComplaints/utils/err"
	"go.lwh.com/linweihao/customerComplaints/utils/rfl"
	// "go.lwh.com/linweihao/customerComplaints/utils/time"
	"go.lwh.com/linweihao/customerComplaints/factory/factoryOfSSH"
	"go.lwh.com/linweihao/customerComplaints/utils/goroutine"
	"go.lwh.com/linweihao/customerComplaints/utils/ssh"
)

var driver string = "mysql"

type EntityDB struct {
	DBSqlx    *sqlx.DB
	EntitySSH *ssh.EntitySSH
}
type TypeEntityData interface{}
type EntityConfigMysql struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	Charset  string
}

/*
Init
*/

func InitDB() (entityDb *EntityDB) {
	entitySSH := getEntitySSH()
	dbSqlx := getDBSqlx()
	entityDb = &EntityDB{
		DBSqlx:    dbSqlx,
		EntitySSH: entitySSH,
	}
	return entityDb
}

func getEntitySSH() (entitySSH *ssh.EntitySSH) {
	entitySSH = nil
	isNeedSSH := isNeedSSH()
	if isNeedSSH {
		// entitySSH = ssh.InitSSHForMysql()
		entitySSH = factoryOfSSH.MakeEntityOfSSHForMysql()
	}
	return entitySSH
}

func isNeedSSH() (isNeedSSH bool) {
	isNeedSSH = env.IsNeedSSH()
	return isNeedSSH
}

func getDBSqlx() (dbSqlx *sqlx.DB) {
	m := getEntityConfigMysql()
	dsn := getDSN(m)
	dbSqlx, errDb := sqlx.Open(driver, dsn)
	err.ErrCheck(errDb)
	// rfl.ShowType(dbSqlx)
	// ping := dbSqlx.Ping
	// fmt.Println(ping)
	return dbSqlx
}

func getEntityConfigMysql() (m *EntityConfigMysql) {
	paramsMysql := env.GetParamsMysql()
	m = &EntityConfigMysql{}
	rfl.ToEntityFromAttr(paramsMysql, m)
	return m
}

func getDSN(m *EntityConfigMysql) (dsn string) {
	mysqlHost := m.Host
	mysqlPort := m.Port
	mysqlUser := m.User
	mysqlPassword := m.Password
	mysqlDbname := m.Dbname
	mysqlCharset := m.Charset
	isNeedSSH := isNeedSSH()
	network := ssh.NetworkTCP
	if isNeedSSH {
		network = ssh.NetworkTCPSSH
	}
	networkFull := fmt.Sprintf(
		"%s(%s:%s)",
		network,
		mysqlHost,
		mysqlPort)
	// fmt.Println(tcp)
	dsn = fmt.Sprintf(
		"%s:%s@%s/%s?charset=%s",
		mysqlUser,
		mysqlPassword,
		networkFull,
		mysqlDbname,
		mysqlCharset)
	// fmt.Println(dsn)
	return dsn
}

/*
Exec
*/

func (self *EntityDB) CloseDB() {
	dbSqlx := self.DBSqlx
	dbSqlx.Close()
	// time.ShowTimeAndMsg("DB close success")
	isNeedSSH := env.IsNeedSSH()
	if isNeedSSH {
		entitySSH := self.EntitySSH
		entitySSH.CloseSSH()
		// time.ShowTimeAndMsg("SSH close success")
	}
	return
}

func GetArrAttrForExcel[T interface{}](self *EntityDB, arrEntity []T, query string, userIdCc int) (arrAttrForExcel rfl.ArrAttrForExcel) {
	dbSqlx := self.DBSqlx
	// fmt.Println(query)
	errSqlxSelect := dbSqlx.Select(
		&arrEntity,
		query,
		userIdCc)
	err.ErrCheck(errSqlxSelect)
	// fmt.Println(arrEntity)
	arrAttr := rfl.ToArrAttrFromArrEntity(arrEntity)
	// fmt.Println(arrAttr)
	arrAttrForExcel = rfl.ToArrAttrForExcelFromArrAttr(arrAttr)
	// fmt.Println(arrAttrForExcel)
	return arrAttrForExcel
}

func GetArrAttrForExcelUseGoroutine[T interface{}](self *EntityDB, entityChannel *goroutine.EntityChannel, arrEntity []T, query string, userIdCc int) {
	arrAttrForExcel := GetArrAttrForExcel(
		self,
		arrEntity,
		query,
		userIdCc)
	entityChannel.WriteOnce(arrAttrForExcel)
	// time.ShowTimeAndMsg("channel write success")
	return
}
