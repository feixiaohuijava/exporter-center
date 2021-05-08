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

func HandlerUrlNginx(gaugeIns *prometheus.GaugeVec, req *prometheus.Registry) http.Handler {
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
	currentTime := tool.GetCurTimeTsp()
	beforeOneMinuteTime := tool.GetBeforeMinTimeTsp(urlNginxConfig.UrlNginx.Cycle)

	wg.Add(len(urlNginxConfig.UrlNginx.Urls))

	var client *elastic.Client
	logs.Infoln("es地址", urlNginxConfig.UrlNginx.EsUrl)
	client, err = elastic.NewClient(
		elastic.SetURL(urlNginxConfig.UrlNginx.EsUrl),
		elastic.SetSniff(false),
	)
	if err != nil {
		logs.Errorln("获取es客户端错误,原因:", err)
		return nil
	}

	for _, item := range urlNginxConfig.UrlNginx.Urls {
		go GetUrlNginxData(gaugeIns, client, urlNginxConfig.UrlNginx, item, beforeOneMinuteTime, currentTime, wg.Done)
	}

	wg.Wait()
	resultHandler := promhttp.HandlerFor(req, promhttp.HandlerOpts{})
	return resultHandler
}

func GetUrlNginxData(gaugeIns *prometheus.GaugeVec, client *elastic.Client, urlNginx configStruct.UrlNginx,
	urlRequest configStruct.UrlRequest, beforeOneMinuteTime time.Time, currentTime time.Time, done func()) {
	defer func() {
		if e := recover(); e != nil {
			logs.Errorln("捕获异常,错误原因", e)
			defer done()
		}
	}()

	var searchResult *elastic.SearchResult
	var err error
	ctx := context.Background()
	logs.Infoln(urlNginx.EsUrl)
	uriTermQuery := elastic.NewTermQuery("uri.keyword", urlRequest.Url)
	requestTimeQuery := elastic.NewRangeQuery("request_time").Gte(urlRequest.RequestTime)
	timeRangeFilter := elastic.NewRangeQuery("@timestamp").Gte(beforeOneMinuteTime).Lte(currentTime)
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(uriTermQuery)
	boolQuery.Must(requestTimeQuery)
	boolQuery.Must(timeRangeFilter)
	data := elastic.NewTermsAggregation().Field("uri.keyword")
	searchResult, err = client.Search(urlNginx.EsIndex).
		Query(boolQuery).
		From(0).Size(100).
		Aggregation("uri.keyword", data).
		Pretty(true).Do(ctx)
	if err != nil {
		logs.Infoln("查询出错!")
		return
	}

	logs.Infoln(searchResult)

	//此处用于打印结果
	//type Msg struct {
	//	Uri string `json:"uri"`
	//	RequestTime  float64 `json:"request_time"`
	//}
	//var msg Msg
	//for _, item := range searchResult.Each(reflect.TypeOf(msg)) {
	//	t := item.(Msg)
	//	logs.Info(t.Uri, t.RequestTime)
	//}

	agg, found := searchResult.Aggregations.Terms("uri.keyword")
	if !found {
		logs.Errorln("没有找到聚合数据!")
	}
	for _, bucket := range agg.Buckets {
		bucketValue := bucket.Key
		logs.Infoln(bucketValue.(string), bucket.DocCount)
		gaugeIns.WithLabelValues(bucketValue.(string), strconv.FormatFloat(urlRequest.RequestTime, 'f', 3, 64),
			urlNginx.EsUrl+"/"+urlNginx.EsIndex).
			Set(float64(bucket.DocCount))
	}
	defer done()
}
