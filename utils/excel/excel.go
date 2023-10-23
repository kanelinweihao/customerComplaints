package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"go.lwh.com/linweihao/customerComplaints/utils/err"
	"go.lwh.com/linweihao/customerComplaints/utils/file"
	"go.lwh.com/linweihao/customerComplaints/utils/rfl"
)

var extExcel string = "xlsx"
var arrCellRow = map[int]string{
	1:  "A",
	2:  "B",
	3:  "C",
	4:  "D",
	5:  "E",
	6:  "F",
	7:  "G",
	8:  "H",
	9:  "I",
	10: "J",
	11: "K",
	12: "L",
	13: "M",
	14: "N",
	15: "O",
	16: "P",
	17: "Q",
	18: "R",
	19: "S",
	20: "T",
	21: "U",
	22: "V",
	23: "W",
	24: "X",
	25: "Y",
	26: "Z",
}

/*
Public
*/

func SetExcelNew(pathFile string) {
	checkExt(pathFile)
	entityExcel := setFileNew()
	saveFile(entityExcel, pathFile)
	return
}

func WriteExcel(pathFile string, arrAttrForExcel rfl.ArrAttrForExcel, paramsForExcel rfl.ArrParams) {
	// fmt.Println(arrAttrForExcel)
	entityExcel := openFile(pathFile)
	defer closeFile(entityExcel)
	sheetName := "Sheet1"
	batchSetCell(
		entityExcel,
		sheetName,
		arrAttrForExcel,
		paramsForExcel)
	saveFile(entityExcel, pathFile)
	return
}

func ReadExcel(pathFile string) (arrAttrFromExcel rfl.ArrAttrForExcel) {
	entityExcel := openFile(pathFile)
	defer closeFile(entityExcel)
	sheetName := "Sheet1"
	arrAttrFromExcel = batchGetCell(entityExcel, sheetName)
	return arrAttrFromExcel
}

/*
Private
*/

func setFileNew() (entityExcel *excelize.File) {
	entityExcel = excelize.NewFile()
	return entityExcel
}

func saveFile(entityExcel *excelize.File, pathFile string) {
	errExcelSave := entityExcel.SaveAs(pathFile)
	err.ErrCheck(errExcelSave)
	return
}

func openFile(pathFile string) (entityExcel *excelize.File) {
	entityExcel, errExcelOpen := excelize.OpenFile(pathFile)
	err.ErrCheck(errExcelOpen)
	return entityExcel
}

func closeFile(entityExcel *excelize.File) {
	errExcelClose := entityExcel.Close()
	err.ErrCheck(errExcelClose)
	return
}

func setSheetNew(entityExcel *excelize.File, num int) {
	numString := rfl.IntToStr(num)
	sheetNameNew := "Sheet" + numString
	numSheet, errSheet := entityExcel.NewSheet(sheetNameNew)
	err.ErrCheck(errSheet)
	entityExcel.SetActiveSheet(numSheet)
	return
}

func setCellValue(entityExcel *excelize.File, sheetName string, cellPosition string, value string) {
	entityExcel.SetCellValue(
		sheetName,
		cellPosition,
		value)
	return
}

func getCellValue(entityExcel *excelize.File, sheetName string, cellPosition string) (cellValue string) {
	cellValue, errExcelGetCell := entityExcel.GetCellValue(sheetName, cellPosition)
	err.ErrCheck(errExcelGetCell)
	return cellValue
}

func checkExt(pathFile string) {
	ext := file.GetExt(pathFile)
	if ext == extExcel {
		return
	}
	msgError := fmt.Sprintf(
		"文件[%s]格式错误,正确后缀应是[%s],当前实际是[%s]",
		pathFile,
		extExcel,
		ext)
	err.ErrPanic(msgError)
}

