package configStruct

type AliyunConfig struct {
	Aliyun Aliyun
}
type Aliyun struct {
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	RegionId        string `yaml:"regionId"`
	TtsCode         string `yaml:"ttsCode"`
}
