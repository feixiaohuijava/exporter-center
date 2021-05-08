package handlers

import (
	"bytes"
	"encoding/json"
	"exporter-center/config"
	"exporter-center/config/configStruct"
	"exporter-center/logs"
	"exporter-center/tool"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/goinggo/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

func NotOperatorHandler(gaugeIns *prometheus.GaugeVec, req *prometheus.Registry) http.Handler {
	//最开始做好清除动作
	logs.Infoln("do clear metric data before each metrics http request!")
	gaugeIns.Reset()
	//获取配置
	var operatorNotScanConfig configStruct.OperatorNotScanConfig
	operatorNotScanStruct := config.GetYamlConfig("config_operatornotscan", &operatorNotScanConfig)
	err := mapstructure.Decode(operatorNotScanStruct, &operatorNotScanConfig)
	if err != nil {
		logs.Errorln("发生错误,原因:", err)
		return nil
	}
	GetOpsScanMonitor(gaugeIns, operatorNotScanConfig)
	resultHandler := promhttp.HandlerFor(req, promhttp.HandlerOpts{})
	return resultHandler
}

type OpsResult struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data []EachData `json:"data"`
}

type EachData struct {
	Count      int    `json:"count"`
	ModuleCode string `json:"moduleCode"`
	ModuleName string `json:"moduleName"`
}

//{"code":1,"msg":"请求成功","data":[{"moduleCode":"express_collection","moduleName":"快件揽收","count":0},{"moduleCode":"receipt","moduleName":"入仓扫描","count":0},{"moduleCode":"arrival","moduleName":"到件扫描","count":0},{"moduleCode":"collection_proxy","moduleName":"代理点代收","count":0},{"moduleCode":"cabinet_deliveryout","moduleName":"快递柜取出","count":0},{"moduleCode":"cabinet_entering","moduleName":"快递柜入库","count":0},{"moduleCode":"cabinet_sign","moduleName":"快递柜出库","count":0},{"moduleCode":"deliveryout","moduleName":"派件出仓","count":0},{"moduleCode":"loading","moduleName":"装车扫描","count":0},{"moduleCode":"unloading","moduleName":"卸车扫描","count":0},{"moduleCode":"send","moduleName":"发件扫描","count":0},{"moduleCode":"packing_list","moduleName":"建包扫描-包内件","count":0},{"moduleCode":"remain_storage","moduleName":"留仓扫描","count":0},{"moduleCode":"unpack","moduleName":"拆包扫描","count":0},{"moduleCode":"sign","moduleName":"签收扫描","count":0},{"moduleCode":"uncar","moduleName":"解车扫描","count":0},{"moduleCode":"pack_send_scan_list","moduleName":"封发车扫描","count":0}],"fail":false,"succ":true}
func GetOpsScanMonitor(gaugeIns *prometheus.GaugeVec, operatorNotScanConfig configStruct.OperatorNotScanConfig) {
	var byteData []byte
	var noticeFlag bool
	var tempData []EachData
	currentTimeStr := time.Now().Format("2006-01-02 15:04:05")
	noticeFlag = false
	palyload := bytes.NewBuffer(byteData)
	var params []map[string]string
	params = append(params, map[string]string{"dateTime": currentTimeStr})
	body, statusCode := tool.HttpGetJson(operatorNotScanConfig.OperatorNotScan.Url, palyload, "", params)
	logs.Infoln("获取请求的数据body:", body)
	if statusCode != 200 {
		logs.Errorln("获取" + operatorNotScanConfig.OperatorNotScan.Url + "的接口返回状态码:" + string(rune(statusCode)))
	}
	v := OpsResult{}
	json.Unmarshal([]byte(body), &v)
	code := v.Code
	msg := v.Msg
	data := v.Data
	logs.Infoln("获取OpsResult的数据:", v.Code)
	logs.Infoln("获取OpsResult的数据:", v.Msg)
	if code == 1 && msg == "请求成功" && len(data) > 0 {
		for _, item := range data {
			if item.Count == 0 {
				if item.ModuleCode != "" && item.ModuleCode != "collection_proxy" && item.ModuleCode != "remain_storage" &&
					item.ModuleCode != "receipt" && item.ModuleCode != "deliveryout" && item.ModuleCode != "cabinet_entering" &&
					item.ModuleCode != "express_collection" {
					noticeFlag = true
					tempData = append(tempData, EachData{item.Count, item.ModuleCode, item.ModuleName})
				}
			}
		}
	}
	currentHour := time.Now().Hour()
	if len(tempData) > 0 {
		if currentHour <= operatorNotScanConfig.OperatorNotScan.EffectTimeMin || currentHour >= operatorNotScanConfig.OperatorNotScan.EffectTimeMax {
			if noticeFlag {
				var tempBuffer bytes.Buffer
				tempBuffer.WriteString("**15分钟扫描类型未操作预警** \n>" + "时间: " + currentTimeStr + "\n\n")
				for _, item := range tempData {
					gaugeIns.WithLabelValues(operatorNotScanConfig.OperatorNotScan.Url, item.ModuleCode, item.ModuleName, strconv.Itoa(item.Count)).Set(1)
					tempBuffer.WriteString("模块名称:" + item.ModuleName + "\n\n" + "模块Code:" + item.ModuleCode + "\n\n" + "<font color='warning'>--------------------</font>\n\n")
				}
				robotContent := tempBuffer.String()
				// 同时进行打电话和发送消息通知
				tool.SendMessage(robotContent, "pro-操作平台")
				err := PhoneCall(operatorNotScanConfig, "未扫描操作")
				if err != nil {
					return
				}
			} else {
				gaugeIns.WithLabelValues(operatorNotScanConfig.OperatorNotScan.Url, "empty", "empty", "empty").Set(0)
			}
		} else {
			gaugeIns.WithLabelValues(operatorNotScanConfig.OperatorNotScan.Url, "empty", "empty", "empty").Set(0)
		}
	}
}

func PhoneCall(operatorNotScanConfig configStruct.OperatorNotScanConfig, voiceName string) error {
	var aliyunConfig configStruct.AliyunConfig
	aliyunStruct := config.GetYamlConfig("config_aliyun", &aliyunConfig)
	err := mapstructure.Decode(aliyunStruct, &aliyunConfig)
	if err != nil {
		logs.Errorln("发生错误,原因:", err)
		return err
	}
	client, err := dyvmsapi.NewClientWithAccessKey(aliyunConfig.Aliyun.RegionId, aliyunConfig.Aliyun.AccessKeyId, aliyunConfig.Aliyun.AccessKeySecret)
	if err != nil {
		logs.Errorln(err.Error())
		return err
	}
	request := dyvmsapi.CreateSingleCallByTtsRequest()
	request.Scheme = "https"
	request.TtsCode = aliyunConfig.Aliyun.TtsCode
	request.TtsParam = "{'alertName':'" + voiceName + "'}"
	for _, phoneUser := range operatorNotScanConfig.OperatorNotScan.User {
		request.CalledNumber = phoneUser.Phone
		response, err := client.SingleCallByTts(request)
		if err != nil {
			logs.Errorln(err.Error())
			logs.Errorln(response)
			return err
		}
	}
	return nil
}
