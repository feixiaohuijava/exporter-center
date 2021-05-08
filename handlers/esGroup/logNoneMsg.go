package esGroup

import (
	"context"
	"encoding/json"
	"exporter-center/config"
	"exporter-center/config/configStruct"
	"exporter-center/logs"
	"exporter-center/tool"
	"github.com/goinggo/mapstructure"
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thedevsaddam/gojsonq/v2"
	"net/http"
	"reflect"
	"sync"
	"time"
)

type SearchIndex struct {
	Name          string `json:"name"`
	ElasticSearch ElasticSearch
}

type ElasticSearch struct {
	Name     string `json:"name"`
	Domain   string `json:"domain"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func LogNoneMessage(gaugeIns *prometheus.GaugeVec, req *prometheus.Registry) http.Handler {
	logs.Infoln("do clear metric data before each metrics http request!")
	gaugeIns.Reset()
	logs.Infoln("开始获取时间")
	currentTime := time.Now()
	beforeOneMinuteTime := time.Now().Add(-time.Minute * 1)
	var wg sync.WaitGroup
	var logNoneMsgConfig configStruct.LogNoneMsgConfig
	logNoneMesStruct := config.GetYamlConfig("config_lognone", &logNoneMsgConfig)
	err := mapstructure.Decode(logNoneMesStruct, &logNoneMsgConfig)
	if err != nil {
		logs.Errorln("解析配置文件出错!原因:", err)
	}
	logs.Infoln("从配置文件config_es_expoter获取到的配置:", logNoneMsgConfig.LogNoneMsg)

	wg.Add(len(logNoneMsgConfig.LogNoneMsg.Apps))
	for _, app := range logNoneMsgConfig.LogNoneMsg.Apps {
		go AppGoRoutine(gaugeIns, app, beforeOneMinuteTime, currentTime, wg.Done)
	}
	wg.Wait()
	resultHandler := promhttp.HandlerFor(req, promhttp.HandlerOpts{})
	return resultHandler
}

func AppGoRoutine(gaugeIns *prometheus.GaugeVec, app string, beforeOneMinuteTime time.Time, currentTime time.Time, done func()) {

	defer func() {
		// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
		if e := recover(); e != nil {
			// 某个服务调用协程报错，可以在这里打印一些错误日志
			logs.Errorln("错误了,开始进行捕获操作!")
			logs.Errorln(e)
			defer done()
		}
	}()

	// 根据app查询es的url
	appData := tool.GetAppEs(app)
	resultData := gojsonq.New().FromString(appData).From("results.[0]").From("searchIndex").Get()
	var searchIndex SearchIndex
	data, _ := json.Marshal(resultData)
	json.Unmarshal(data, &searchIndex)
	logs.Infoln("从cmdb获取的es数据", searchIndex)
	if searchIndex == (SearchIndex{}) {
		logs.Errorln("应用名", app)
	}
	GetLogNoneData(gaugeIns, app, searchIndex, beforeOneMinuteTime, currentTime)
	defer done()
}

func GetLogNoneData(gaugeIns *prometheus.GaugeVec, app string, searchIndex SearchIndex, beforeOneMinuteTime time.Time, currentTime time.Time) {
	indexName := searchIndex.Name
	esUrl := searchIndex.ElasticSearch.Domain
	username := searchIndex.ElasticSearch.Username
	password := searchIndex.ElasticSearch.Password
	var searchResult *elastic.SearchResult
	ctx := context.Background()
	var client *elastic.Client
	var err error
	if username != "" && password != "" {
		logs.Infoln("账号密码都存在")
		logs.Infoln(esUrl, username, password)
		client, err = elastic.NewClient(
			elastic.SetURL("http://"+esUrl),
			elastic.SetSniff(false),
			elastic.SetBasicAuth(username, password),
		)
	} else {
		logs.Infoln("账号密码不存在")
		client, err = elastic.NewClient(
			elastic.SetURL(esUrl),
			elastic.SetSniff(false),
		)
	}
	if err != nil {
		logs.Errorln("获取es客户端错误,其地址是:", esUrl)
		logs.Errorln("获取es客户端错误,其索引是:", indexName)
		logs.Errorln("app名称是:", app)
		logs.Errorln("获取es客户端错误,原因:", err)
	}

	//进行term查询,此处term是代表匹配查询,这里的appName需要添加keyword,否则value值会被拆分
	termQuery := elastic.NewTermQuery("appName.keyword", app)
	timeRangeFilter := elastic.NewRangeQuery("@timestamp").Gte(beforeOneMinuteTime).Lte(currentTime)
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(termQuery)
	boolQuery.Must(timeRangeFilter)
	searchResult, err = client.Search(indexName).Query(boolQuery).From(0).Size(100).Pretty(true).Do(ctx)
	if err != nil {
		logs.Errorln("查询结果出错,原因:", err)
	}
	type Msg struct {
		AppName string `json:"appName"`
		Logmsg  string `json:"logmsg"`
		Logger  string `json:"logger"`
		Level   string `json:"level"`
	}
	var msg Msg
	dataLength := searchResult.Each(reflect.TypeOf(msg))
	if len(dataLength) != 0 {
		gaugeIns.WithLabelValues(app, "http://"+esUrl+"/"+indexName).Set(float64(len(dataLength)))
	} else {
		gaugeIns.WithLabelValues(app, "http://"+esUrl+"/"+indexName).Set(float64(0))
	}
	//查看每条数据
	//for index, item := range searchResult.Each(reflect.TypeOf(msg)) {
	//	t := item.(Msg)
	//	fmt.Println(t.Logger, t.AppName, t.Level, index)
	//}
}
