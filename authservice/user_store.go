package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
)

type UserStore interface {
	GetUserByID(ctx context.Context, userID string) (*types.User, error)
}

type SqlUserStore struct {
	db *sql.DB
}

func NewSqlUserStore(db *sql.DB) *SqlUserStore {
	return &SqlUserStore{
		db: db,
	}
}

func (s *SqlUserStore) GetUserByID(ctx context.Context, userID string) (*types.User, error) {
	query := `SELECT USERID , FIRSTNAME, LASTNAME, EMAIL, PHONE, PASSWORD WHERE USERID = ?`

	row := s.db.QueryRowContext(ctx, query, userID)
	user := &types.User{}

	err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("unexpected error occured")
	}

	return user, nil
}
