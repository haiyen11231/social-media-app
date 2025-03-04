package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/haiyen11231/social-media-app.git/configs"
	"github.com/haiyen11231/social-media-app.git/internal/app/authen_and_post_service"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"google.golang.org/grpc"
)

var path = flag.String("cfg", "/app/configs/files/test.yml", "path to config file of this service")

func main() {
	flag.Parse()

	cfg, err := configs.GetAuthenAndPostConfig(*path)
	if err != nil {
		log.Fatalf("Failed to get config: %s", err)
	}

	service, err := authen_and_post_service.NewAuthenAndPostService(cfg)
	if err != nil {
		log.Fatalf("Failed to init server: %s", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}

	aapServer := grpc.NewServer()
	authen_and_post.RegisterAuthenticateAndPostServer(aapServer, service)

	log.Printf("gRPC AAP Service server started on port %d", cfg.Port)
	if err := aapServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
	// // Graceful shutdown
	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer stop()

	// go func() {
	// 	log.Printf("gRPC AAP Service server started on port %d", cfg.Port)
	// 	if err := aapServer.Serve(lis); err != nil {
	// 		log.Fatalf("Failed to serve: %s", err)
	// 	}
	// }()

	// <-ctx.Done()
	// log.Println("Shutting down gRPC AAP Service server...")
	// aapServer.GracefulStop()
	// log.Println("gRPC AAP Service server stopped.")
}