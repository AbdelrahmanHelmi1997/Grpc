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
	"google.golang.org/grpc/metadata"
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
func (s *server) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	var user model.User
	var foundUser model.User
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

	CreateUser := model.User{
		ID:        primitive.NewObjectID(),
		FirstName: user.FirstName,
		Username:  user.Username,
		Password:  user.Password,
	}

	dataBase.UsersDB.InsertOne(ctx, CreateUser)

	dataBase.UsersDB.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)

	token, _ := Helper.GenerateAllTokens(foundUser.ID, foundUser.Username, foundUser.FirstName)

	return &proto.CreateUserResponse{Message: "Created", Token: token}, nil

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

	token, _ := Helper.GenerateAllTokens(foundUser.ID, foundUser.Username, foundUser.FirstName)

	return &proto.LoginResponse{Message: "successs", Token: token}, nil

}

func (s *server) GetUserInfo(ctx context.Context, request *proto.UserInfoRequest) (*proto.UserInfoResponse, error) {
	var values []string
	var AccessToken string
	var user model.User

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values = md.Get("auth")
	}

	if len(values) > 0 {
		AccessToken = values[0]
	}

	claims, _ := Helper.ValidateToken(AccessToken)

	err := dataBase.UsersDB.FindOne(ctx, bson.M{"_id": claims.ID}).Decode(&user)

	if err != nil {
		return &proto.UserInfoResponse{Message: "User Not Found"}, nil
	}
	return &proto.UserInfoResponse{Message: "Success", Id: user.ID.Hex(), Name: user.FirstName, Username: user.Username}, nil

}
