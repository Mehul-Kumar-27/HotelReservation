package main

import (
	"context"
	"database/sql"
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

type AuthInterface interface {
	LoginService(ctx context.Context, req *auth.Login) (*auth.LoginResponse, error)
	JwtAuthService(ctx context.Context, req *auth.JwToken) (*auth.JwTokenResponse, error)
}

type Auth struct {
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
		"userid":   user.UserID,
		"email":    user.Email,
		"firsname": user.FirstName,
		"lastname": user.LastName,
		"phone":    user.Phone,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(secretKey)
}

func (a *Auth) LoginService(ctx context.Context, req *auth.Login) (*auth.LoginResponse, error) {
	loginPayload := NewLoginPayload(req.Userid, req.Email, req.Password)

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
					Status: 400,
					Body:   "unexpected error occured",
				},
				Acesstoken: "",
			}, nil
		}
		return &auth.LoginResponse{
			Response: &auth.Response{
				Status: 200,
				Body:   "unauthorized user",
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
