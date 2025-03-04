package service

import (
	"fmt"
	"log"

	"github.com/haiyen11231/social-media-app.git/configs"
	aapClient "github.com/haiyen11231/social-media-app.git/internal/client/authen_and_post"
	nfsClient "github.com/haiyen11231/social-media-app.git/internal/client/newsfeed"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/newsfeed"
)

type WebService struct {
	authenAndPostClient authen_and_post.AuthenticateAndPostClient
	newsfeedClient      newsfeed.NewsfeedClient
}

func NewWebService(cfg *configs.WebappConfig) (*WebService, error) {
	authenAndPostClient, err := aapClient.NewClient(cfg.AuthenAndPost.Hosts[0])
	if err != nil {
		log.Printf("Error creating AuthenAndPost client: %v\n", err)
		return nil, fmt.Errorf("failed to create AuthenAndPost client: %w", err)
	}

	newsfeedClient, err := nfsClient.NewClient(cfg.Newsfeed.Hosts[0])
	if err != nil {
		log.Printf("Error creating Newsfeed client: %v\n", err)
		return nil, fmt.Errorf("failed to create Newsfeed client: %w", err)
	}

	return &WebService{
		authenAndPostClient: authenAndPostClient,
		newsfeedClient:      newsfeedClient,
	}, nil
}