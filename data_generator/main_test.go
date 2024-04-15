package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSql struct {
	mock.Mock
}

func (m *MockSql) Open(driver, dsn string) (*sql.DB, error) {

	args := m.Called(driver, dsn)
	
	return args.Get(0).(*sql.DB), nil
}

func (m *MockSql) PingContext(ctx context.Context, db *sql.DB) error {
	args := m.Called(ctx, db)
	return args.Error(0)
}

func TestConnetToSql(t *testing.T) {
	config := &Config{
		User:     "username",
		Password: "password",
		Host:     "localhost",
		Port:     3306,
		Database: "testdb",
	}

	mockSqlConnect := new(MockSql)

	mockSqlConnect.On("Open", "mysql", "username:password@tcp(localhost:3306)/testdb").Return(&sql.DB{}, nil)

	ctx := context.Background()

	mockSqlConnect.On("PingContext", ctx, &sql.DB{}).Return(nil)

	// Call the function being tested with the mock object
	db, err := connectToMySQL(ctx, config, mockSqlConnect)

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that the database connection object is not nil
	assert.NotNil(t, db)

	// Assert that all expectations were met
	mockSqlConnect.AssertExpectations(t)
}
