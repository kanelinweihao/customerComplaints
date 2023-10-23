package file

import (
	"fmt"
	"go.lwh.com/linweihao/customerComplaints/utils/err"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	// "go.lwh.com/linweihao/customerComplaints/utils/rfl"
)

/*
Ext
*/

func GetExt(fileName string) (ext string) {
	/*kDot := strings.LastIndex(fileName, ".")
	  kBegin := kDot + 1
	  ext = fileName[kBegin:]*/
	extWithDot := filepath.Ext(fileName)
	ext = extWithDot[1:]
	// fmt.Println(ext)
	return ext
}

/*
Dir
*/

func HasDir(path string) (isExisted bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(pathDir string) {
	isExisted, errDirExisted := HasDir(pathDir)
	err.ErrCheck(errDirExisted)
	if isExisted {
		/*msgError := fmt.Sprintf(
		    "The dir of |%s| is existed",
		    pathDir)
		err.ErrPanic(msgError)*/
		return
	}
	errDirMk := os.MkdirAll(pathDir, os.ModePerm)
	err.ErrCheck(errDirMk)
	return
}

/*
Path
*/

func IsPathAbs(filePath string) (isAbs bool) {
	isAbs = filepath.IsAbs(filePath)
	return isAbs
}

// 相对路径转绝对路径(实现对windows\分隔符的兼容)
func GetFilePathAbs(filePath string) (filePathAbs string) {
	isAbs := IsPathAbs(filePath)
	if isAbs {
		filePathAbs = filePath
		return filePathAbs
	}
	filePathAbs, errAbs := filepath.Abs(filePath)
	err.ErrCheck(errAbs)
	return filePathAbs
}

// 绝对路径转相对路径(实现对embed静态文件打包的兼容)
func GetFilePathRel(filePath string) (filePathRel string) {
	isAbs := IsPathAbs(filePath)
	if !isAbs {
		filePath = GetFilePathAbs(filePath)
	}
	filePathBase, errBase := filepath.Abs("./")
	err.ErrCheck(errBase)
	filePathRel, errRel := filepath.Rel(filePathBase, filePath)
	err.ErrCheck(errRel)
	return filePathRel
}

// 文件路径强制用/作分隔符(实现对embed静态文件打包的兼容)
func GetFilePathEmbed(filePath string) (filePathEmbed string) {
	filePathRel := GetFilePathRel(filePath)
	filePathEmbed = filepath.ToSlash(filePathRel)
	return filePathEmbed
}

/*
FileRead
*/

func ReadFileAsArrayByte(pathFile string) (arrayByte []byte) {
	fs, errFileOpen := os.OpenFile(
		pathFile,
		os.O_RDONLY,
		0666)
	err.ErrCheck(errFileOpen)
	// rfl.ShowType(fs)
	// *os.File
	defer fs.Close()
	arrayByte, errFileRead := ioutil.ReadAll(fs)
	err.ErrCheck(errFileRead)
	// fmt.Println(arrayByte)
	// rfl.ShowType(arrayByte)
	return arrayByte
}

func ReadFileAsString(pathFile string) (strFromFile string) {
	arrayByte := ReadFileAsArrayByte(pathFile)
	strFromFile = string(arrayByte)
	fmt.Println(strFromFile)
	return strFromFile
}

func readFileAsArrayString(pathFile string) (arrayStr []string) {
	strFromFile := ReadFileAsString(pathFile)
	arrayStr = strings.Split(strFromFile, "\n")
	// fmt.Println(arrayStr)
	return arrayStr
}
