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
			log.Panicln("Welocome to the data generator CLI.")
			log.Println("Creating user table..")
			databasehandler := NewDatabaseHandlerForSQL(db)
			err := databasehandler.databaseHandlerInterface.CreateUserTable()
			if err != nil {
				log.Fatalln("Error creating user table: ", err)
			}

		},
	}
	return cmd
}

func NewUserDataGeneratorCLI() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Generate user data",
		Long:  `This command is used to generate user data for the hotel reservation system.`,
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			count, err := flags.GetInt("count")
			if err != nil {
				log.Fatalf("Error while parsing count flag: %v", err)
			}

			log.Printf("Generating %d users", count)
		},
	}

	cmdFlags := cmd.Flags()
	cmdFlags.Int("count", 10, "Number of users to generate")
	return cmd
}

func NewCli(db *sql.DB) *cobra.Command {
	dgen := NewDataGeneratorCli(db)
	userCommands := NewUserDataGeneratorCLI()
	dgen.AddCommand(userCommands)

	return dgen
}
