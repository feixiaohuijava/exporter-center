package toolElasticSearch

import (
	"exporter-center/logs"
	"github.com/olivere/elastic/v7"
	"strings"
)

func GetClientEs(url string) *elastic.Client {
	var err error
	var esClient *elastic.Client
	if strings.Contains(url, "http://ops-es.jms.com") || strings.Contains(url, "http://customer-es.jms.com") {
		esClient, err = elastic.NewClient(
			elastic.SetURL(url),
			elastic.SetSniff(false),
			elastic.SetBasicAuth("elastic", "yl123456c0m"),
		)
	} else {
		esClient, err = elastic.NewClient(
			elastic.SetURL(url),
			elastic.SetSniff(false),
		)
	}
	if err != nil {
		logs.Errorln("获取客户端出错,原因:", err)
		return nil
	}
	return esClient

}
