package authen_and_post

import (
	"log"

	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient creates a gRPC client with built-in round-robin load balancing
func NewClient(serviceName string) (authen_and_post.AuthenticateAndPostClient, error) {
	// Use DNS resolver to discover multiple hosts
	conn, err := grpc.Dial(
		"dns:///" + serviceName, // "localhost:50051" if havent used Kubernetes Headless Service
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // Enable client-side LB
	)
	// conn, err := grpc.Dial(
	// 	"dns:///my-grpc-service.default.svc.cluster.local:50051",
	// 	grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )

	if err != nil {
		return nil, err
	}
	defer conn.Close()

	log.Println("gRPC client of authen_and_post service connected with round-robin load balancing")

	return authen_and_post.NewAuthenticateAndPostClient(conn), nil
}