package main

import "github.com/spf13/cobra"

func NewDataGeneratorCli() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dgen",
		Short: "Data generator for hotel reservation",
		Long: `This CLI application is used to generate the data
		for the hotel reservation system. It can be used to generate user, hotel, and booking data.`,
	}

	cmd.AddCommand(NewCli())
	return cmd
}

func NewUserDataGeneratorCLI() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Generate user data",
		Long:  `This command is used to generate user data for the hotel reservation system.`,
	}

	cmdFlags := cmd.Flags()
	cmdFlags.Int("users", 10, "Number of users to generate")
	return cmd
}

func NewCli() *cobra.Command {
	dgen := NewDataGeneratorCli()

	return dgen
}
