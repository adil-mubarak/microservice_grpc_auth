package user

import (
	"context"
	"fmt"
	"microservice_grpc_auth/models"
	"microservice_grpc_auth/pb/auth"
	"microservice_grpc_auth/tokenjwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

type AuthServiceServer struct {
	auth.UnimplementedAuthServiceServer
	DB *gorm.DB
}

func (S *AuthServiceServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.AuthResponse, error) {

	var userCount int64
	err := S.DB.Model(&models.User{}).Where("user_name = ?",req.Username).Count(&userCount).Error
	if err != nil{
		return &auth.AuthResponse{
			Message: "Failed to check user existence",
			Success: false,
		},err
	}
	if userCount > 0{
		return &auth.AuthResponse{
			Message: "Username already exists",
			Success: false,
		},nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &auth.AuthResponse{
			Message: "Falied to hash password",
			Success: false,
		}, err
	}

	user := &models.User{
		UserName: req.Username,
		Password: string(hashedPassword),
	}

	if err := db.Create(user).Error; err != nil {
		return &auth.AuthResponse{
			Message: "Falied to register user",
			Success: false,
		}, err
	}

	return &auth.AuthResponse{
		Message: "User registered successfully",
		Success: true,
	}, nil

}

func (S *AuthServiceServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.AuthResponse, error) {

	user, err := findUserbyUsername(req.Username)

	if err != nil {
		return &auth.AuthResponse{
			Message: "User not found",
			Success: false,
		}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &auth.AuthResponse{
			Message: "Invalid credentials",
			Success: false,
		}, err
	}

	token, err := tokenjwt.GenerateJWT(user.ID, user.UserName, user.Password)
	if err != nil {
		return &auth.AuthResponse{
			Message: "Failed to generate token",
			Success: false,
		}, err
	}

	return &auth.AuthResponse{
		Token:   token,
		Message: "Login succesful",
		Success: true,
	}, nil
}

// func (s *AuthServiceServer) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest)(*auth.AuthResponse,error){

// }

func findUserbyUsername(username string) (*models.User, error) {
	var user models.User

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return &user, nil
}
