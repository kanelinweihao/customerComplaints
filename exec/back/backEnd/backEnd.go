package backEnd

import (
	// "fmt"
	"go.lwh.com/linweihao/customerComplaints/exec/back/cacheSet"
	"go.lwh.com/linweihao/customerComplaints/exec/back/dataGet"
	"go.lwh.com/linweihao/customerComplaints/exec/back/dataPut"
	"go.lwh.com/linweihao/customerComplaints/utils/err"
	"go.lwh.com/linweihao/customerComplaints/utils/rfl"
	"go.lwh.com/linweihao/customerComplaints/utils/time"
)

func ExecBackEnd(paramsIn rfl.Params) (paramsOut rfl.Params) {
	defer err.ThrowError()
	strUserId := paramsIn["UserId"].(string)
	userId := rfl.StrToInt(strUserId)
	boxData, boxParams := dataGet.GetDataFromDB(userId)
	// fmt.Println(boxData)
	// fmt.Println(boxParams)
	dataPut.PutDataToExcel(userId, boxData, boxParams)
	cacheSet.SetCacheOfLogOfPutSuccess(userId)
	paramsOut = paramsIn
	paramsOut["MsgOut"] = "Success"
	// fmt.Println(paramsOut)
	time.ShowTimeAndMsg("OK")
	return paramsOut
}
