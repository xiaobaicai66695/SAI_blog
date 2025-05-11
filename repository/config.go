package repository

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
}

type Database struct {
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Name      string `yaml:"name"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parse_time"`
	Loc       string `yaml:"loc"`
}
type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func LoadConfig() (Config, error) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("无法读取 YAML 文件: %v", err)
	}

	// 解析 YAML 文件到 Configuration 结构体
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("无法解析 YAML 文件: %v", err)
		return Config{}, err
	}
	return config, nil
}
