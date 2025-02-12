package configs

type WebappConfig struct {
	Port          int `yaml:"port"`
	AuthenAndPost struct {
		Hosts []string `yaml:"hosts"`
	} `yaml:"authen_and_post"`
	Newsfeed struct {
		Hosts []string `yaml:"hosts"`
	} `yaml:"newsfeed"`
}