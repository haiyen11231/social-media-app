package configs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type appConfigs struct {
	MySQLConfig          MySQLConfig           `yaml:"mysql"`
	RedisConfig          RedisConfig           `yaml:"redis"`
	MinioConfig          MinioConfig           `yaml:"minio"`
	KafkaConfig          *KafkaConfig          `yaml:"kafka"`
	AuthenAndPostConfig  *AuthenAndPostConfig  `yaml:"authen_and_post_config"`
	NewsfeedConfig       *NewsfeedConfig       `yaml:"newsfeed_config"`
	WebappConfig         *WebappConfig         `yaml:"webapp_config"`
}

func loadAppConfigs(path string) (*appConfigs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file (path: %s): %w", path, err)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file (path: %s): %w", path, err)
	}

	config := &appConfigs{}
	if err := yaml.Unmarshal(fileContent, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file (path: %s): %w", path, err)
	}

	log.Printf("Loaded config: %+v", config)
	return config, nil
}

func GetAuthenAndPostConfig(path string) (*AuthenAndPostConfig, error) {
	config, err := loadAppConfigs(path)
	if err != nil {
		return nil, err
	}
	if config.AuthenAndPostConfig == nil {
		return nil, fmt.Errorf("authen_and_post_config is missing in config file")
	}
	return config.AuthenAndPostConfig, nil
}

func GetNewsfeedConfig(path string) (*NewsfeedConfig, error) {
	config, err := loadAppConfigs(path)
	if err != nil {
		return nil, err
	}
	if config.NewsfeedConfig == nil {
		return nil, fmt.Errorf("newsfeed_config is missing in config file")
	}
	return config.NewsfeedConfig, nil
}

func GetWebappConfig(path string) (*WebappConfig, error) {
	config, err := loadAppConfigs(path)
	if err != nil {
		return nil, err
	}
	if config.WebappConfig == nil {
		return nil, fmt.Errorf("webapp_config is missing in config file")
	}
	return config.WebappConfig, nil
}