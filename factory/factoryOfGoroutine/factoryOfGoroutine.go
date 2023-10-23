package factoryOfGoroutine

import (
	"go.lwh.com/linweihao/customerComplaints/utils/goroutine"
)

func MakeEntityOfGoroutine() (entityChannel *goroutine.EntityChannel) {
	entityChannel = goroutine.InitEntityChannel()
	return entityChannel
}
