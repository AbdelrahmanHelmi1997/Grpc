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

#**dataBase Folder**
----------------------

**database.go**

-In this file I create DB() funtion which's helps me to connect with mongodb

-----------------------------------------------------------------------------------------------------------------------------------------------------------

#**Helper Folder**
------------------------

-**Helper.go**

  - In this file I created to funtions : 
 
     -HashPassword --> which hash the password the came with the user's request "When creating"
     
     -Verify Password --> used to compare between the hashed password in the database and the password came with the login request and validate that the password is correct 
     
-#**HelperToken.go**

   -GenerateAllToken--> Used to generate a token for the useer
   
   -ValidateToken--> Used to validate the token formate and check if the token is expired or not









