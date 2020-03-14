package config

import (
	"fmt"
	"go-im/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var GlobalConfig = &model.AppConfig{}

/**
初始化配置
*/
func Initconfig() error {
	yamlFile, err := ioutil.ReadFile("config/app.yaml")
	if err != nil {
		fmt.Println("==读取app.yaml文件失败=========:", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, GlobalConfig)
	if err != nil {
		fmt.Println("==解析app.yaml文件失败=========")
		return err
	}
	return nil
}
