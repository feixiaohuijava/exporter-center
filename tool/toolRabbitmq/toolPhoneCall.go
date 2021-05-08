package tool

import (
	"bytes"
	"encoding/json"
	"exporter-center/config/configStruct"
	"exporter-center/models/webhookmanager"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/chenhg5/collection"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func PhoneCall(yamlConfig configStruct.AliyunConfig, phoneUsers []PhoneUser, logger *logrus.Logger, voiceName string, dbConnection *gorm.DB) {
	client, err := dyvmsapi.NewClientWithAccessKey(yamlConfig.Aliyun.RegionId, yamlConfig.Aliyun.AccessKeyId, yamlConfig.Aliyun.AccessKeySecret)
	if err != nil {
		logger.Error(err.Error())
	}
	request := dyvmsapi.CreateSingleCallByTtsRequest()
	request.Scheme = "https"
	request.TtsCode = yamlConfig.Aliyun.TtsCode
	// voiceName 这里前后要加上单引号,不然就不是字符串了
	request.TtsParam = "{'alertName':'" + voiceName + "'}"
	result := removeDuplicate(phoneUsers)
	var calledArray []string
	for _, phoneUser := range result {
		if !collection.Collect(calledArray).Contains(phoneUser.Phone) {
			var buff bytes.Buffer
			var result bool
			calledArray = append(calledArray, phoneUser.Phone)
			logger.Info("准备要打电话名字是:" + phoneUser.UserName + " ,电话是:" + phoneUser.Phone)
			request.CalledNumber = phoneUser.Phone
			response, err := client.SingleCallByTts(request)
			if err != nil {
				logger.Error(err.Error())
			}
			logger.Info("response is %#v\n", response)
			// 电话记录
			buff.WriteString("Code:")
			buff.WriteString(response.Code)
			buff.WriteString(",")
			buff.WriteString("RequestId:")
			buff.WriteString(response.RequestId)
			buff.WriteString(",")
			buff.WriteString("Message:")
			buff.WriteString(response.Message)
			buff.WriteString(",")
			buff.WriteString("CallId:")
			buff.WriteString(response.CallId)
			resultMsg := buff.String()
			if response.Code == "OK" {
				result = true
			} else {
				result = false
			}
			insertIns := webhookmanager.WebhookmanagerPhonerecord{
				Username:    phoneUser.UserName,
				Phone:       phoneUser.Phone,
				AlertName:   "",
				Result:      result,
				ResultMsg:   resultMsg,
				Project:     "",
				VoiceName:   voiceName,
				CreatedTime: time.Now(),
			}
			webhookmanager.InsertPhoneRecord(dbConnection, insertIns)
		} else {
			logger.Info("重复打电话名字是:" + phoneUser.UserName + " ,电话是:" + phoneUser.Phone)
		}
	}
}

func removeDuplicate(personList []PhoneUser) []PhoneUser {
	resultMap := map[string]bool{}
	for _, v := range personList {
		data, _ := json.Marshal(v)
		resultMap[string(data)] = true
	}
	var result []PhoneUser
	for k := range resultMap {
		var t PhoneUser
		json.Unmarshal([]byte(k), &t)
		result = append(result, t)
	}
	return result
}
