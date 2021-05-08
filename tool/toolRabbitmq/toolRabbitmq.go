package tool

import (
	"bytes"
	"encoding/json"
	"exporter-center/config/configStruct"
	"exporter-center/models/webhookmanager"
	"exporter-center/tool"
	"github.com/goinggo/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/gojsonq/v2"
	"gorm.io/gorm"
	"strconv"

	"exporter-center/config"
	"github.com/streadway/amqp"
)

type PhoneUser struct {
	UserName string `json:"username"`
	Phone    string `json:"phone"`
}
type CustomeConsumer struct {
	AlertRuleId int      `json:"alertRule_id"`
	CallResult  []string `json:"callResult"`
}

type RegularConsumer struct {
	AlertRuleId int        `json:"alertRule_id"`
	CallResult  CallResult `json:"callResult"`
}

type CallResult struct {
	Dev     []PhoneUser `json:"dev"`
	Ops     []PhoneUser `json:"ops"`
	Test    []PhoneUser `json:"test"`
	Product []PhoneUser `json:"product"`
}

type Rabbitmq struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func failOnError(err error, msg string, logger *logrus.Logger) {
	if err != nil {
		logger.Error("%s: %s", msg, err)
	}
}

func GetMsgFromRabbitmq(logger *logrus.Logger, dbConnection *gorm.DB) {

	defer func() {
		// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
		if e := recover(); e != nil {
			// 某个服务调用协程报错，可以在这里打印一些错误日志
			logger.Errorln("从rabbitmq获取消息消费错误了,开始进行捕获操作!")
			logger.Errorln(e)
		}
	}()

	var databaseConfig configStruct.DatabaseConfig
	var aliyunConfig configStruct.AliyunConfig
	var err error

	db := tool.GetDbConnection()

	// 争对db配置进行解析
	dbStruct := config.GetYamlConfig("config_db", &databaseConfig)
	err = mapstructure.Decode(dbStruct, &databaseConfig)
	if err != nil {
		logger.Errorln("在rabbitmq consumer中解析db配置文件出错!")
	}
	// 针对aliyun配置进行解析
	aliyunStruct := config.GetYamlConfig("config_aliyun", &aliyunConfig)
	err = mapstructure.Decode(aliyunStruct, &aliyunConfig)

	if err != nil {
		logger.Errorln("在rabbitmq conisumer中解析aliyun配置出错!")
	}

	conn, err := amqp.Dial("amqp://" + databaseConfig.Rabbitmq.User + ":" +
		databaseConfig.Rabbitmq.Password + "@" + databaseConfig.Rabbitmq.Host + ":" + strconv.Itoa(databaseConfig.Rabbitmq.Port) + "/")
	failOnError(err, "连接rabbitmq失败,", logger)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "打开channel失败", logger)
	defer ch.Close()

	//q, err := ch.QueueDeclare(
	//	"devops", // name
	//	false,   // durable
	//	false,   // delete when unused
	//	false,   // exclusive
	//	false,   // no-wait
	//	nil,     // arguments
	//)
	//failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		"phone_center", // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		logger.Errorln("Failed to register a consumer", err)
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			msgFromRabbit := string(d.Body)
			logger.Info("开始从rabbitmq接收消息:", msgFromRabbit)
			// get callType
			callType := gojsonq.New().FromString(msgFromRabbit).Find("callType")
			logger.Info("获取到的callType:", callType)
			var phoneUser []PhoneUser
			if callType == "regular" {
				// 获取电话
				var regularConsumer RegularConsumer
				json.Unmarshal([]byte(msgFromRabbit), &regularConsumer)
				//根据alertRuleId查找该应用
				alertRule := webhookmanager.FindAlertById(db, regularConsumer.AlertRuleId, logger)
				noticeKeyValue := alertRule.NotiteKeyValue
				voiceName := alertRule.VoiceName
				logger.Info("noticeKeyValue", noticeKeyValue)
				// 添加开发
				if alertRule.CallDevFlag {
					for _, dev := range regularConsumer.CallResult.Dev {
						phoneUser = append(phoneUser, PhoneUser{dev.UserName, dev.Phone})
					}
				}
				// 添加测试
				if alertRule.CallTestFlag {
					for _, test := range regularConsumer.CallResult.Test {
						phoneUser = append(phoneUser, PhoneUser{test.UserName, test.Phone})
					}
				}
				// 添加产品
				if alertRule.CallProductFlag {
					for _, product := range regularConsumer.CallResult.Product {
						phoneUser = append(phoneUser, PhoneUser{product.UserName, product.Phone})
					}
				}
				// do phone action
				PhoneCall(aliyunConfig, phoneUser, logger, voiceName, dbConnection)
			} else if callType == "custome" {
				var customeConsumer CustomeConsumer
				json.Unmarshal([]byte(msgFromRabbit), &customeConsumer)
				logger.Info(customeConsumer.CallResult)
				jsonArrayCallResult, _ := json.Marshal(customeConsumer.CallResult)
				phoneUser := FindPhoneByUser(jsonArrayCallResult)
				logger.Info("phoneUser:", phoneUser)
				//根据alertRuleId查找该应用
				alertRule := webhookmanager.FindAlertById(db, customeConsumer.AlertRuleId, logger)
				voiceName := alertRule.VoiceName
				PhoneCall(aliyunConfig, phoneUser, logger, voiceName, dbConnection)
			}
		}
	}()
	logger.Info("等待接收消息")
	<-forever
}

func FindPhoneByUser(callUser []byte) []PhoneUser {
	var resultMap map[string][]PhoneUser
	var byteData []byte
	payload := bytes.NewBuffer(byteData)
	var params []map[string]string
	demo := map[string]string{"callUser": string(callUser)}
	params = append(params, demo)
	body, statusCode := tool.HttpGetJson("http://devops-admin.jtexpress.com.cn/api/account/findphonebyuser/", payload, "", params)
	if statusCode != 200 {
		panic("异常")
	}
	json.Unmarshal([]byte(body), &resultMap)
	return resultMap["results"]
}
