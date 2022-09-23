package server

import (
	"context"
	"test/Helper"
	"test/dataBase"
	"test/model"
	"test/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	var user model.User
	var foundUser model.User
	user.FirstName, user.Username, user.Password = request.GetName(), request.GetUsername(), request.GetPassword()

	count, err := dataBase.UsersDB.CountDocuments(ctx, bson.M{"username": user.Username})
	if err != nil {
		return nil, status.Error(409, "Error while Checking user")
	}

	if count > 0 {
		return nil, status.Error(409, "User already exists")
	}

	password := Helper.HashPassword(user.Password)
	user.Password = password

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

func (s *Server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	var user model.User
	var foundUser model.User
	user.Username, user.Password = request.GetUsername(), request.GetPassword()

	err := dataBase.UsersDB.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)

	if err != nil {
		return nil, status.Error(400, "Username is incorrect")
	}

	passwordIsValid := Helper.VerifyPassword(user.Password, foundUser.Password)

	if passwordIsValid != true {

		return nil, status.Error(400, "Password is incorrect")
	}

	token, _ := Helper.GenerateAllTokens(foundUser.ID, foundUser.Username, foundUser.FirstName)
	return &proto.LoginResponse{Message: "successs", Token: token}, nil

}

func (s *Server) GetUserInfo(ctx context.Context, request *proto.UserInfoRequest) (*proto.UserInfoResponse, error) {
	var metaDataArray []string
	var AccessToken string
	var user model.User

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		metaDataArray = md.Get("auth")
	}

	if len(metaDataArray) > 0 {
		AccessToken = metaDataArray[0]
	}

	claims, msg := Helper.ValidateToken(AccessToken)

	if msg != "" {
		return nil, status.Error(403, msg)
	}
	err := dataBase.UsersDB.FindOne(ctx, bson.M{"_id": claims.ID}).Decode(&user)

	if err != nil {
		return nil, status.Error(400, "User not found")
	}

	return &proto.UserInfoResponse{Message: "Success", Id: user.ID.Hex(), Name: user.FirstName, Username: user.Username}, nil

}
