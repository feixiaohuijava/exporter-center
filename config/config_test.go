package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func Test_DBConfig(t *testing.T) {

	sysType := runtime.GOOS
	fmt.Println(sysType)
	//var demo YmalConfig

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(dir)
	filePath := dir + "\\config_operatornotscan.yaml"
	fmt.Println(filePath)

	yamlFile, err := ioutil.ReadFile(filePath)
	fmt.Println(yamlFile)
	if err != nil {
		panic(err)
	}
	//err = yaml.UnmarshalStrict(yamlFile, &demo)
	if err != nil {
		panic(err)
	}
	//fmt.Println(demo.Database.Name)
	//fmt.Println(demo.Watchins)
	//for _, temp := range demo.Watchins {
	//	fmt.Println(temp.Cluster)
	//	fmt.Println(temp.Namespace)
	//}
}

// get current project's root path
// return path not contain the exec file
func GetProjectRoot() string {
	var (
		path string
		err  error
	)
	defer func() {
		if err != nil {
			panic(fmt.Sprintf("GetProjectRoot error :%+v", err))
		}
	}()
	path, err = filepath.Abs(filepath.Dir(os.Args[0]))
	return path
}

// get configure file path
//func GetConfPath() string {
//	return GetProjectRoot() + confDir
//}
