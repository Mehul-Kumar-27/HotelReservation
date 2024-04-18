package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Mehul-Kumar-27/HotelReservation/proto/gen/auth"
	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("123456789")

type LoginPayload struct {
	UserId   string `json:"userid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLoginPayload(userid, email, password string) *LoginPayload {
	return &LoginPayload{
		UserId:   userid,
		Email:    email,
		Password: password,
	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	Userid    string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	EXP       int64
}

func NewUserClaims(userid, firstname, lastname, email, phone string, exp int64) *UserClaims {
	return &UserClaims{
		Userid:    userid,
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		Phone:     phone,
		EXP:       exp,
	}
}

type AuthInterface interface {
	LoginService(ctx context.Context, req *auth.Login) (*auth.LoginResponse, error)
	JwtAuthService(ctx context.Context, req *auth.JwToken) (*auth.JwTokenResponse, error)
}

type Auth struct {
	auth.UnimplementedAuthServiceServer
	db        *sql.DB
	userStore UserStore
}

func NewAuth(db *sql.DB) *Auth {
	return &Auth{
		db:        db,
		userStore: NewSqlUserStore(db),
	}
}

func (a *Auth) HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (a *Auth) VerifyPassword(password, hashedpassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
	return err == nil
}

func (a *Auth) CreateToken(user *types.User) (string, error) {

	claims := jwt.MapClaims{
		"registered": jwt.RegisteredClaims{},
		"userid":     user.UserID,
		"email":      user.Email,
		"firsname":   user.FirstName,
		"lastname":   user.LastName,
		"phone":      user.Phone,
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (a *Auth) LoginService(ctx context.Context, req *auth.Login) (*auth.LoginResponse, error) {
	log.Println("Recieved the request for login")
	loginPayload := NewLoginPayload(req.GetUserid(), req.GetEmail(), req.GetPassword())
	log.Println(loginPayload.UserId)
	log.Println(loginPayload.Password)
	user, err := a.userStore.GetUserByID(ctx, loginPayload.UserId)
	if err != nil {
		return &auth.LoginResponse{
			Response: &auth.Response{
				Status: 400,
				Body:   "user not found",
			},
			Acesstoken: "",
		}, nil
	}

	verified := a.VerifyPassword(loginPayload.Password, user.Password)
	if verified {
		/////// Generate the jwt token and return it
		token, err := a.CreateToken(user)
		if err != nil {
			return &auth.LoginResponse{
				Response: &auth.Response{
					Status: 500,
					Body:   "unexpected error occured",
				},
				Acesstoken: "",
			}, nil
		}
		log.Println(token)
		return &auth.LoginResponse{
			Response: &auth.Response{
				Status: 200,
				Body:   "authenticated",
			},
			Acesstoken: token,
		}, nil

	}

	return &auth.LoginResponse{
		Response: &auth.Response{
			Status: 400,
			Body:   "unauthorized user",
		},
		Acesstoken: "",
	}, nil

}

func (a *Auth) JwtAuthService(ctx context.Context, req *auth.JwToken) (*auth.JwTokenResponse, error) {
	tokenRequest := req.GetToken()
	log.Println("Request for the jwt auth")
	log.Println(tokenRequest)
	token, err := jwt.ParseWithClaims(
		tokenRequest,
		&UserClaims{},
		func(tokenRequest *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		log.Printf("Inside the err %v", err)

		return &auth.JwTokenResponse{
				Response: &auth.Response{
					Status: 400,
					Body:   "token expired",
				},

				Userid: "",
			},
			nil
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		log.Printf("Inside the not ok err %v", err)

		return &auth.JwTokenResponse{
				Response: &auth.Response{
					Status: 400,
					Body:   "unauthorized user",
				},

				Userid: "",
			},
			nil
	}

	return &auth.JwTokenResponse{
			Response: &auth.Response{
				Status: 200,
				Body:   "jwt auth success",
			},

			Userid: claims.Userid,
		},
		nil

}
