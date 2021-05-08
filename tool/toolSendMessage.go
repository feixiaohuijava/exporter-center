package tool

import (
	"bytes"
	"encoding/json"
	"exporter-center/logs"
)

type RobotResult struct {
	Results []Robot `json:"results"`
}

type Robot struct {
	Id               int      `json:"id"`
	RobotType        string   `json:"robot_type"`
	RobotUrlLabel    []string `json:"robot_url_label"`
	RobotName        string   `json:"robot_name"`
	RobotUrl         string   `json:"robot_url"`
	RobotDescription string   `json:"robot_description"`
	Channel          string   `json:"channel"`
	Project          string   `json:"project"`
	CreatedTime      string   `json:"createdTime"`
	UpdateTime       string   `json:"updateTime"`
}

func SendMessage(robotContent string, robotName string) {
	logs.Infoln("开始进行发送消息操作")
	robotString := FindRobot(robotName)
	logs.Infoln("获取操作平台的机器人地址:", robotString)
	var tempRobot = RobotResult{}
	json.Unmarshal([]byte(robotString), &tempRobot)
	if len(tempRobot.Results) != 1 {
		logs.Infoln("机器人值不唯一!")
		panic("获取机器人有有误!")
	}
	type DDPlayload struct {
		Msgtype  string            `json:"msgtype"`
		Markdown map[string]string `json:"markdown"`
		At       map[string]string `json:"at"`
	}
	for _, item := range tempRobot.Results {
		var tempMarkdown = map[string]string{"title": "监控", "text": "firing" + robotContent}
		var tempAt = map[string]string{"atMobiles": "None", "isAtAll": "True"}
		var payload = make(map[string]interface{})
		payload["msgtype"] = "markdown"
		payload["markdown"] = tempMarkdown
		payload["at"] = tempAt
		bytesData, err := json.Marshal(payload)
		Checkerr(err)
		HttpPost(item.RobotUrl, bytes.NewReader(bytesData), AuthPost{})
	}
}
