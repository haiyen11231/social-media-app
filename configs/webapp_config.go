package configs

type WebappConfig struct {
	Port          int          `yaml:"port"`
	NginxHost     string       `yaml:"nginx_host"`
	AuthenAndPost ServiceHosts `yaml:"authen_and_post"`
	Newsfeed      ServiceHosts `yaml:"newsfeed"`
}

type ServiceHosts struct {
	Hosts []string `yaml:"hosts"`
}
