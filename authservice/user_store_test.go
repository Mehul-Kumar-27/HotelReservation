package main

import (
	"context"
	"errors"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserGetUserById_UserFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	expectedUserID := "123"

	expectedUser := &types.User{
		UserID:    expectedUserID,
		FirstName: "Firstname",
		LastName:  "lastname",
		Email:     "dummy@email",
		Phone:     "9999999999",
		Password:  "123456",
	}

	rows := sqlmock.NewRows([]string{"USERID", "FIRSTNAME", "LASTNAME", "EMAIL", "PHONE", "PASSWORD"})
	rows.AddRow(expectedUser.UserID, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.Phone, expectedUser.Password)

	mock.ExpectQuery(`SELECT USERID, FIRSTNAME, LASTNAME, EMAIL, PHONE, PASSWORD FROM USERS WHERE USERID = ?`).
		WithArgs(expectedUserID).
		WillReturnRows(rows)

	store := NewSqlUserStore(db)
	ctx := context.Background()
	user, err := store.GetUserByID(ctx, expectedUserID)

	assert.Nil(t, err)
	assert.Equal(t, expectedUser, user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there where unfulfilled expectations : %v", err)
	}
}

func TestGetUserById_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	expectedUserID := "some-random-id"

	insertRandomUser := &types.User{
		UserID:    "123456789",
		FirstName: "Firstname",
		LastName:  "lastname",
		Email:     "dummy@email",
		Phone:     "9999999999",
		Password:  "123456",
	}

	rows := sqlmock.NewRows([]string{"USERID", "FIRSTNAME", "LASTNAME", "EMAIL", "PHONE", "PASSWORD"})
	rows.AddRow(insertRandomUser.UserID, insertRandomUser.FirstName, insertRandomUser.LastName, insertRandomUser.Email, insertRandomUser.Phone, insertRandomUser.Password)

	mock.ExpectQuery(`SELECT USERID, FIRSTNAME, LASTNAME, EMAIL, PHONE, PASSWORD FROM USERS WHERE USERID = ?`).
    WithArgs(expectedUserID).
    WillReturnError(errors.New("sql.ErrNoRows"))
		

	store := NewSqlUserStore(db)
	ctx := context.Background()
	user, err := store.GetUserByID(ctx, expectedUserID)
	
	assert.Equal(t, errors.New("sql.ErrNoRows"), err)
	assert.Nil(t, user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there where unfulfilled expectations : %v", err)
	}
}
