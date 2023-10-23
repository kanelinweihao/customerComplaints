package factoryOfSSH

import (
	"go.lwh.com/linweihao/customerComplaints/utils/ssh"
)

func MakeEntityOfSSHForMysql() (entitySSHMysql *ssh.EntitySSH) {
	entitySSHMysql = ssh.InitSSHForMysql()
	return entitySSHMysql
}

func MakeEntityOfSSHForRedis() (entitySSHRedis *ssh.EntitySSH) {
	entitySSHRedis = ssh.InitSSHForRedis()
	return entitySSHRedis
}
