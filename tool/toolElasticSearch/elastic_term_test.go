package toolElasticSearch

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"reflect"
	"testing"
	"time"
)

func Test_ElasticTermSearch(t *testing.T) {
	var searchResult *elastic.SearchResult
	currentTime := time.Now()
	beforeOneMinuteTime := time.Now().Add(-time.Minute * 1)

	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL("http://ops-es.jms.com"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "yl123456c0m"),
	)
	if err != nil {
		panic(err)
	}

	//进行term查询,此处term是代表匹配查询,这里的appName需要添加keyword,否则value值会被拆分,如果要多个条件查询就得使用boolQuery
	termQuery := elastic.NewTermQuery("appName.keyword", "yl-web-operatingplatform")
	timeRangeFilter := elastic.NewRangeQuery("@timestamp").Gte(beforeOneMinuteTime).Lte(currentTime)
	//searchResult, err = client.Search().Index("jms-ops-pro-applog-*").Query(timeRangeFilter).Query(termQuery).From(0).Size(1000).Pretty(true).Do(ctx)
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(termQuery)
	boolQuery.Must(timeRangeFilter)
	searchResult, err = client.Search("jms-ops-pro-applog-*").Query(boolQuery).Pretty(true).Do(ctx)
	if err != nil {
		panic(err)
	}
	type Msg struct {
		AppName string `json:"appName"`
		Logmsg  string `json:"logmsg"`
		Logger  string `json:"logger"`
		Level   string `json:"level"`
	}
	var msg Msg
	for index, item := range searchResult.Each(reflect.TypeOf(msg)) {
		t := item.(Msg)
		fmt.Println(t.Logger, t.AppName, t.Level, index)
	}
}
