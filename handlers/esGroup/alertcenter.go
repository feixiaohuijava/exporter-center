package esGroup

import (
	"context"
	"exporter-center/logs"
	"exporter-center/tool"
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"reflect"
	"strconv"
)

func HandlerAlertCenter(gaugeIns *prometheus.GaugeVec, req *prometheus.Registry) http.Handler {
	logs.Infoln("do clear metric data before each metrics http request!")
	gaugeIns.Reset()

	var esClient *elastic.Client
	var err error
	cycle := 5
	esUrl := "http://10.30.6.146:9200"
	esIndex := "devops-*"

	esClient, err = elastic.NewClient(
		elastic.SetURL(esUrl),
		elastic.SetSniff(false),
	)
	if err != nil {
		logs.Errorln("获取es客户端错误,原因:", err)
		return nil
	}

	currentTime := tool.GetCurTimeTsp()
	// cycle is 1
	beforeOneMinuteTime := tool.GetBeforeMinTimeTsp(cycle)

	var searchResult *elastic.SearchResult

	ctx := context.Background()
	boolQuery := elastic.NewBoolQuery()
	timeRangeFilter := elastic.NewRangeQuery("@timestamp").Gte(beforeOneMinuteTime).Lte(currentTime)
	boolQuery.Filter(timeRangeFilter)

	//logFileFilter := elastic.NewMatchQuery("log.file.path", "alertCenter.log")
	//
	logFileFilter := elastic.NewTermQuery("log.file.path.keyword", "/data/devopslog/alertlog/alertCenter.log")
	boolQuery.Must(logFileFilter)

	searchResult, err = esClient.Search().Index(esIndex).Query(boolQuery).From(0).Size(200).Do(ctx)
	if err != nil {
		logs.Errorln("es search err, the reason was:", err)
	}

	//此处用于打印结果
	type Msg struct {
		Message string `json:"message"`
	}
	logs.Infoln("start print result")
	var msg Msg
	//for _, item := range searchResult.Each(reflect.TypeOf(msg)) {
	//	t := item.(Msg)
	//	logs.Infoln(t.Message)
	//	gaugeIns.WithLabelValues(esUrl+esIndex, strconv.Itoa(cycle), "alertcenter").Set(1)
	//}
	lengthMsg := len(searchResult.Each(reflect.TypeOf(msg)))
	if lengthMsg != 0 {
		// send message to dingding
		tool.SendMessage("alert center", "devops专属")
	}
	gaugeIns.WithLabelValues(esUrl+esIndex, strconv.Itoa(cycle), "alertcenter").Set(float64(lengthMsg))
	resultHandler := promhttp.HandlerFor(req, promhttp.HandlerOpts{})
	return resultHandler
}
