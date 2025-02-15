package newsfeed

import (
	"log"

	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/newsfeed"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CreateNewsfeedClient creates a gRPC client and connects to NGINX server
func CreateNewsfeedClient (nginxAddr string) (newsfeed.NewsfeedClient, error) {
	conn, err := grpc.Dial(nginxAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect to NGINX server at %s: %v", nginxAddr, err)
		return nil, err
	}
	defer conn.Close()

	return newsfeed.NewNewsfeedClient(conn), nil
}