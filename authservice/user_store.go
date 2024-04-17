package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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
	query := `SELECT USERID, FIRSTNAME, LASTNAME, EMAIL, PHONE, PASSWORD FROM USERS WHERE USERID = ?`
	row := s.db.QueryRowContext(ctx, query, userID)

	user := &types.User{}
	err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		log.Printf("Error processing the user get request: %v", err)
		return nil, fmt.Errorf("unexpected error occurred")
	}

	return user, nil
}
