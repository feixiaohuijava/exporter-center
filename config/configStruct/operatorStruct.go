package configStruct

type OperatorNotScanConfig struct {
	OperatorNotScan OperatorNotScan
}

type OperatorNotScan struct {
	Url           string `yaml:"url"`
	EffectTimeMax int    `yaml:"effectTimeMax"`
	EffectTimeMin int    `yaml:"effectTimeMin"`
	User          []User
}

type User struct {
	Username string `yaml:"username"`
	Phone    string `yaml:"phone"`
}
