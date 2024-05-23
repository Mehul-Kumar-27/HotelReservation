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
			if err != nil {
				log.Fatalf("Error fetching the number of bookings to generate")
			}

			///////////////////////////////////////////////////////////////////////////////////////////

			databaseHandler := NewDatabaseHandlerForSQL(db)

			//////////////// Creating all the tables required
			log.Println("Creating Tables Required  ................")
			err = databaseHandler.databaseHandlerInterface.CreateUserTable()
			if err != nil {
				log.Fatalln("Error creating user table: ", err)
			}
			log.Println("✅ User table created successfully")

			err = databaseHandler.databaseHandlerInterface.CreateHotelTable()
			if err != nil {
				log.Fatalln("Error creating hotel table: ", err)
			}

			log.Println("✅ Hotel table created successfully")

			err = databaseHandler.databaseHandlerInterface.CreateBookingsTable()
			if err != nil {
				log.Fatalf("Error generating the bookings table %v", err)
			}

			log.Println("✅ Bookings table created successfully")

			err = databaseHandler.databaseHandlerInterface.CreateUserBookingTable()
			if err != nil {
				log.Fatalf("Error creating the user booking table %v", err)
			}

			log.Println("✅ User booking table created successfully")

			err = databaseHandler.databaseHandlerInterface.CreateReviewsTable()
			if err != nil {
				log.Fatalf("Error creating the reviews table %v", err)
			}

			log.Println("✅ Reviews table created successfully")

			////////////////////////////////////////////////////////////////////////////
			////////////   Generated all the tables required

			err = databaseHandler.databaseHandlerInterface.CreateFakeUsers(count)
			if err != nil {
				log.Fatalf("Error creating fake users: %v", err)
			}

			log.Println("Created the hotel table")

			log.Println("Adding fake data to the hotel table")
			err = databaseHandler.databaseHandlerInterface.CreateHotelFakeData(hotels)
			if err != nil {
				log.Fatalln("Error adding data to the hotel table: ", err)
			}

			////////////////////////////////////

			err = databaseHandler.databaseHandlerInterface.CreateBookingFakeData(bookings)
			if err != nil {
				log.Println("Error generating the fake booking data")
			}

			log.Printf("Successfully generated  the required data, exiting the datagenerator")
			<-cmd.Context().Done()

		},
	}
	cmdFlags := cmd.Flags()
	cmdFlags.Int("u", 50, "Number of users to generate")
	cmdFlags.Int("h", 15, "Number of hotels to generate")
	cmdFlags.Int("b", 5, "Number of reviews to generate for a hotel")
	return cmd
}

func NewCli(db *sql.DB) *cobra.Command {
	dgen := NewDataGeneratorCli(db)

	return dgen
}
