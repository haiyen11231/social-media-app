package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/app/web_app/service"
)

func AddUserRouter(r *gin.RouterGroup, svc *service.WebService) {
	userRouter := r.Group("/users")
	
	userRouter.POST("/", svc.SignUp)
	userRouter.POST("/login", svc.LogIn)
	userRouter.POST("/authen", svc.AuthenticateUser)
	userRouter.POST("/refresh-token", svc.RefreshToken)
	userRouter.Use(svc.AuthenticateUser)
	userRouter.PUT("/", svc.EditUser)
}