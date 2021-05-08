package configStruct

type DatabaseConfig struct {
	DataBase DataBase
	Rabbitmq Rabbitmq
}
type DataBase struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type Rabbitmq struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
