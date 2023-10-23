package ssh

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go.lwh.com/linweihao/customerComplaints/config/env"
	"go.lwh.com/linweihao/customerComplaints/utils/err"
	"go.lwh.com/linweihao/customerComplaints/utils/rfl"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
	// "go.lwh.com/linweihao/customerComplaints/utils/file"
	// "go.lwh.com/linweihao/customerComplaints/utils/ip"
	// "go.lwh.com/linweihao/customerComplaints/utils/time"
)

var NetworkTCP string = "tcp"
var NetworkTCPSSH string = "tcp+ssh"
var TimeoutSSH time.Duration = time.Second * time.Duration(5)

type EntitySSH struct {
	clientDial *ssh.Client
	network    string
	address    string
}
type EntityClientSSH struct {
	Host           string
	Port           string
	User           string
	TypeAuth       string
	Password       string
	PathPrivateKey string
}

/*
Init
*/

func InitSSHForMysql() (entitySSH *EntitySSH) {
	entitySSH = initSSH()
	entitySSH.RegisterDialToMysql()
	return entitySSH
}

func InitSSHForRedis() (entitySSH *EntitySSH) {
	entitySSH = initSSH()
	return entitySSH
}

func initSSH() (entitySSH *EntitySSH) {
	// ip.ShowIP()
	s := getEntityClientSSH()
	// fmt.Println(s.TypeAuth)
	clientDial, network, address := dialSSH(s)
	// fmt.Println(clientDial)
	entitySSH = &EntitySSH{
		clientDial: clientDial,
		network:    network,
		address:    address,
	}
	return entitySSH
}

func getEntityClientSSH() (s *EntityClientSSH) {
	paramsSSH := env.GetParamsSSH()
	s = &EntityClientSSH{}
	rfl.ToEntityFromAttr(paramsSSH, s)
	s.PathPrivateKey = env.GetPathPrivateKey()
	return s
}

func dialSSH(s *EntityClientSSH) (clientDial *ssh.Client, network string, address string) {
	sshHost := s.Host
	sshPort := s.Port
	sshUser := s.User
	sshTypeAuth := s.TypeAuth
	sshPassword := s.Password
	sshPathPrivateKey := s.PathPrivateKey
	network = "tcp"
	address = getAddressToDial(
		sshHost,
		sshPort)
	configSSHClient := getConfigSSHClient(
		sshUser,
		sshTypeAuth,
		sshPassword,
		sshPathPrivateKey)
	clientDial, errDial := ssh.Dial(
		network,
		address,
		configSSHClient)
	err.ErrCheck(errDial)
	// time.ShowTimeAndMsg("Dial set client success")
	// ip.ShowIP()
	return clientDial, network, address
}

func getAddressToDial(sshHost string, sshPort string) (address string) {
	address = fmt.Sprintf(
		"%s:%s",
		sshHost,
		sshPort)
	return address
}

func getConfigSSHClient(sshUser string, sshTypeAuth string, sshPassword string, sshPathPrivateKey string) (configSSHClient *ssh.ClientConfig) {
	hostKeyCallBack := ssh.InsecureIgnoreHostKey()
	timeoutSSH := TimeoutSSH
	entitySSHConfig := ssh.ClientConfig{
		User:            sshUser,
		HostKeyCallback: hostKeyCallBack,
		Timeout:         timeoutSSH,
	}
	sshAuth := getConfigSSHAuth(
		sshTypeAuth,
		sshPassword,
		sshPathPrivateKey)
	entitySSHConfig.Auth = sshAuth
	configSSHClient = &entitySSHConfig
	return configSSHClient
}

func getConfigSSHAuth(sshTypeAuth string, sshPassword string, sshPathPrivateKey string) (sshAuth []ssh.AuthMethod) {
	switch sshTypeAuth {
	case "Password":
		sshAuth = []ssh.AuthMethod{
			ssh.Password(sshPassword),
		}
	case "Key":
		// arrayBytePrivateKey := file.ReadFileAsArrayByte(sshPathPrivateKey)
		arrayBytePrivateKey := env.GetArrayBytePrivateKey(sshPathPrivateKey)
		signer, errParse := ssh.ParsePrivateKey(arrayBytePrivateKey)
		err.ErrCheck(errParse)
		sshAuth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	default:
		msgError := fmt.Sprintf(
			"The typeAuth |%s| is invalid of |%s|",
			sshTypeAuth,
			"dialSSH")
		err.ErrPanic(msgError)
	}
	return sshAuth
}

/*
Exec
*/

func (self *EntitySSH) GetClientDial() (clientDial *ssh.Client) {
	return self.clientDial
}

func (self *EntitySSH) GetNetwork() (network string) {
	return self.network
}

func (self *EntitySSH) GetAddress() (address string) {
	return self.address
}

func (self *EntitySSH) SetAddress(address string) {
	self.address = address
	return
}

func (self *EntitySSH) CloseSSH() {
	clientDial := self.GetClientDial()
	clientDial.Close()
	return
}

func (self *EntitySSH) RegisterDialToMysql() {
	funcDialMysql := self.DialForMysql
	networkDialMysql := NetworkTCPSSH
	mysql.RegisterDialContext(
		networkDialMysql,
		funcDialMysql)
	// time.ShowTimeAndMsg("Dial register mysql success")
	// ip.ShowIP()
	return
}

func (self *EntitySSH) DialForMysql(context context.Context, address string) (coon net.Conn, errDial error) {
	// fmt.Println(context)
	// fmt.Println(address)
	network := self.network
	coon, errDial = self.clientDial.Dial(network, address)
	err.ErrCheck(errDial)
	return coon, errDial
}

func (self *EntitySSH) DialForRedis() (coon net.Conn, errDial error) {
	network := self.network
	address := self.address
	// fmt.Println(network)
	// fmt.Println(address)
	coon, errDial = self.clientDial.Dial(network, address)
	err.ErrCheck(errDial)
	return coon, errDial
}
