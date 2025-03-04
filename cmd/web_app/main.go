package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/haiyen11231/social-media-app.git/configs"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
)

var path = flag.String("cfg", "/app/configs/files/test.yml", "path to config file of this service")

func main() {
	flag.Parse()

	cfg, err := configs.GetWebappConfig(*path)
	fmt.Println(cfg)
	if err != nil {
		log.Fatalf("Failed to get config: %s", err)
	}

	webService, err := service.NewWebService(cfg)
	if err != nil {
		log.Fatalf("Failed to init server: %s", err)
	}

	webController := &web_app.WebController{
		WebService: *webService,
		Port:       cfg.Port,
	}

	log.Printf("Web app started on port %d", cfg.Port)
	webController.Run()
	// // Graceful shutdown
	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer stop()

	// go func() {
	// 	log.Printf("Web app started on port %d", cfg.Port)
	// 	webController.Run()
	// }()

	// <-ctx.Done()
	// log.Println("Shutting down web app...")
	// // Add any cleanup logic here
	// log.Println("Web app stopped.")
}