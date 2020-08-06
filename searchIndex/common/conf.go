package common

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Conf *configuration

var ConfPath string

//profile variables
type configuration struct {
	SonicHost       string `yaml:"sonic_host"`
	SonicPassword   string `yaml:"sonic_password"`
	RootPath        string `yaml:"root_path"`
	PostPath        string `yaml:"post_path"`
	IndexAllInit    bool   `yaml:"index_all_init"`
	SonicCollection string `yaml:"sonic_collection"`
	SonicBucket     string `yaml:"sonic_bucket"`
}

func (c *configuration) getConf() *configuration {

	path := flag.String("p", "conf.yaml", "配置文件的路径")

	flag.Parse()

	yamlFile, err := ioutil.ReadFile(*path)
	if err != nil {
		log.Error(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Error(err.Error())
	}
	return c
}

func init() {
	log.Info("初始化配置文件")
	conf := configuration{}
	Conf = conf.getConf()
}
