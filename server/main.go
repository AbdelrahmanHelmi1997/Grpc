package main

import (
	"context"
	"net"
	"test/Helper"
	"test/dataBase"
	"test/model"
	"test/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedUserServer
}

func main() {

	dataBase.DB()
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterUserServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}
func (s *server) CreateUSer(ctx context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	var user model.User
	user.FirstName, user.Username, user.Password = request.GetName(), request.GetUsername(), request.GetPassword()

	count, err := dataBase.UsersDB.CountDocuments(ctx, bson.M{"username": user.Username})
	if err != nil {
		return &proto.CreateUserResponse{Message: "Error while Checking user"}, nil
	}

	password := Helper.HashPassword(user.Password)
	user.Password = password

	if count > 0 {
		return &proto.CreateUserResponse{Message: "User is already exsits"}, nil
	}

	token, _ := Helper.GenerateAllTokens(user.Username, user.FirstName)
	user.Token = token

	CreateUser := model.User{
		ID:        primitive.NewObjectID(),
		FirstName: user.FirstName,
		Username:  user.Username,
		Password:  user.Password,
		Token:     user.Token,
	}

	dataBase.UsersDB.InsertOne(ctx, CreateUser)

	return &proto.CreateUserResponse{Message: "Created"}, nil

}

func (s *server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	var user model.User
	var foundUser model.User
	user.Username, user.Password = request.GetUsername(), request.GetPassword()

	err := dataBase.UsersDB.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
	if err != nil {
		return &proto.LoginResponse{Message: "Email is invalid"}, nil
	}

	passwordIsValid, _ := Helper.VerifyPassword(user.Password, foundUser.Password)

	if passwordIsValid != true {

		return &proto.LoginResponse{Message: "Password is Wrong"}, nil
	}

	token, _ := Helper.GenerateAllTokens(foundUser.Username, foundUser.FirstName)

	Helper.UpdateAllTokens(token, foundUser.ID)
	err = dataBase.UsersDB.FindOne(ctx, bson.M{"_id": foundUser.ID}).Decode(&foundUser)

	return &proto.LoginResponse{Message: "successs", Token: foundUser.Token}, nil

}

func (s *server) GetUserInfo(ctx context.Context, request *proto.UserInfoRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	user.Token = request.GetToken()
	err := dataBase.UsersDB.FindOne(ctx, bson.M{"token": user.Token}).Decode(&user)

	if err != nil {
		return &proto.UserInfoResponse{Message: "User Not Found"}, nil
	}

	return &proto.UserInfoResponse{Message: "Success", Name: user.FirstName, Usernasme: user.Username}, nil

}
