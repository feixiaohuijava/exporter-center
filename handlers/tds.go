package handlers

import (
	"encoding/json"
	"exporter-center/config"
	"exporter-center/config/configStruct"
	"exporter-center/logs"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/goinggo/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"strings"
)

type TdsResult struct {
	TotalRecordCount         int `json:"TotalRecordCount"`
	PageNumber               int `json:"PageNumber"`
	PageRecordCount          int `json:"PageRecordCount"`
	SynchronizationInstances []EachIns
	RequestId                string `json:"RequestId"`
}

type EachIns struct {
	Status                 string `json:"Status"`
	SynchronizationJobName string `json:"SynchronizationJobName"`
	CreateTime             string `json:"CreateTime"`
	SynchronizationJobId   string `json:"SynchronizationJobId"`
}

func TdsHandler(gaugeIns *prometheus.GaugeVec, req *prometheus.Registry) http.Handler {
	logs.Infoln("do clear metric data before each metrics http request!")
	gaugeIns.Reset()

	var aliyunConfig configStruct.AliyunConfig
	var err error
	aliyunStruct := config.GetYamlConfig("config_aliyun", &aliyunConfig)
	err = mapstructure.Decode(aliyunStruct, &aliyunConfig)
	if err != nil {
		logs.Errorln("发生错误,原因:", err)
		return nil
	}
	var client *sdk.Client
	client, err = sdk.NewClientWithAccessKey("cn-shanghai", aliyunConfig.Aliyun.AccessKeyId, aliyunConfig.Aliyun.AccessKeySecret)
	if err != nil {
		logs.Errorln("发生错误,原因:", err)
		return nil
	}
	// 根据业务需要,目前就只是请求两次,每次请求100条
	totalPage := []int{1, 2}
	for _, page := range totalPage {
		request := requests.NewCommonRequest() // 构造一个公共请求。
		request.Method = "POST"
		request.Scheme = "https"
		request.Domain = "dts.aliyuncs.com"
		request.Version = "2020-01-01"
		request.ApiName = "DescribeSynchronizationJobs"
		request.QueryParams["RegionId"] = "cn-shanghai"
		request.QueryParams["PageSize"] = "100"
		request.QueryParams["PageNum"] = strconv.Itoa(page)
		response, err := client.ProcessCommonRequest(request)
		if err != nil {
			logs.Errorln("调用tds列表api接口出错,原因:", err)
			logs.Errorln("发生错误,原因:", err)
			return nil
		}
		var tdsResult TdsResult
		json.Unmarshal([]byte(response.GetHttpContentString()), &tdsResult)
		for _, item := range tdsResult.SynchronizationInstances {
			env := ""
			if strings.HasPrefix(item.SynchronizationJobName, "uat") {
				env = "uat"
			} else if strings.HasPrefix(item.SynchronizationJobName, "pro") {
				env = "pro"
			} else {
				env = "test"
			}
			if item.Status == "InitializeFailed" || item.Status == "Failed" {
				gaugeIns.WithLabelValues(item.SynchronizationJobName, item.SynchronizationJobId, item.Status, env).Set(1)
			} else {
				gaugeIns.WithLabelValues(item.SynchronizationJobName, item.SynchronizationJobId, item.Status, env).Set(0)
			}
		}
	}
	resultHandler := promhttp.HandlerFor(req, promhttp.HandlerOpts{})
	return resultHandler
}
