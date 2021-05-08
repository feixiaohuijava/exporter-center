package esGroup

import (
	"bytes"
	"context"
	"exporter-center/config"
	"exporter-center/config/configStruct/esExporer"
	"exporter-center/logs"
	"exporter-center/tool"
	"exporter-center/tool/toolElasticSearch"
	"github.com/goinggo/mapstructure"
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func ProBigdataEslog(gaugeIns *prometheus.GaugeVec, req *prometheus.Registry) http.Handler {
	log.Infoln("do clear metric data before each metrics http request!")
	gaugeIns.Reset()
	logs.Infoln("开始获取时间")

	// 获取配置
	var esExporterConfig esExporer.EsExporter
	esExporerStruct := config.GetYamlConfig("config_es_expoter", &esExporterConfig)
	err := mapstructure.Decode(esExporerStruct, &esExporterConfig)
	if err != nil {
		logs.Errorln("读取配置出错了!", err)
		panic(err)
	}

	var wg sync.WaitGroup
	urls := esExporterConfig.ProBigdataLog.Urls
	cycle := esExporterConfig.ProBigdataLog.Cycle
	metricNames := esExporterConfig.ProBigdataLog.MetricNames
	for _, url := range urls {
		lteTime := tool.GetCurTimeTsp()
		gteTime := tool.GetBeforeMinTimeTsp(cycle)
		// 获取es客户端
		esClient := toolElasticSearch.GetClientEs(url.Name)
		if esClient == nil {
			logs.Errorln("获取es客户端为空!")
			continue
		}
		for _, metricName := range metricNames {
			wg.Add(1)
			go GetEsProbigdata(gaugeIns, esClient, gteTime, lteTime, metricName, url, cycle, wg.Done)
		}
	}
	resultHandler := promhttp.HandlerFor(req, promhttp.HandlerOpts{})
	return resultHandler
}

func GetEsProbigdata(gaugeIns *prometheus.GaugeVec, esClient *elastic.Client, gteTime time.Time, lteTime time.Time,
	metricName esExporer.MetricName, url esExporer.Url, cycle int, done func()) {
	//捕获错误
	defer func() {
		// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
		if e := recover(); e != nil {
			// 某个服务调用协程报错，可以在这里打印一些错误日志
			logs.Errorln("错误了,开始进行捕获操作!", e)
			defer done()
		}
	}()

	ctx := context.Background()
	var searchResult *elastic.SearchResult
	var query string
	var buff bytes.Buffer
	type Msg struct {
		AppName string `json:"appName"`
		Logmsg  string `json:"logmsg"`
		Logger  string `json:"logger"`
		Level   string `json:"level"`
	}
	name := metricName.Name
	fields := metricName.Fields
	// must查询
	boolQ := elastic.NewBoolQuery()
	for num, field := range fields {
		// 这个NewMatchPhraseQuery很诡异,如果是[rocketMQ消息推送], shardingKey:PULL_FAILED可以查出，如果是[rocketMQ消息推送], shardingKey:P这种就查不出
		// 看来日志里面的内容也要区分,但是NewMatchQuery就可以实现。 请一定要理解这两者区别
		if field.FieldType == "fuzzy" {
			logs.Infoln("start fuzzy filter")
			//boolQ.Filter(elastic.NewMatchPhraseQuery(field.Key, field.Value))
			//boolQ.Filter(elastic.NewMatchNoneQuery())
			if field.Logic == true {
				logs.Infoln("fuzzy true")
				boolQ.Must(elastic.NewMatchPhraseQuery(field.Key, field.Value))
			} else {
				logs.Infoln("fuzzy false")
				boolQ.MustNot(elastic.NewMatchPhraseQuery(field.Key, field.Value))
				buff.WriteString("NOT ")
			}
			buff.WriteString(field.Key)
			buff.WriteString("=")
			buff.WriteString(field.Value)
			if num != (len(fields) - 1) {
				buff.WriteString(" AND ")
			}
		}
		if field.FieldType == "specific" {
			if field.Logic == true {
				logs.Infoln("specific true")
				boolQ.Must(elastic.NewMatchQuery(field.Key, field.Value))
			} else {
				logs.Infoln("specific false")
				boolQ.MustNot(elastic.NewMatchQuery(field.Key, field.Value))
				buff.WriteString("NOT ")
			}
			buff.WriteString(field.Key)
			buff.WriteString(":")
			buff.WriteString(field.Value)
			if num != (len(fields) - 1) {
				buff.WriteString(" AND ")
			}
		}

	}
	query = buff.String()
	timeRangeFilter := elastic.NewRangeQuery("@timestamp").Gte(gteTime).Lte(lteTime)
	boolQ.Filter(timeRangeFilter)
	data := elastic.NewTermsAggregation().Field("appName.keyword")
	searchResult, err := esClient.Search().
		Index(url.Index).
		Query(boolQ).
		From(0).
		Size(200).
		Aggregation("appName.keyword", data).
		Do(ctx)
	if err != nil {
		logs.Infoln("es查询出错,原因:", err)
	}
	//此处用于打印结果
	//var msg Msg
	//for _, item := range searchResult.Each(reflect.TypeOf(msg)) {
	//	t := item.(Msg)
	//	logger.Info(t.Logger, t.AppName, t.Level)
	//}
	agg, found := searchResult.Aggregations.Terms("appName.keyword")
	if !found {
		logs.Infoln("没有聚合结果")
	}
	for _, item := range agg.Buckets {
		app := item.Key
		docCount := item.DocCount
		gaugeIns.WithLabelValues(app.(string), strconv.Itoa(cycle), name, query, url.Name+"/"+url.Index).Set(float64(docCount))
	}
	defer done()

}
