package api

import (
	"context"
	"log"

	"github.com/Mehul-Kumar-27/HotelReservation/proto/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoginPayload struct {
	Userid   string `json:"userid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	status      int
	response    string
	accesstoken string
}

func NewLoginPayload(userid, email, password string) *LoginPayload {
	return &LoginPayload{
		Userid:   userid,
		Email:    email,
		Password: password,
	}
}

type JWTPayload struct {
	Token string
}

func NewJWTPayload(token string) *JWTPayload {
	return &JWTPayload{
		Token: token,
	}
}

type Auth struct {
}

func NewAuth() *Auth {

	return &Auth{}
}

type AuthRequest interface {
	LoginService(ctx context.Context, loginPayload LoginPayload) AuthResponse
	JwtAuthService(ctx context.Context, jwtPayload JWTPayload) AuthResponse
}

func (a *Auth) LoginService(ctx context.Context, loginPayload LoginPayload) AuthResponse {
	log.Println("Sending the login request")
	conn, err := grpc.Dial("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return AuthResponse{
			status:      500,
			response:    "unexpected error occured",
			accesstoken: "",
		}
	}

	defer conn.Close()

	client := auth.NewAuthServiceClient(conn)
	response, err := client.LoginService(ctx, &auth.Login{Userid: loginPayload.Userid, Email: loginPayload.Email, Password: loginPayload.Password})
	if err != nil {
		return AuthResponse{
			status:      400,
			response:    "unauthorized",
			accesstoken: "",
		}
	}

	log.Println(response.GetAcesstoken())

	return AuthResponse{
		status:      int(response.Response.GetStatus()),
		response:    response.Response.GetBody(),
		accesstoken: response.GetAcesstoken(),
	}

}

func (a *Auth) JwtAuthService(ctx context.Context, jwtPayload JWTPayload) AuthResponse {
	log.Println("Sending the login request")
	conn, err := grpc.Dial("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return AuthResponse{
			status:      500,
			response:    "unexpected error occured",
			accesstoken: "",
		}
	}
	defer conn.Close()
	client := auth.NewAuthServiceClient(conn)

	response, err := client.JwtAuthService(ctx, &auth.JwToken{Token: jwtPayload.Token})
	if err != nil {
		return AuthResponse{
			status:      400,
			response:    "unauthorized user",
			accesstoken: "",
		}
	}

	return AuthResponse{
		status:      int(response.Response.GetStatus()),
		response:    response.Response.GetBody(),
		accesstoken: response.GetUserid(),
	}
}
