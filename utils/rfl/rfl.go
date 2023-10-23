package rfl

import (
	"fmt"
	"go.lwh.com/linweihao/customerComplaints/utils/err"
	"reflect"
	"strconv"
)

type Attr map[string]interface{}
type ArrAttr map[int]Attr
type AttrForExcel map[string]string
type ArrAttrForExcel map[int]AttrForExcel
type BoxData map[int]ArrAttrForExcel
type Params map[string]interface{}
type ArrParams map[string]Params
type BoxParams[T int | string] map[T]ArrParams

func ShowType(value interface{}) {
	typeOfValue := reflect.TypeOf(value)
	fmt.Println(typeOfValue)
	return
}

func GetTypeName(value interface{}) (typeName string) {
	typeOfValue := reflect.TypeOf(value)
	typeName = typeOfValue.Name()
	return typeName
}

func GetValue(value interface{}) (valueOfValue reflect.Value) {
	valueOfValue = reflect.ValueOf(value)
	return valueOfValue
}

/*int->string*/

func IntToStr(num int) (str string) {
	str = strconv.Itoa(num)
	return str
}

/*string->int*/

func StrToInt(str string) (num int) {
	num, errAtoi := strconv.Atoi(str)
	err.ErrCheck(errAtoi)
	return num
}

/*Entity->Map*/

func ToAttrFromEntity(entity interface{}) (attr map[string]interface{}) {
	attr = Attr{}
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)
	// fmt.Println(t)
	// fmt.Println(v)
	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Name
		value := v.Field(i).Interface()
		attr[key] = value
	}
    // fmt.Println(attr)
	// typeOfAttr := reflect.TypeOf(attr)
	// fmt.Println(typeOfAttr)
	return attr
}

func ToArrAttrFromArrEntity[T interface{}](arrEntity []T) (arrAttr ArrAttr) {
	arrAttr = ArrAttr{}
	for k, entity := range arrEntity {
		attr := ToAttrFromEntity(entity)
		// fmt.Println(attr)
		arrAttr[k] = attr
	}
	// fmt.Println(arrAttr)
	return arrAttr
}

/*Map->Entity*/

func ToEntityFromAttr(attr map[string]interface{}, entity interface{}) {
	structValue := reflect.ValueOf(entity).Elem()
	for k, v := range attr {
		structFieldValue := structValue.FieldByName(k)
		if !structFieldValue.IsValid() {
			continue
		}
		if !structFieldValue.CanSet() {
			msgError := fmt.Sprintf(
				"The field |%s| can not set",
				k)
			err.ErrPanic(msgError)
		}
		val := reflect.ValueOf(v)
		if structFieldValue.Type() != val.Type() {
			msgError := fmt.Sprintf(
				"The value |%s| of field |%s| is invalid, th type is |%s|, it should be |%s|",
				v,
				k,
				val.Type(),
				structFieldValue.Type())
			err.ErrPanic(msgError)
		}
		structFieldValue.Set(val)
	}
	return
}

/*ArrAttr->ArrAttrForExcel*/

func ToAttrForExcelFromAttr(attr Attr) (attrForExcel AttrForExcel) {
	attrForExcel = AttrForExcel{}
	for key, value := range attr {
		// fmt.Println(value)
		typeName := GetTypeName(value)
		// fmt.Println(typeName)
		// isString := (typeName == "string")
		// fmt.Println(isString)
		var valueForExcel string
		switch typeName {
		case "int":
			valueForExcel = fmt.Sprintf(
				"%d",
				value)
		case "bool":
			valueForExcel = fmt.Sprintf(
				"%t",
				value)
		case "string":
			valueForExcel = fmt.Sprintf(
				"%s",
				value)
		default:
			valueForExcel = fmt.Sprintf(
				"%s",
				value)
		}
		attrForExcel[key] = valueForExcel
	}
	// fmt.Println(attr)
	// fmt.Println(attrForExcel)
	// typeOfAttrForExcel := reflect.TypeOf(attrForExcel)
	// fmt.Println(typeOfAttrForExcel)
	return attrForExcel
}

func ToArrAttrForExcelFromArrAttr(arrAttr ArrAttr) (arrAttrForExcel ArrAttrForExcel) {
	arrAttrForExcel = ArrAttrForExcel{}
	for k, attr := range arrAttr {
		attrForExcel := ToAttrForExcelFromAttr(attr)
		// fmt.Println(attrForExcel)
		arrAttrForExcel[k] = attrForExcel
	}
	// fmt.Println(arrAttrForExcel)
	return arrAttrForExcel
}

/*array_merge*/

func ParamsMerge(params1 Params, params2 Params) (params3 Params) {
	params3 = Params{}
	for k1, v1 := range params1 {
		params3[k1] = v1
	}
	for k2, v2 := range params2 {
		params3[k2] = v2
	}
	return params3
}
