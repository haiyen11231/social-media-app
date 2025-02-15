package authen_and_post

import (
	"log"

	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient creates a gRPC client and connects to NGINX server
func NewClient(nginxHost string) (authen_and_post.AuthenticateAndPostClient, error) {
	conn, err := grpc.Dial(nginxHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect to NGINX server at %s: %v", nginxHost, err)
		return nil, err
	}
	defer conn.Close()

	log.Println("gRPC client of authen_and_post service is connected to NGINX server...")
	return authen_and_post.NewAuthenticateAndPostClient(conn), nil
}


// Qs: When do I need to use pointer? When to use defer?