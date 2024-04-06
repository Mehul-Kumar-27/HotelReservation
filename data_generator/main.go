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

// func NewConfig() *Config {
// 	return &Config{
// 		User:     "mehul",
// 		Password: "mehulpassword",
// 		Host:     "localhost",
// 		Port:     3306,
// 		Database: "reservation",
// 	}
// }

func main() {
	// Connect to MySQL database
	log.Println("This is a CLI application")
	config := NewConfig()
	if config.User == "" || config.Password == "" || config.Host == "" || config.Database == "" {
		log.Fatal("Missing required environment variables")
	}
	// Connect to MySQL database with retries and context cancellation
	db, err := connectToMySQL(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

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

func connectToMySQL(ctx context.Context, config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User, config.Password, config.Host, config.Port, config.Database)

	var count int
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("connection context cancelled")
		default:
			db, err := sql.Open("mysql", dsn)
			if err != nil {
				if count < 10 {
					log.Println("Error connecting to database. Retrying in 3 seconds...")
					log.Printf("Error encountered is: %v", err)
					time.Sleep(3 * time.Second)
					count++
				} else {
					return nil, fmt.Errorf("failed to connect to database after %d retries: %w", count, err)
				}
			} else {
				// Try to ping the database with context cancellation
				pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer pingCancel()

				err := db.PingContext(pingCtx)
				if err != nil {
					if count < 10 {
						log.Println("Error pinging the database. Retrying in 3 seconds...")
						db.Close() // Close the connection before retrying
						time.Sleep(3 * time.Second)
						count++
					} else {
						db.Close() // Close the connection on final failure
						return nil, fmt.Errorf("failed to ping database after %d retries: %w", count, err)
					}
				} else {
					log.Println("Connected to the SQL database")
					return db, nil
				}
			}
		}
	}
}
