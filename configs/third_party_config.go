package configs

type MySQLConfig struct {
	DSN                       string `yaml:"dsn"`
	DefaultStringSize         uint   `yaml:"defaultStringSize"`
	DisableDatetimePrecision  bool   `yaml:"disableDatetimePrecision"`
	DontSupportRenameIndex    bool   `yaml:"dontSupportRenameIndex"`
	SkipInitializeWithVersion bool   `yaml:"skipInitializeWithVersion"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type MinioConfig struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	BucketName string `yaml:"bucket_name"`
	UseSSL     bool   `yaml:"use_ssl"`
}

type KafkaConfig struct {
	Brokers          []string `yaml:"brokers"`
	Topic            string   `yaml:"topic"`
	ZookeeperConnect string   `yaml:"zookeeper_connect"`
}