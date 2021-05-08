package phoneGroup

import (
	"bytes"
	"exporter-center/config"
	"exporter-center/config/configStruct"
	"exporter-center/logs"
	"exporter-center/models/webhookmanager"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/gin-gonic/gin"
	"github.com/goinggo/mapstructure"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type RequestBody struct {
	AlertCall AlertCall `json:"alertCall" form:"alertCall"`
}

type AlertCall struct {
	AlertName string `json:"alertName"`
	VoiceName string `json:"voiceName"`
	Project   string `json:"project"`
	User      []User `json:"user"`
}

type User struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

// @Tags 电话接口
// @Summary 电话接口
// @Description {"alertCall": {"alertName": "监控告警", "voiceName": "运维平台", "project": "平台运维", "user": [{"username": "feixiaohui", "phone": "13048856263"}]}}
// @Accept application/json
// @Produce application/json
// @Param alertCall body string false "参数"
// @Param Content-Type header string false "application/json"
// @Param Authorization header string false "JWT Token"
// @Router /phone/callphone [post]
func CallPhone(c *gin.Context, dbConnection *gorm.DB) {
	var requestBody RequestBody
	var err error
	if c.ShouldBind(&requestBody) == nil {
		if requestBody.AlertCall.AlertName == "" || requestBody.AlertCall.VoiceName == "" || requestBody.AlertCall.Project == "" {
			logs.Errorln("电话失败,参数错误!", requestBody)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "电话失败,参数错误!"})
			return
		}
		for _, item := range requestBody.AlertCall.User {
			if item.Username == "" || item.Phone == "" {
				logs.Errorln("参数错误!", requestBody.AlertCall.User)
				c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "电话失败,参数错误!"})
				return
			}
		}
		var aliyunConfig configStruct.AliyunConfig
		aliyunStruct := config.GetYamlConfig("config_aliyun", &aliyunConfig)
		err = mapstructure.Decode(aliyunStruct, &aliyunConfig)
		if err != nil {
			logs.Errorln("获取阿里云客户端出错,原因:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "msg": "获取阿里云客户端出错!"})
		}
		defer func() {
			// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
			if e := recover(); e != nil {
				// 某个服务调用协程报错，可以在这里打印一些错误日志
				logs.Errorln("错误了,开始进行捕获操作!", e)
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "msg": e})
			}
		}()
		PhoneCall(aliyunConfig.Aliyun, dbConnection, requestBody)
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "results": "电话成功!"})
	}
}

func PhoneCall(aliyun configStruct.Aliyun, dbConnection *gorm.DB, requestBody RequestBody) {

	client, err := dyvmsapi.NewClientWithAccessKey(aliyun.RegionId, aliyun.AccessKeyId, aliyun.AccessKeySecret)
	if err != nil {
		logs.Errorln("获取阿里云客户端错误:", err)
		panic(err)
	}
	request := dyvmsapi.CreateSingleCallByTtsRequest()
	request.Scheme = "https"
	request.TtsCode = aliyun.TtsCode
	request.TtsParam = "{'alertName':'" + requestBody.AlertCall.VoiceName + "'}"
	for _, user := range requestBody.AlertCall.User {

		var buff bytes.Buffer
		var result bool

		request.CalledNumber = user.Phone
		response, err := client.SingleCallByTts(request)
		if err != nil {
			logs.Errorln("执行电话接口出错,原因:", err)
			panic(err)
		}
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
		// 进行记录操作
		insertIns := webhookmanager.WebhookmanagerPhonerecord{
			Username:    user.Username,
			Phone:       user.Phone,
			AlertName:   requestBody.AlertCall.AlertName,
			Result:      result,
			ResultMsg:   resultMsg,
			Project:     requestBody.AlertCall.Project,
			VoiceName:   requestBody.AlertCall.VoiceName,
			CreatedTime: time.Now(),
		}
		webhookmanager.InsertPhoneRecord(dbConnection, insertIns)
		logs.Infoln(user.Username)
		logs.Infoln(user.Phone)
		logs.Infoln(requestBody.AlertCall.AlertName)
		logs.Infoln(requestBody.AlertCall.VoiceName)
		logs.Infoln(requestBody.AlertCall.Project)
		logs.Infoln("执行电话回调Code:", response.Code, "执行电话回调RequestId:", response.RequestId, "执行电话回调Message:", response.Message, "执行电话回调CallId:", response.CallId)
	}
}
