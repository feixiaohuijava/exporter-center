package esGroup

//
//import (
//	"bytes"
//	"encoding/json"
//	"exporter-center/config"
//	"exporter-center/config/configStruct"
//	"exporter-center/tool"
//	"github.com/goinggo/mapstructure"
//	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/client_golang/prometheus/promhttp"
//	"github.com/sirupsen/logrus"
//	"github.com/thedevsaddam/gojsonq/v2"
//	"net/http"
//	"strconv"
//	"strings"
//	"sync"
//)
//
//type Buckets struct {
//	Bucket []Bucket
//}
//type Bucket struct {
//	Key      string `json:"key"`
//	DocCount int    `json:"doc_count"`
//}
//
//func LogMessageAppName(gaugeIns *prometheus.GaugeVec, req *prometheus.Registry, logger *logrus.Logger) http.Handler {
//	//最开始做好清除动作
//	logger.Info("do clear metric data before each metrics http request!")
//	gaugeIns.Reset()
//	var wg sync.WaitGroup
//	//获取配置
//	var logMessageAppNameConfig configStruct.LogMessageAppNameConfig
//	logMessageAppNameStruct := config.GetYamlConfig("config_es_expoter", &logMessageAppNameConfig)
//	err := mapstructure.Decode(logMessageAppNameStruct, &logMessageAppNameConfig)
//	if err != nil {
//		panic(err)
//	}
//	logger.Info(logMessageAppNameConfig.LogMessageAppName)
//	for _, item := range logMessageAppNameConfig.LogMessageAppName {
//		Urls := item.Urls
//		cycle := item.Cycle
//		metricNames := item.MetricNames
//		for _, eachUrls := range Urls {
//			for _, metricName := range metricNames {
//				wg.Add(1)
//				//go GetEsData(gaugeIns, eachUrls, cycle, metricName.Name, metricName.Logmsgs, logger, wg.Done)
//			}
//		}
//	}
//	wg.Wait()
//	resultHandler := promhttp.HandlerFor(req, promhttp.HandlerOpts{})
//	return resultHandler
//}
//
//func GetEsData(gaugeIns *prometheus.GaugeVec, url string, cycle int, metricName string, logmsg string, logger *logrus.Logger, done func()) {
//	gtTime := "now-" + strconv.Itoa(cycle) + "m"
//	var buf bytes.Buffer
//	buf.WriteString(`{"track_total_hits": "true", "query": {"bool": {
//            "must": [{"query_string": {"analyze_wildcard": "true", "query": "`)
//	buf.WriteString(logmsg)
//	buf.WriteString(`"}},
//                     {"range": {"@timestamp": {"gt": `)
//	buf.WriteString(`"`)
//	buf.WriteString(gtTime)
//	buf.WriteString(`",`)
//	buf.WriteString(`"lt": "now"}}}]}},
//        "aggs": {"groups": {"terms": {"field": "appName.keyword", "size": 200, "order": {"_count": "desc"}}}},
//        "size": 0}`)
//	payload := bytes.NewBuffer([]byte(buf.String()))
//	var authPost tool.AuthPost
//	if strings.Contains(url, "http://ops-es.jms.com") || strings.Contains(url, "http://customer.jms.com") {
//		authPost.Username = "elastic"
//		authPost.Password = "yl123456c0m"
//	}
//	logger.Info("用户名密码:", authPost.Username, authPost.Password)
//	returnData, statusCode := tool.HttpPost(url, payload, authPost)
//	if statusCode != 200 {
//		logger.Error("请求状态码不是200,而是", statusCode)
//		panic("请求失败")
//	}
//	var bucket []Bucket
//	//bucketsInterface, err  := gojsonq.New().FromString(returnData).FindR("aggregations.groups.buckets")
//	//bucketsInterface := gojsonq.New().FromString(returnData).From("aggregations").From("groups").From("buckets").Select("key", "doc_count").Get()
//	bucketsInterface := gojsonq.New().FromString(returnData).From("aggregations").From("groups").From("buckets").Get()
//	logger.Info("bucketsInterface", bucketsInterface)
//	data, _ := json.Marshal(bucketsInterface)
//	//fmt.Println(string(data))
//	json.Unmarshal(data, &bucket)
//	//logger.Info(bucket)
//	var key string
//	var docCount int
//	for _, eachBucket := range bucket {
//		key = eachBucket.Key
//		docCount = eachBucket.DocCount
//		logger.Info(key)
//		logger.Info(docCount)
//		gaugeIns.WithLabelValues(key, strconv.Itoa(cycle), metricName, logmsg, url).Set(float64(docCount))
//	}
//	defer done()
//}