func batchSetCell(entityExcel *excelize.File, sheetName string, arrAttrForExcel rfl.ArrAttrForExcel, paramsForExcel rfl.ArrParams) {
	var cellRow string
	var cellColumn int
	var cellPosition string
	var cellValue string
	for k, attrForExcel := range arrAttrForExcel {
		// fmt.Println(k)
		// fmt.Println(attr)
		cellColumn = k + 2
		// fmt.Println(cellColumn)
		for field, value := range attrForExcel {
			// fmt.Println(field)
			// fmt.Println(value)
			cellRow = getCellRow(paramsForExcel, field)
			if cellRow == "未知表格列序号" {
				continue
			}
			cellPosition = fmt.Sprintf(
				"%s%d",
				cellRow,
				cellColumn)
			// fmt.Println(cellPosition)
			cellValue = value
			setCellValue(
				entityExcel,
				sheetName,
				cellPosition,
				cellValue)
			if k == 0 {
				batchSetCellOfTitle(
					entityExcel,
					sheetName,
					cellRow,
					paramsForExcel,
					field)
			}
		}
	}
}

func getCellRow(paramsForExcel rfl.ArrParams, field string) (cellRow string) {
	// fmt.Println(field)
	// fmt.Println(paramsForExcel)
	sort := getSortByField(paramsForExcel, field)
	// fmt.Println(sort)
	cellRow = getCellRowBySort(sort)
	return cellRow
}

func getSortByField(paramsForExcel rfl.ArrParams, field string) (sort int) {
	sort = 0
	sortFromMap := paramsForExcel[field]["sort"]
	sortInt, isInt := sortFromMap.(int)
	if isInt {
		sort = sortInt
	}
	return sort
}

func getCellRowBySort(sort int) (cellRow string) {
	if sort == 0 {
		return "未知表格列序号"
	}
	if sort > 26 {
		return "未知表格列序号"
	}
	cellRow = arrCellRow[sort]
	// fmt.Println(cellRow)
	return cellRow
}

func getValueTitle(paramsForExcel rfl.ArrParams, field string) (valueTitle string) {
	valueTitle = "未知字段"
	titleFromMap := paramsForExcel[field]["title"]
	strTitle, isString := titleFromMap.(string)
	if isString {
		valueTitle = strTitle
	}
	return valueTitle
}

func batchSetCellOfTitle(entityExcel *excelize.File, sheetName string, cellRow string, paramsForExcel rfl.ArrParams, field string) {
	cellColumnTitle := 1
	cellPositionTitle := fmt.Sprintf(
		"%s%d",
		cellRow,
		cellColumnTitle)
	valueTitle := getValueTitle(paramsForExcel, field)
	entityExcel.SetCellValue(
		sheetName,
		cellPositionTitle,
		valueTitle)
	return
}

func batchGetCell(entityExcel *excelize.File, sheetName string) (arrAttrFromExcel rfl.ArrAttrForExcel) {
	arrayRows, errExcelGetRows := entityExcel.GetRows(sheetName)
	err.ErrCheck(errExcelGetRows)
	// rfl.ShowType(arrayRows)
	// fmt.Println(arrayRows)
	arrAttrFromExcel = rfl.ArrAttrForExcel{}
	arrayRowTitle := arrayRows[0]
	// fmt.Println(arrayRowTitle)
	for k, arrayRow := range arrayRows {
		// fmt.Println(k)
		// fmt.Println(arrayRow)
		// rfl.ShowType(arrayRow)
		if k == 0 {
			continue
		}
		attrFromExcel := rfl.AttrForExcel{}
		for sort, value := range arrayRow {
			// fmt.Println(sort)
			// fmt.Println(value)
			// rfl.ShowType(value)
			title := getTitleBySort(arrayRowTitle, sort)
			attrFromExcel[title] = value
		}
		arrAttrFromExcel[k] = attrFromExcel
	}
	// fmt.Println(arrAttrFromExcel)
	return arrAttrFromExcel
}

func getTitleBySort(arrayRowTitle []string, sort int) (title string) {
	// fmt.Println(arrayRowTitle)
	var arrTitle = map[int]string{}
	for sort, title := range arrayRowTitle {
		// fmt.Println(sort)
		// fmt.Println(title)
		arrTitle[sort] = title
	}
	// fmt.Println(mapTitle)
	title = arrTitle[sort]
	// fmt.Println(title)
	return title
}
