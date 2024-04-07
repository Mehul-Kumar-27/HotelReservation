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

		},
	}
	cmdFlags := cmd.Flags()
	cmdFlags.Int("u", 20_000, "Number of users to generate")
	cmdFlags.Int("h", 5000, "Number of hotels to generate")
	return cmd
}

func NewCli(db *sql.DB) *cobra.Command {
	dgen := NewDataGeneratorCli(db)

	return dgen
}
