package toolElasticSearch

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"reflect"
	"testing"
	"time"
)

func Test_ElasticSearch(t *testing.T) {

	currentTime := time.Now()
	beforeOneMinuteTime := time.Now().Add(-time.Minute * 120)

	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL("http://ops-es.jms.com"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "yl123456c0m"),
	)
	if err != nil {
		panic(err)
	}
	var searchResult *elastic.SearchResult
	type Msg struct {
		AppName string `json:"appName"`
		Logmsg  string `json:"logmsg"`
		Logger  string `json:"logger"`
		Level   string `json:"level"`
	}
	// 这个NewMatchPhraseQuery很诡异,如果是[rocketMQ消息推送], shardingKey:PULL_FAILED可以查出，如果是[rocketMQ消息推送], shardingKey:P这种就查不出
	// 看来日志里面的内容也要区分,但是NewMatchQuery就可以实现。 请一定要理解这两者区别
	termQuery := elastic.NewMatchPhraseQuery("level", "ERROR")
	//demo := elastic.NewMatchPhraseQuery("logger", "com.yl.ops.schedule.ReceiptListHandler")
	matchPhraseQuery := elastic.NewMatchQuery("logmsg", "自定义异常")
	timeRangeFilter := elastic.NewRangeQuery("@timestamp").Gte(beforeOneMinuteTime).Lte(currentTime)
	boolQ := elastic.NewBoolQuery()
	boolQ.Filter(matchPhraseQuery)
	boolQ.Filter(termQuery)
	boolQ.Filter(timeRangeFilter)
	data := elastic.NewTermsAggregation().Field("appName.keyword")
	searchResult, err = client.Search().
		Index("jms-ops-uat-applog-*").
		Query(boolQ).
		From(0).
		Size(200).
		Aggregation("appName.keyword", data).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	var msg Msg
	for _, item := range searchResult.Each(reflect.TypeOf(msg)) {
		t := item.(Msg)
		fmt.Println(t.Logger, t.AppName, t.Level)
	}
	agg, found := searchResult.Aggregations.Terms("appName.keyword")
	if !found {
		fmt.Println("没有聚合结果")
	}
	fmt.Println(agg)
	for _, item := range agg.Buckets {
		app := item.Key
		docCount := item.DocCount
		fmt.Println(app, docCount)
	}

}
