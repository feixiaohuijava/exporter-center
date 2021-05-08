package tool

import (
	"exporter-center/config"
	"exporter-center/config/configStruct"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/goinggo/mapstructure"
	"testing"
)

func Test_Tds(t *testing.T) {
	var aliyunConfig configStruct.AliyunConfig
	aliyunStruct := config.GetYamlConfig("config_aliyun", &aliyunConfig)
	err := mapstructure.Decode(aliyunStruct, &aliyunConfig)
	if err != nil {
		panic(err)
	}
	client, err := sdk.NewClientWithAccessKey("cn-shanghai", aliyunConfig.Aliyun.AccessKeyId, aliyunConfig.Aliyun.AccessKeySecret)
	if err != nil {
		// Handle exceptions
		panic(err)
	}
	fmt.Println(client)

	request := requests.NewCommonRequest() // 构造一个公共请求。
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "DescribeSynchronizationJobs"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["PageSize"] = "100"
	request.QueryParams["PageNum"] = "1"
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
}
