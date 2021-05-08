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

func HandlerAvUvPv(gaugeIns *prometheus.GaugeVec, req *prometheus.Registry) http.Handler {
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

	avUvPvConfig := urlNginxConfig.AvUvPv

	var client *elastic.Client
	logs.Infoln("es地址", avUvPvConfig.EsUrl)

	currentTime := tool.GetCurTimeTsp()
	beforeOneMinuteTime := tool.GetBeforeMinTimeTsp(avUvPvConfig.Cycle)

	client, err = elastic.NewClient(
		elastic.SetURL(avUvPvConfig.EsUrl),
		elastic.SetSniff(false),
	)
	if err != nil {
		logs.Errorln("获取es客户端错误,原因:", err)
		return nil
	}

	for _, metricName := range avUvPvConfig.MetricNames {
		if metricName == "avg(request_time)" {
			wg.Add(1)
			go requestAvg(gaugeIns, client, avUvPvConfig, beforeOneMinuteTime, currentTime, wg.Done)
		} else if metricName == "uv" {
			wg.Add(1)
			go requestUv(gaugeIns, client, avUvPvConfig, beforeOneMinuteTime, currentTime, wg.Done)
		} else if metricName == "pv" {
			wg.Add(1)
			go requestPv(gaugeIns, client, avUvPvConfig, beforeOneMinuteTime, currentTime, wg.Done)
		}
	}

	wg.Wait()
	resultHandler := promhttp.HandlerFor(req, promhttp.HandlerOpts{})
	return resultHandler
}

func requestAvg(gaugeIns *prometheus.GaugeVec, client *elastic.Client, avUvPvConfig configStruct.AvUvPv,
	beforeOneMinuteTime time.Time, currentTime time.Time, done func()) {
	defer func() {
		if e := recover(); e != nil {
			logs.Errorln("捕获异常,错误原因", e)
			defer done()
		}
	}()

	var searchResult *elastic.SearchResult
	var err error
	ctx := context.Background()

	boolQuery := elastic.NewBoolQuery()
	timeRangeFilter := elastic.NewRangeQuery("@timestamp").Gte(beforeOneMinuteTime).Lte(currentTime)
	boolQuery.Filter(timeRangeFilter)

	data := elastic.NewAvgAggregation().Field("request_time")

	searchResult, err = client.Search().
		Index(avUvPvConfig.EsIndex).
		Query(boolQuery).
		Size(0).
		Aggregation("request_time", data).Do(ctx)
	if err != nil {
		logs.Infoln("查询出错!", err)
		defer done()
		return
	}
	agg, found := searchResult.Aggregations.Avg("request_time")
	if !found {
		logs.Infoln("没有聚合结果")
		defer done()
		return
	}
	gaugeIns.WithLabelValues(avUvPvConfig.EsUrl+avUvPvConfig.EsIndex, strconv.Itoa(avUvPvConfig.Cycle),
		"avg(request_time)", "empty").Set(*agg.Value)
	defer done()
}

func requestUv(gaugeIns *prometheus.GaugeVec, client *elastic.Client, avUvPvConfig configStruct.AvUvPv,
	beforeOneMinuteTime time.Time, currentTime time.Time, done func()) {
	defer func() {
		if e := recover(); e != nil {
			logs.Errorln("捕获异常,错误原因", e)
			defer done()
		}
	}()

	var searchResult *elastic.SearchResult
	var err error
	ctx := context.Background()

	boolQuery := elastic.NewBoolQuery()
	timeRangeFilter := elastic.NewRangeQuery("@timestamp").Gte(beforeOneMinuteTime).Lte(currentTime)
	boolQuery.Filter(timeRangeFilter)

	data := elastic.NewCardinalityAggregation().Field("geoip.ip")

	searchResult, err = client.Search().
		Index(avUvPvConfig.EsIndex).
		Query(boolQuery).
		Size(0).
		Aggregation("uv", data).Do(ctx)
	if err != nil {
		logs.Infoln("查询出错!", err)
		defer done()
		return
	}
	agg, found := searchResult.Aggregations.Cardinality("uv")
	if !found {
		logs.Infoln("没有聚合结果")
		defer done()
		return
	}
	gaugeIns.WithLabelValues(avUvPvConfig.EsUrl+avUvPvConfig.EsIndex, strconv.Itoa(avUvPvConfig.Cycle),
		"uv", "empty").Set(*agg.Value)
	defer done()
}

func requestPv(gaugeIns *prometheus.GaugeVec, client *elastic.Client, avUvPvConfig configStruct.AvUvPv,
	beforeOneMinuteTime time.Time, currentTime time.Time, done func()) {
	defer func() {
		if e := recover(); e != nil {
			logs.Errorln("捕获异常,错误原因", e)
			defer done()
		}
	}()

	//var searchResult *elastic.SearchResult
	var err error
	ctx := context.Background()

	boolQuery := elastic.NewBoolQuery()

	boolQuery.Must(elastic.NewMatchAllQuery())

	timeRangeFilter := elastic.NewRangeQuery("@timestamp").Gte(beforeOneMinuteTime).Lte(currentTime)
	boolQuery.Filter(timeRangeFilter)

	data := elastic.NewTermsAggregation().Field("log.file.path.keyword")
	searchResult, err := client.Search().
		Index(avUvPvConfig.EsIndex).
		Query(boolQuery).Aggregation("pv", data).
		Size(0).Do(ctx)
	if err != nil {
		logs.Infoln("查询出错!", err)
		defer done()
		return
	}
	agg, found := searchResult.Aggregations.Terms("pv")
	if !found {
		logs.Infoln("没有聚合结果!")
		defer done()
		return
	}
	for _, item := range agg.Buckets {
		app := item.Key
		docCount := item.DocCount
		gaugeIns.WithLabelValues(avUvPvConfig.EsUrl+avUvPvConfig.EsIndex, strconv.Itoa(avUvPvConfig.Cycle),
			"pv", app.(string)).Set(float64(docCount))
	}

	defer done()
}
