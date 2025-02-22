package configs

type AuthenAndPostConfig struct {
	Port  int         `yaml:"port"`
	MySQL MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
	Minio MinioConfig `yaml:"minio"`
}