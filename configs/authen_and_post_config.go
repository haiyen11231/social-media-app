package configs

import (
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
)

type AuthenAndPostConfig struct {
	Port  int           `yaml:"port"`
	MySQL mysql.Config  `yaml:"mysql"`
	Redis redis.Options `yaml:"redis"`
	Minio MinioConfig   `yaml:"minio"`
}