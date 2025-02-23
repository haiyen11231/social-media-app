package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
)

func AddNewsfeedRouter(r *gin.RouterGroup, svc *service.WebService) {
	newsfeedRouter := r.Group("/newsfeed")
	
	newsfeedRouter.Use(svc.AuthenticateUser)
	newsfeedRouter.GET("/", )
}