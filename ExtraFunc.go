package main

import(
	"fmt"
	m "itchat4go/model"
	"io/ioutil"
	"github.com/go-yaml/yaml"
)


var (
	uuid       string
	err        error
	loginMap   m.LoginMap
	contactMap map[string]m.User
	groupMap   map[string][]m.User 
)


type conf struct {
	Time     	string `yaml:"time"`
	Date        string `yaml:"date"`
	Everyday  	int    `yaml:"everyDay"`
	WechatName  string `yaml:"wechatName"`
	CityName    string `yaml:"cityName"`
	NameList    []string `yaml:"list"`
}


func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

//panicErr 输出错误信息
func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}