package tmpl

import (
	"fmt"
	"go.lwh.com/linweihao/customerComplaints/config/env"
	"go.lwh.com/linweihao/customerComplaints/utils/err"
	"go.lwh.com/linweihao/customerComplaints/utils/rfl"
	"go.lwh.com/linweihao/customerComplaints/utils/time"
	"html/template"
	"net/http"
	"net/url"
	"os/exec"
)

var paramsRoute = rfl.BoxParams[string]{}
var commands = map[string]string{
	"windows": "cmd /c start",
	"darwin":  "open",
	"linux":   "xdg-open",
}
var entityBackEnd backEndExecer

type backEndExecer interface {
	ExecBackEnd(paramsIn rfl.Params) (paramsFromBackend rfl.Params)
}

/*
Html
*/

func SetTmplAndServer(paramsRouteToSet rfl.BoxParams[string], entityBackEndToSet backEndExecer) {
	if len(paramsRouteToSet) == 0 {
		return
	}
	entityBackEnd = entityBackEndToSet
	paramsRoute = paramsRouteToSet
	setTmplToHandle()
	time.Sleep(100, "ms")
	// OpenExplorer()
	// time.Sleep(100, "ms")
	setServerToListen()
	return
}

/*
Tmpl
*/

func setTmplToHandle() {
	for route, _ := range paramsRoute {
		// fmt.Println(route)
		http.HandleFunc(
			route,
			setTmpl)
	}
	// time.ShowTimeAndMsg("Tmpl set success")
	return
}

func setTmpl(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	route, ok := getRoute(req)
	if !ok {
		return
	}
	paramsIn, ok := execFromTmpl(route, req)
	paramsOutAppend := rfl.Params{}
	if ok {
		paramsOutAppend = entityBackEnd.ExecBackEnd(paramsIn)
	}
	ok = execSetToTmpl(route, resp, paramsOutAppend)
	if !ok {
		return
	}
	return
}

/*
FromTmpl
*/

func execFromTmpl(route string, req *http.Request) (paramsIn rfl.Params, ok bool) {
	paramsFromTmpl := req.URL.Query()
	// fmt.Println(paramsFromTmpl)
	// rfl.ShowType(paramsFromTmpl)
	if len(paramsFromTmpl) == 0 {
		return nil, false
	}
	paramsInDefault, ok := getParamsInDefault(route)
	if !ok {
		return nil, false
	}
	paramsIn = getParamsFromTmpl(paramsFromTmpl, paramsInDefault)
	return paramsIn, true
}

func getParamsInDefault(route string) (paramsInDefault rfl.Params, ok bool) {
	paramsInDefault, ok = paramsRoute[route]["paramsInDefault"]
	return paramsInDefault, ok
}

func getParamsFromTmpl(paramsFromTmpl url.Values, paramsInDefault rfl.Params) (paramsIn rfl.Params) {
	var valueFromTmpl string
	var valueIn interface{}
	paramsIn = rfl.Params{}
	for field, valueDefault := range paramsInDefault {
		valueFromTmpl = paramsFromTmpl.Get(field)
		// fmt.Println(valueFromTmpl)
		// rfl.ShowType(valueFromTmpl)
		if valueFromTmpl == "" {
			valueIn = valueDefault
		} else {
			valueIn = valueFromTmpl
		}
		paramsIn[field] = valueIn
	}
	// fmt.Println(paramsIn)
	return paramsIn
}

/*
ToTmpl
*/

func getRoute(req *http.Request) (route string, ok bool) {
	route = req.URL.Path
	// fmt.Println(route)
	ok = false
	for routeValid, _ := range paramsRoute {
		// fmt.Println(routeValid)
		if route == routeValid {
			ok = true
		}
	}
	if !ok {
		return "", false
	}
	return route, true
}

func execSetToTmpl(route string, resp http.ResponseWriter, paramsOutAppend rfl.Params) (ok bool) {
	pathTmpl, ok := getPathTmplByRoute(route)
	if !ok {
		return false
	}
	tmpl := getTmpl(pathTmpl)
	paramsOutDefault, ok := getParamsOutDefault(route)
	if !ok {
		return false
	}
	paramsOutAppendAll := getParamsOutAppendAll(paramsOutAppend)
	paramsOut := getParamsToTmpl(paramsOutDefault, paramsOutAppendAll)
	// fmt.Println(paramsOut)
	errTmplExec := tmpl.Execute(resp, paramsOut)
	err.ErrCheck(errTmplExec)
	return true
}

func getPathTmplByRoute(route string) (pathTmpl string, ok bool) {
	pathTmpl, ok = paramsRoute[route]["pathTmpl"]["pathTmpl"].(string)
	if !ok {
		return "", false
	}
	// fmt.Println(pathTmpl)
	return pathTmpl, true
}

func getTmpl(pathTmpl string) (tmpl *template.Template) {
	fsPublic, patternsTmpl := env.GetFSOfTmpl(pathTmpl)
	tmpl, errTmplParse := template.ParseFS(fsPublic, patternsTmpl)
	err.ErrCheck(errTmplParse)
	// rfl.ShowType(tmpl)
	return tmpl
}

func getParamsOutDefault(route string) (paramsOutDefault rfl.Params, ok bool) {
	paramsOutDefault, ok = paramsRoute[route]["paramsOutDefault"]
	return paramsOutDefault, ok
}

func getParamsOutAppendAll(paramsOutAppend rfl.Params) (paramsOutAppendAll rfl.Params) {
	paramsOutFromENV := env.GetParamsTmpl()
	paramsOutAppendAll = rfl.ParamsMerge(paramsOutFromENV, paramsOutAppend)
	return paramsOutAppendAll
}

func getParamsToTmpl(paramsOutDefault rfl.Params, paramsOutAppendAll rfl.Params) (paramsOut rfl.Params) {
	paramsOut = rfl.Params{}
	for field, valueDefault := range paramsOutDefault {
		valueOut := valueDefault
		valueOutAppend, ok := paramsOutAppendAll[field].(string)
		if ok {
			valueOut = valueOutAppend
		}
		paramsOut[field] = valueOut
	}
	// fmt.Println(paramsOutAppendAll)
	// fmt.Println(paramsOut)
	return paramsOut
}

/*
Explorer
*/

func OpenExplorer() {
	osName := env.GetOSName()
	// fmt.Println(osName)
	if osName != "windows" {
		return
	}
	command, _ := commands[osName]
	domain := env.GetDomain()
	url := fmt.Sprintf(
		"%s://%s/%s",
		"http",
		domain,
		"index")
	run := fmt.Sprintf(
		"%s %s",
		command,
		url)
	// fmt.Println(run)
	cmd := exec.Command(run)
	byteOutput, errCmd := cmd.Output()
	fmt.Println(byteOutput)
	err.ErrCheck(errCmd)
	msg := "Open Explorer Success"
	time.ShowTimeAndMsg(msg)
	return
}

/*
Listen
*/

func setServerToListen() {
	address := env.GetDomain()
	server := http.Server{
		Addr: address,
	}
	// time.ShowTimeAndMsg("Server listen success")
	msg := fmt.Sprintf(
		"Ready!Please open |%s/%s| in your browser to continue",
		address,
		"index")
	msg = fmt.Sprintf(
		"准备就绪!请在浏览器中打开|%s|继续操作",
		address)
	time.ShowTimeAndMsg(msg)
	errServerListen := server.ListenAndServe()
	err.ErrCheck(errServerListen)
	return
}
