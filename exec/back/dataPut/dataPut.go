package dataPut

import (
	"fmt"
	"go.lwh.com/linweihao/customerComplaints/config/env"
	"go.lwh.com/linweihao/customerComplaints/utils/excel"
	"go.lwh.com/linweihao/customerComplaints/utils/file"
	"go.lwh.com/linweihao/customerComplaints/utils/rfl"
	"go.lwh.com/linweihao/customerComplaints/utils/time"
)

var ext = "xlsx"
var ArrFileNamePrefix = map[int]string{
	1: "给艳艳_客诉_首发留存",
	2: "给艳艳_客诉_寄售买入",
	3: "给艳艳_客诉_寄售卖出",
	4: "给艳艳_客诉_水晶领取",
}

func PutDataToExcel(userId int, boxData rfl.BoxData, boxParams rfl.BoxParams[int]) {
	// time.ShowTimeAndMsg("Data put begin")
	PathDirPutExcel := env.GetPathDirPutExcel()
	file.CreateDir(PathDirPutExcel)
	PathDir := getPathDir(userId)
	file.CreateDir(PathDir)
	for i := 1; i <= 4; i++ {
		go PutOne(
			userId,
			boxData,
			boxParams,
			PathDir,
			i)
	}
	time.ShowTimeAndMsg("Data put success")
	return
}

func PutOne(userId int, boxData rfl.BoxData, boxParams rfl.BoxParams[int], PathDir string, num int) {
	dataNum := boxData[num]
	paramsNum := boxParams[num]
	if len(dataNum) == 0 {
		return
	}
	fileNamePrefixNum := ArrFileNamePrefix[num]
	fileNameNum := getFileName(fileNamePrefixNum, userId)
	pathFileNum := PathDir + "/" + fileNameNum
	pathFileNum = file.GetFilePathAbs(pathFileNum)
	// fmt.Println(pathFileNum)
	excel.SetExcelNew(pathFileNum)
	// fmt.Println(dataNum)
	// fmt.Println(paramsNum)
	excel.WriteExcel(pathFileNum, dataNum, paramsNum)
	// arrAttrFromExcel := excel.ReadExcel(pathFileNum)
	// fmt.Println(arrAttrFromExcel)
	return
}

func getFileName(fileNamePrefix string, userId int) (fileName string) {
	suffix := time.GetSuffix()
	strUserId := rfl.IntToStr(userId)
	remarkUser := "用户" + strUserId
	fileName = fmt.Sprintf(
		"%s_%s_%s.%s",
		fileNamePrefix,
		remarkUser,
		suffix,
		ext)
	// fmt.Println(fileName)
	return fileName
}

func getPathDir(userId int) (PathDir string) {
	PathDirPutExcel := env.GetPathDirPutExcel()
	suffix := time.GetSuffix()
	dirNamePrefix := "数据导出"
	strUserId := rfl.IntToStr(userId)
	remarkUser := "用户" + strUserId
	dirName := fmt.Sprintf(
		"%s_%s_%s",
		dirNamePrefix,
		remarkUser,
		suffix)
	PathDir = PathDirPutExcel + "/" + dirName
	PathDir = file.GetFilePathAbs(PathDir)
	// fmt.Println(PathDir)
	return PathDir
}
