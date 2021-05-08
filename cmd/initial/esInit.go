package initial

import (
	"exporter-center/logs"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	gaugeIns *prometheus.GaugeVec
	registry *prometheus.Registry
)

func GetGaugeRegistry(gaugeName string, helpName string, metricLabels []string) (*prometheus.GaugeVec, *prometheus.Registry) {
	logs.Infoln("开始进行" + gaugeName + "和registry的初始化操作!")
	gaugeIns = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: gaugeName, Help: helpName}, metricLabels)
	registry = prometheus.NewRegistry()
	registry.MustRegister(gaugeIns)
	return gaugeIns, registry
}
