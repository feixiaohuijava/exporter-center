package tool

import (
	"exporter-center/config"
	"exporter-center/config/configStruct"
	"github.com/goinggo/mapstructure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

func GetDbConnection() *gorm.DB {

	var databaseConfig configStruct.DatabaseConfig
	var err error
	var db *gorm.DB
	dbStruct := config.GetYamlConfig("config_db", &databaseConfig)
	err = mapstructure.Decode(dbStruct, &databaseConfig)

	if err != nil {
		panic(err)
	}
	dsn := databaseConfig.DataBase.User +
		":" + databaseConfig.DataBase.Password +
		"@tcp(" + databaseConfig.DataBase.Host +
		":" + strconv.Itoa(databaseConfig.DataBase.Port) +
		")/" + databaseConfig.DataBase.Name +
		"?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
