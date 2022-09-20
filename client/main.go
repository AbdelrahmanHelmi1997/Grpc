package main

import (
	"fmt"
	"log"
	"net/http"
	"test/model"
	"test/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewUserClient(conn)
	g := gin.Default()

	g.POST("/CreateUser", func(ctx *gin.Context) {
		var user model.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, "Bad Request")
			return
		}
		req := &proto.CreateUserRequest{Name: user.FirstName, Username: user.Username, Password: user.Password}
		if response, err := client.CreateUSer(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprint(response.Message),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

	})

	g.POST("/Login", func(ctx *gin.Context) {
		var user model.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, "Bad Request")
			return
		}
		req := &proto.LoginRequest{Username: user.Username, Password: user.Password}
		if response, err := client.Login(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"Message": fmt.Sprint(response.Message),
				"Token":   fmt.Sprint(response.Token),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

	})

	g.GET("GetUser", func(ctx *gin.Context) {

		userToken := ctx.Request.Header.Get("token")

		if userToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			ctx.Abort()
			return
		}
		req := &proto.UserInfoRequest{Token: userToken}

		if response, err := client.GetUserInfo(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"Message":  fmt.Sprint(response.Message),
				"Name":     fmt.Sprint(response.Name),
				"UserName": fmt.Sprint(response.Usernasme),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

	})
	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
