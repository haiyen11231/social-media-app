package configs

type NewsfeedConfig struct {
	Port  int         `yaml:"port"`
	MySQL MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
	Minio MinioConfig `yaml:"minio"`
	Kafka KafkaConfig `yaml:"kafka"`
}