package factoryOfDB

import (
	"go.lwh.com/linweihao/customerComplaints/utils/db"
)

func MakeEntityOfDB() (entityDB *db.EntityDB) {
	entityDB = db.InitDB()
	return entityDB
}
