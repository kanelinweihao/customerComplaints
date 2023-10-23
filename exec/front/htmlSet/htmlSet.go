package htmlSet

import (
	// "fmt"
	"go.lwh.com/linweihao/customerComplaints/exec/back/backEnd"
	"go.lwh.com/linweihao/customerComplaints/utils/rfl"
	"go.lwh.com/linweihao/customerComplaints/utils/tmpl"
)

var paramsRoute = rfl.BoxParams[string]{
	"/": {
		"pathTmpl": {
			"pathTmpl": "./public/tmpl/index.tmpl",
		},
		"paramsOutDefault": {
			"ModTitle":      "未知标题",
			"Version":       "v1.0.0",
			"PathDirOutput": "./downloads",
			"UserId":        "1005704",
			"MsgOut":        "Ready",
		},
		"paramsInDefault": {
			"UserId": "1005704",
		},
	},
	"/index": {
		"pathTmpl": {
			"pathTmpl": "./public/tmpl/index.tmpl",
		},
		"paramsOutDefault": {
			"ModTitle":      "未知标题",
			"Version":       "v1.0.0",
			"PathDirOutput": "./downloads",
			"UserId":        "1005704",
			"MsgOut":        "Ready",
		},
		"paramsInDefault": {
			"UserId": "1005704",
		},
	},
}

type EntityBackEnd struct{}

func (self *EntityBackEnd) ExecBackEnd(paramsIn rfl.Params) (paramsFromBackend rfl.Params) {
	paramsFromBackend = backEnd.ExecBackEnd(paramsIn)
	return paramsFromBackend
}

func GetParamsRoute() (paramsRoute rfl.BoxParams[string]) {
	return paramsRoute
}

func SetHtml() {
	entityBackEnd := &EntityBackEnd{}
	tmpl.SetTmplAndServer(paramsRoute, entityBackEnd)
	return
}
