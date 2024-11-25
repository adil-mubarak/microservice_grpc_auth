package main

import (
	"fmt"
	"log"
	"microservice_grpc_auth/database"
	"microservice_grpc_auth/pb/auth"
	"microservice_grpc_auth/user"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	db,err = database.ConnectDatabase()
	if err != nil{
		log.Fatalf("Failed to connect to the database: %v",err)
	}

	server := grpc.NewServer()

	authServiceServer := &user.AuthServiceServer{
		DB: db,
	}
	auth.RegisterAuthServiceServer(server, authServiceServer)

	reflection.Register(server)

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("Server is running on port :8080")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
