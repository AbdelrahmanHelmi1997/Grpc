# GrpcServer
- Grpc server used to create user , login user  & get user info

#**Server Folder**
--------------------

-**GrpcFunctions.go**

- The server have three functions :

    - "CreateUser" ---> Which takes the request and create a user and save him in the DataBase and returns Access token for the user
  
    - "Login" ---> takes user's credentials as a request and search for the user in the DataBase by user name and when find the user returns Access token
  
    - "GetUserInfo" -----> Takes Access token as MetaData and extract user's Id from the token and search for user in DataBase and returns user's saved data

-**GrpcConnection.go**

-This file have one funtion "GrpcServerConnection" which helps me to connect with the local server I created 
  
  --------------------------------------------------------------------------------------------------------------------------------------------------------
