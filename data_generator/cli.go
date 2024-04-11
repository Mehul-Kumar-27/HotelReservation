package main

import (
	"database/sql"
	"log"

	"github.com/spf13/cobra"
)

func NewDataGeneratorCli(db *sql.DB) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dgen",
		Short: "Data generator for hotel reservation",
		Long: `This CLI application is used to generate the data
		for the hotel reservation system. It can be used to generate user, hotel, and booking data.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Welocome to the data generator CLI.")
			flags := cmd.Flags()
			count, err := flags.GetInt("u")
			if err != nil {
				log.Fatalf("Error while parsing count flag: %v", err)
			}
			hotels, err := flags.GetInt("h")
			if err != nil {
				log.Fatalln("Error parsing hotel flag: ", err)
			}

			bookings, err := flags.GetInt("b")
			if err != nil{
				log.Fatalf("Error fetching the number of bookings to generate")
			}


			databaseHandler := NewDatabaseHandlerForSQL(db)
			log.Println("Creating user table..")
			err = databaseHandler.databaseHandlerInterface.CreateUserTable()
			if err != nil {
				log.Fatalln("Error creating user table: ", err)
			}

			log.Println("User table created successfully")
			err = databaseHandler.databaseHandlerInterface.CreateFakeUsers(count)
			if err != nil {
				log.Fatalf("Error creating fake users: %v", err)
			}

			////////////////////////////
			log.Println("Creating the hotel table")
			err = databaseHandler.databaseHandlerInterface.CreateHotelTable()
			if err != nil {
				log.Fatalln("Error creating user table: ", err)
			}
			log.Println("Created the hotel table")

			log.Println("Adding fake data to the hotel table")
			err = databaseHandler.databaseHandlerInterface.CreateHotelFakeData(hotels)
			if err != nil {
				log.Fatalln("Error adding data to the hotel table: ", err)
			}

			////////////////////////////////////
			log.Printf("Generating the bookings data")
			err = databaseHandler.databaseHandlerInterface.CreateBookingsTable()
			if err != nil {
				log.Fatalf("Error generating the bookings table %v", err)
			}

			log.Println("Bookings Table Generated")

			err = databaseHandler.databaseHandlerInterface.CreateBookingFakeData(bookings)
			if err != nil {
				log.Println("Error generating the fake booking data")
			}
			log.Println("Generated the fake booking data")
		},
	}
	cmdFlags := cmd.Flags()
	cmdFlags.Int("u", 20_000, "Number of users to generate")
	cmdFlags.Int("h", 5000, "Number of hotels to generate")
	cmdFlags.Int("b", 50, "Number of reviews to generate for a hotel")
	return cmd
}

func NewCli(db *sql.DB) *cobra.Command {
	dgen := NewDataGeneratorCli(db)

	return dgen
}
