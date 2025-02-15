package service

import (
	aapClient "github.com/haiyen11231/social-media-app.git/internal/client/authen_and_post"
	nfClient "github.com/haiyen11231/social-media-app.git/internal/client/newsfeed"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/newsfeed"
)

type WebService struct {
	authenAndPostClient authen_and_post.AuthenticateAndPostClient
	newsfeedClient newsfeed.NewsfeedClient
}

func NewWebService(nginxHost string) (*WebService, error) {
	authenAndPostClient, err := aapClient.NewClient(nginxHost)
	if err != nil {
		return nil, err
	}

	newsfeedClient, err := nfClient.NewClient(nginxHost)
	if err != nil {
		return nil, err
	}

	return &WebService{
		authenAndPostClient: authenAndPostClient,
		newsfeedClient: newsfeedClient,
	}, nil
}