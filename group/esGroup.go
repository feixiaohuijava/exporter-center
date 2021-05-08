package group

import (
	"exporter-center/cmd/initial"
	"exporter-center/handlers"
	"exporter-center/handlers/esGroup"
	"github.com/gin-gonic/gin"
)

func EsGroup(esRouter *gin.RouterGroup) {
	esRouter.GET("/probigdata/metrics", func(context *gin.Context) {
		labels := []string{"url"}
		guageIns, registryIns := initial.GetGaugeRegistry("probigdata", "pro bigdata metrics", labels)
		wrapper := gin.WrapH(esGroup.ProBigdataEslog(guageIns, registryIns))
		wrapper(context)
	})
	esRouter.GET("/notoperater/metrics", func(context *gin.Context) {
		labels := []string{"url", "ModuleCode", "ModuleName", "Count"}
		gaugeIns, registryIns := initial.GetGaugeRegistry("NotOperatedIn15minitues", "Not Operated In 15 minitues", labels)
		wrapper := gin.WrapH(handlers.NotOperatorHandler(gaugeIns, registryIns))
		wrapper(context)
	})
	esRouter.GET("/logmessage/appname/metrics", func(context *gin.Context) {
		labels := []string{"appName", "cycle", "metricName", "query", "url"}
		gaugeIns, registryIns := initial.GetGaugeRegistry("logMessageAppName", "logMessageAppName", labels)
		wrapper := gin.WrapH(esGroup.LogMessageSdk(gaugeIns, registryIns))
		wrapper(context)
	})
	esRouter.GET("lognone/appname/metrics", func(context *gin.Context) {
		labels := []string{"appName", "url"}
		gaugeIns, registryIns := initial.GetGaugeRegistry("logNone", "logNone", labels)
		wrapper := gin.WrapH(esGroup.LogNoneMessage(gaugeIns, registryIns))
		wrapper(context)
	})
	esRouter.GET("/dts/metrics", func(context *gin.Context) {
		labels := []string{"SynchronizationJobName", "SynchronizationJobId", "Status", "env"}
		gaugeIns, registryIns := initial.GetGaugeRegistry("dts", "dts", labels)
		wrapper := gin.WrapH(handlers.TdsHandler(gaugeIns, registryIns))
		wrapper(context)
	})
	esRouter.GET("/urlnginxtimeout/metrics", func(context *gin.Context) {
		labels := []string{"url", "requestTime", "esUrl"}
		gaugeIns, registryIns := initial.GetGaugeRegistry("urlNginx", "urlNginx", labels)
		wrapper := gin.WrapH(esGroup.HandlerUrlNginx(gaugeIns, registryIns))
		wrapper(context)
	})
	esRouter.GET("/request_response/metrics", func(context *gin.Context) {
		labels := []string{"url", "cycle", "metricName", "log_file_path", "input_value"}
		gaugeIns, registryIns := initial.GetGaugeRegistry("request_response", "request_response", labels)
		wrapper := gin.WrapH(esGroup.HandlerResRspNginx(gaugeIns, registryIns))
		wrapper(context)
	})
	esRouter.GET("/requesttime_uv_pv/metrics", func(context *gin.Context) {
		labels := []string{"url", "cycle", "metricName", "log_file_path"}
		gaugeIns, registryIns := initial.GetGaugeRegistry("requesttime_uv_pv", "es export", labels)
		wrapper := gin.WrapH(esGroup.HandlerAvUvPv(gaugeIns, registryIns))
		wrapper(context)
	})
	// for alert-center
	esRouter.GET("/alertcenter/metrics", func(context *gin.Context) {
		labels := []string{"url", "cycle", "metricName"}
		gaugeIns, registryIns := initial.GetGaugeRegistry("alertcenter", "alertcenter", labels)
		wrapper := gin.WrapH(esGroup.HandlerAlertCenter(gaugeIns, registryIns))
		wrapper(context)
	})
}
