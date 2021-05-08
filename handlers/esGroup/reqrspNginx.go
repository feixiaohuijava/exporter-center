package esGroup

import (
	"context"
	"exporter-center/config"
	"exporter-center/config/configStruct"
	"exporter-center/logs"
	"exporter-center/tool"
	"github.com/goinggo/mapstructure"
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func HandlerResRspNginx(gaugeIns *prometheus.GaugeVec, req *prometheus.Registry) http.Handler {
	logs.Infoln("do clear metric data before each metrics http request!")
	gaugeIns.Reset()

	var wg sync.WaitGroup
	var urlNginxConfig configStruct.UrlNginxConfig
	logNoneMesStruct := config.GetYamlConfig("config_urlnginx", &urlNginxConfig)
	err := mapstructure.Decode(logNoneMesStruct, &urlNginxConfig)
	logs.Infoln("获取到的配置地址:", urlNginxConfig)
	if err != nil {
		logs.Errorln("解析错误,原因:", err)
		return nil
	}

	reqRspNginx := urlNginxConfig.ReqRspNginx
	var client *elastic.Client
	logs.Infoln("es地址", reqRspNginx.EsUrl)

	currentTime := tool.GetCurTimeTsp()
	beforeOneMinuteTime := tool.GetBeforeMinTimeTsp(reqRspNginx.Cycle)

	client, err = elastic.NewClient(
		elastic.SetURL(reqRspNginx.EsUrl),
		elastic.SetSniff(false),
	)
	if err != nil {
		logs.Errorln("获取es客户端错误,原因:", err)
		return nil
	}
	for _, metricName := range reqRspNginx.MetricNames {
		logs.Infoln("metricName", metricName)
		if metricName == "request_time" {
			logs.Infoln("==========开始request_time")
			for indexRequestTimes, requestTime := range reqRspNginx.RequestTimes {
				logs.Infoln(requestTime, indexRequestTimes)
				wg.Add(1)
				go GetReqRspNginxData(gaugeIns, client, reqRspNginx, metricName, requestTime,
					beforeOneMinuteTime, currentTime, wg.Done)
			}
		} else if metricName == "response_status" {
			logs.Infoln("==========开始response_status")
			for indexResponseStatus, responseStatus := range reqRspNginx.ResponseStatus {
				logs.Infoln(responseStatus, indexResponseStatus)
				wg.Add(1)
				go GetReqRspNginxData(gaugeIns, client, reqRspNginx, metricName, responseStatus,
					beforeOneMinuteTime, currentTime, wg.Done)
			}
		} else {
			logs.Errorln("获取异常数据metricName", metricName)
			wg.Done()
			return nil
		}
	}
	wg.Wait()
	resultHandler := promhttp.HandlerFor(req, promhttp.HandlerOpts{})
	return resultHandler
}

func GetReqRspNginxData(gaugeIns *prometheus.GaugeVec, client *elastic.Client, reqRspNginx configStruct.ReqRspNginx,
	metricName string, reqRsp configStruct.Between, beforeOneMinuteTime time.Time, currentTime time.Time, done func()) {
	defer func() {
		if e := recover(); e != nil {
			logs.Errorln("捕获异常,错误原因", e)
			defer done()
		}
	}()
	var searchResult *elastic.SearchResult
	var err error
	ctx := context.Background()
	//logFilePathQuery := elastic.NewTermQuery("log.file.path.keyword", logFilePath)
	boolQuery := elastic.NewBoolQuery()
	//boolQuery.Must(logFilePathQuery)

	if metricName == "request_time" {
		var requestTimeQuery *elastic.RangeQuery
		if reqRsp.Max == "*" {
			requestTimeQuery = elastic.NewRangeQuery("request_time").Gte(reqRsp.Min)
		} else {
			requestTimeQuery = elastic.NewRangeQuery("request_time").Gte(reqRsp.Min).Lte(reqRsp.Max)
		}
		boolQuery.Must(requestTimeQuery)
	} else if metricName == "response_status" {
		responseQuery := elastic.NewRangeQuery("response").Gte(reqRsp.Min).Lte(reqRsp.Max)
		boolQuery.Must(responseQuery)
	} else {
		panic("出现异常metricName")
	}
	timeRangeFilter := elastic.NewRangeQuery("@timestamp").Gte(beforeOneMinuteTime).Lte(currentTime)
	boolQuery.Filter(timeRangeFilter)

	data := elastic.NewTermsAggregation().Field("log.file.path.keyword")

	searchResult, err = client.Search().
		Index(reqRspNginx.EsIndex).
		Query(boolQuery).From(0).Size(100).
		Aggregation("log.file.path.keyword", data).
		Do(ctx)
	if err != nil {
		logs.Infoln("查询出错!", err)
		defer done()
		return
	}
	logs.Infoln("searchResult", searchResult)

	//此处用于打印结果
	//type Msg struct {
	//	Uri         string  `json:"uri"`
	//	RequestTime float64 `json:"request_time"`
	//	Response    float64 `json:"response"`
	//}
	//var msg Msg
	//for _, item := range searchResult.Each(reflect.TypeOf(msg)) {
	//	t := item.(Msg)
	//	//logs.Infoln("查询出来结果:", t.Uri, t.RequestTime, t.Response)
	//	if metricName == "request_time" {
	//		gaugeIns.WithLabelValues(reqRspNginx.EsUrl+reqRspNginx.EsIndex, strconv.Itoa(reqRspNginx.Cycle), metricName,
	//			logFilePath, reqRsp.Min+" TO "+reqRsp.Max).Set(t.RequestTime)
	//	} else {
	//		gaugeIns.WithLabelValues(reqRspNginx.EsUrl+reqRspNginx.EsIndex, strconv.Itoa(reqRspNginx.Cycle), metricName,
	//			logFilePath, reqRsp.Min+" TO "+reqRsp.Max).Set(t.Response)
	//	}
	//}

	agg, found := searchResult.Aggregations.Terms("log.file.path.keyword")
	if !found {
		logs.Infoln("没有聚合结果")
		defer done()
		return
	}
	for _, item := range agg.Buckets {
		app := item.Key
		docCount := item.DocCount
		if metricName == "request_time" {
			gaugeIns.WithLabelValues(reqRspNginx.EsUrl+reqRspNginx.EsIndex, strconv.Itoa(reqRspNginx.Cycle), metricName,
				app.(string), reqRsp.Min+" TO "+reqRsp.Max).Set(float64(docCount))
		} else {
			gaugeIns.WithLabelValues(reqRspNginx.EsUrl+reqRspNginx.EsIndex, strconv.Itoa(reqRspNginx.Cycle), metricName,
				app.(string), reqRsp.Min+" TO "+reqRsp.Max).Set(float64(docCount))
		}
	}

	defer done()
}
