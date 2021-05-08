package esExporer

type EsExporter struct {
	LogMessageAppName []LogMessageAppName
	ProBigdataLog     ProBigdataLog
}

type Url struct {
	Name  string `yaml:"name"`
	Index string `yaml:"index"`
}

type MetricName struct {
	Name   string `yaml:"name"`
	Fields []Field
}

type Field struct {
	Key       string `yaml:"key"`
	Value     string `yaml:"value"`
	FieldType string `yaml:"fieldtype"`
	Logic     bool   `yaml:"logic"`
}
