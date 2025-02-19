package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/haiyen11231/social-media-app.git/configs"
	"github.com/haiyen11231/social-media-app.git/internal/app/newsfeed_service"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/newsfeed"
	"google.golang.org/grpc"
)

var path = flag.String("cfg", "config.yml", "path to config file of this service")

func main() {
	cfg, err := configs.GetNewsfeedConfig(*path)
	if err != nil {
		log.Fatalf("Failed to get config: %s", err)
	}

	service, err := newsfeed_service.NewNewsfeedService(cfg)
	if err != nil {
		log.Fatalf("Failed to init server: %s", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}

	nsfServer := grpc.NewServer()
	newsfeed.RegisterNewsfeedServer(nsfServer, service)

	log.Printf("gRPC NSF Service server started on port %s", cfg.Port)
	if err := nsfServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}

}