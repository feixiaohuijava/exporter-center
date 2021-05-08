package tool

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

func GetToken(loginUrl string) string {
	byteLogin := []byte(`{ "username": "ylops", "password": "YLpoim09" }`)
	payload := bytes.NewBuffer(byteLogin)
	body, statusCode := HttpPost(loginUrl, payload, AuthPost{})
	if statusCode != 200 {
		panic("登录获取token请求异常:" + body + "返回状态码:" + strconv.Itoa(statusCode))
	}
	var data = []byte(body)
	jsonDataToken := jsoniter.Get(data, "results", "token").ToString()
	return jsonDataToken
}

func FindRobot(robotName string) string {
	token := GetToken("http://devops-admin.jtexpress.com.cn/api/account/login/")
	var byteData []byte
	payload := bytes.NewBuffer(byteData)
	var params []map[string]string
	demo := map[string]string{"robot_name": robotName}
	params = append(params, demo)
	body, statusCode := HttpGetJson("http://devops-admin.jtexpress.com.cn/api/webhookmanager/robotviewset/", payload, token, params)
	if statusCode != 200 {
		panic("异常")
	}
	return body
}

func GetAppEs(appName string) string {
	token := GetToken("http://devops-admin.jtexpress.com.cn/api/account/login/")
	var byteData []byte
	payload := bytes.NewBuffer(byteData)
	var params []map[string]string
	demo := map[string]string{"applicationName": appName, "environmentName": "pro"}
	params = append(params, demo)
	body, statusCode := HttpGetJson("http://devops-admin.jtexpress.com.cn/api/cmdb/applicationinstance/", payload, token, params)
	if statusCode != 200 {
		panic("异常,获取接口的状态不是200")
	}
	return body
}
