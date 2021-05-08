package configStruct

type AdConfig struct {
	Ad Ad
}
type Ad struct {
	Address      string `yaml:"address"`
	Basedn       string `yaml:"basedn"`
	BindUsername string `yaml:"bindUsername"`
	BindPassword string `yaml:"bindPassword"`
}
