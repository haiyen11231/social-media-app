package configs

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
)

type NewsfeedConfig struct {
	Port  int           `yaml:"port"`
	MySQL mysql.Config  `yaml:"mysql"`
	Redis redis.Options `yaml:"redis"`
	Minio MinioConfig   `yaml:"minio"`
	Kafka KafkaConfig   `yaml:"kafka"`
}