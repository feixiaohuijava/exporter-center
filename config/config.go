package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
)

// 传入配置文件的名字,返回的就是yamlConfig这个结构体中携带配置的结构体
func GetYamlConfig(yamlName string, yamlConfig interface{}) interface{} {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	sysType := runtime.GOOS
	filePath := ""
	if sysType == "linux" || sysType == "darwin" {
		filePath = dir + "/config/" + yamlName + ".yaml"
	}
	if sysType == "windows" {
		filePath = dir + "\\config\\" + yamlName + ".yaml"
	}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = yaml.UnmarshalStrict(yamlFile, &yamlConfig)
	if err != nil {
		panic(err)
	}
	return yamlConfig
}
