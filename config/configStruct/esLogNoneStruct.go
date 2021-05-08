package configStruct

type LogNoneMsgConfig struct {
	LogNoneMsg LogNoneMsg
}

type LogNoneMsg struct {
	Apps       []string `yaml:"apps"`
	Cycle      int      `yaml:"cycle"`
	MetricName string   `yaml:"metricName"`
}
