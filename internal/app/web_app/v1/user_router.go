package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
)

func AddUserRouter(r *gin.RouterGroup, svc *service.WebService) {
	userRouter := r.Group("/users")
	
	userRouter.POST("/", )
	userRouter.POST("login", )
	userRouter.PUT("/", )
	userRouter.POST("authen", )
	userRouter.POST("refresh-token", )
}