package configStruct

type UrlNginxConfig struct {
	UrlNginx    UrlNginx
	ReqRspNginx ReqRspNginx
	AvUvPv      AvUvPv
}

type UrlNginx struct {
	EsUrl   string `yaml:"esUrl"`
	EsIndex string `yaml:"esIndex"`
	Cycle   int    `yaml:"cycle"`
	Urls    []UrlRequest
}

type UrlRequest struct {
	Url         string  `yaml:"username"`
	RequestTime float64 `yaml:"requestTime"`
}

type ReqRspNginx struct {
	MetricNames    []string  `yaml:"metricName"`
	Cycle          int       `yaml:"cycle"`
	EsUrl          string    `yaml:"esUrl"`
	EsIndex        string    `yaml:"esIndex"`
	RequestTimes   []Between `yaml:"requestTimes"`
	ResponseStatus []Between `yaml:"responseStatus"`
}

type Between struct {
	Min string `yaml:"min"`
	Max string `yaml:"max"`
}

type AvUvPv struct {
	MetricNames []string `yaml:"metricNames"`
	Cycle       int      `yaml:"cycle"`
	EsUrl       string   `yaml:"esUrl"`
	EsIndex     string   `yaml:"esIndex"`
}
