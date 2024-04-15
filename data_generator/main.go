package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config holds database connection details
type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

// NewConfig returns a new Config
func NewConfig() *Config {
	return &Config{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     3306,
		Database: os.Getenv("MYSQL_DATABASE"),
	}
}

type ConnectToSQL interface {
	Open(driver, dsn string) (*sql.DB, error)
	PingContext(ctx context.Context, db *sql.DB) error
}

type MySql struct {
}

func NewMySqlStruct() *MySql {
	return &MySql{}
}

func (m *MySql) Open(driver, dsn string) (*sql.DB, error) {
	
	return sql.Open(driver, dsn)
}

func (m *MySql) PingContext(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}

func main() {
	// Connect to MySQL database
	log.Println("This is a CLI application")
	config := NewConfig()
	if config.User == "" || config.Password == "" || config.Host == "" || config.Database == "" {
		log.Fatal("Missing required environment variables")
	}

	sqlStruct := NewMySqlStruct()
	// Connect to MySQL database with retries and context cancellation
	db, err := connectToMySQL(context.Background(), config, sqlStruct)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(500)

	dgen := NewCli(db)

	// Handle graceful shutdown on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run the CLI
	go func() {
		err := dgen.ExecuteContext(ctx)
		if err != nil {
			log.Printf("Error running the CLI: %v", err)
		}
	}()
	<-ctx.Done()
	log.Println("Shutting down the data generator")

	// Close the database connection
	if err := db.Close(); err != nil {
		log.Println("Error closing the database:", err)
	}

	log.Println("Exiting application")
	os.Exit(0)

}

func connectToMySQL(ctx context.Context, config *Config, database ConnectToSQL) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User, config.Password, config.Host, config.Port, config.Database)

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = database.Open("mysql", dsn)
		if err != nil {
			fmt.Printf("Error connecting to the database, %v", err)
			fmt.Println("Trying againg in 2 second")
			time.Sleep(5 * time.Second)
			continue
		}

		err = database.PingContext(ctx, db)
		if err == nil {
			fmt.Println("Connection successful!")
			return db, nil
		}

		fmt.Printf("Attempt %d: Error pinging database: %v\n", i+1, err)
		db.Close()
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("error connecting to the sql database, %v", err)
}
