package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/haiyen11231/social-media-app.git/configs"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
)

var path = flag.String("cfg", "test.yml", "path to config file of this service")

func main() {
	flag.Parse()
	cfg, err := configs.GetWebappConfig(*path)
	fmt.Println(cfg)
	if err != nil {
		log.Fatalf("Failed to get config: %s", err)
	}

	// Add NGINX Host!!!
	webService, err := service.NewWebService(cfg.NginxHost)
	if err != nil {
		log.Fatalf("Failed to init server: %s", err)
	}

	webController := &web_app.WebController{
		WebService: *webService,
		Port:       cfg.Port,
	}
	webController.Run()
}