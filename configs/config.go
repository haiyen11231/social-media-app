package configs

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
)

type appConfigs struct {
	MySQLConfig mysql.Config `yaml:"mysql"`
	RedisConfig redis.Options `yaml:"redis"`
	MinioConfig *MinioConfig `yaml:"minio"`
	KafkaConfig *KafkaConfig `yaml:"kafka"`
	AuthenAndPostConfig *AuthenAndPostConfig `yaml:"authen_and_post_config"`
	NewsfeedConfig *NewsfeedConfig `yaml:"newsfeed_config"`
	WebappConfig *WebappConfig `yaml:"webapp_config"`
}

func loadAppConfigs(path string) (*appConfigs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open config file (path: %s): %s", path, err)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file (path: %s): %s", path, err)
	}
	
	config := &appConfigs{}
	if err := yaml.Unmarshal(fileContent, config); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal config file (path: %s): %s", path, err)
	}

	return config, nil
} 

func GetAuthenAndPostConfig(path string) (*AuthenAndPostConfig, error) {
	config, err := loadAppConfigs(path)
	if err != nil {
		return nil, err
	}

	return config.AuthenAndPostConfig, nil
}

func GetNewsfeedConfig(path string) (*NewsfeedConfig, error) {
	config, err := loadAppConfigs(path)
	if err != nil {
		return nil, err
	}

	return config.NewsfeedConfig, nil
}

func GetWebappConfig(path string) (*WebappConfig, error) {
	config, err := loadAppConfigs(path)
	if err != nil {
		return nil, err
	}

	return config.WebappConfig, nil
}