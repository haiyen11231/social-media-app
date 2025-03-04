package service

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/haiyen11231/social-media-app.git/internal/grpc/pb/authen_and_post"
	"github.com/haiyen11231/social-media-app.git/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (svc *WebService) SignUp(ctx *gin.Context) {
	var jsonRequest models.SignUpRequest
	if err := ctx.ShouldBindJSON(&jsonRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: err.Error()})
		return
	}

	log.Println("Sending request to Sign Up")
	response, err := svc.authenAndPostClient.SignUp(ctx, &authen_and_post.SignUpRequest{
		FirstName: jsonRequest.FirstName,
		LastName:  jsonRequest.LastName,
		Dob:       timestamppb.New(jsonRequest.DoB),
		Email:     jsonRequest.Email,
		Username:  jsonRequest.Username,
		Password:  jsonRequest.Password,
	})

	log.Println("Received response from Sign Up")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}


	ctx.JSON(http.StatusOK, response)
}

func (svc *WebService) LogIn(ctx *gin.Context) {
	var jsonRequest models.LogInRequest
	if err := ctx.ShouldBindJSON(&jsonRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: err.Error()})
		return
	}

	response, err := svc.authenAndPostClient.LogIn(ctx, &authen_and_post.LogInRequest{
		Username: jsonRequest.Username,
		Password: jsonRequest.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("refresh_token", response.RefreshToken, 60*60*24, "/", "", true, true)
	ctx.JSON(http.StatusOK, response)
}

func (svc *WebService) EditUser(ctx *gin.Context) {
	userId := ctx.GetUint64("user_id")

	var jsonRequest models.EditUserRequest
	if err := ctx.ShouldBindJSON(&jsonRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.MessageResponse{Message: err.Error()})
		return
	}

	var firstName *string
	if jsonRequest.FirstName != "" {
		firstName = &jsonRequest.FirstName
	}

	var lastName *string
	if jsonRequest.LastName != "" {
		lastName = &jsonRequest.LastName
	}

	var dob *timestamppb.Timestamp
	if !jsonRequest.DoB.IsZero() {
		dob = timestamppb.New(jsonRequest.DoB)
	}

	var password *string
	if jsonRequest.Password != "" {
		password = &jsonRequest.Password
	}

	response, err := svc.authenAndPostClient.EditUser(ctx, &authen_and_post.EditUserRequest{
		UserId:    userId,
		FirstName: firstName,
		LastName:  lastName,
		Dob:       dob,
		Password:  password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (svc *WebService) AuthenticateUser(ctx *gin.Context) {
	token := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, &models.MessageResponse{Message: "Unauthorized"})
		return
	}

	response, err := svc.authenAndPostClient.AuthenticateUser(ctx, &authen_and_post.AuthenticateUserRequest{
		Token: token,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	if !response.IsValid {
		ctx.JSON(http.StatusUnauthorized, &models.MessageResponse{Message: "Unauthorized"})
		return
	}

	ctx.Set("user_id", response.UserId)
	ctx.Next()
}

func (svc *WebService) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &models.MessageResponse{Message: "Unauthorized"})
		return
	}

	response, err := svc.authenAndPostClient.RefreshToken(ctx, &authen_and_post.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}