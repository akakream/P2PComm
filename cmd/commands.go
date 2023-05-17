package cmd

import (
	"errors"

	"github.com/akakream/sailorsailor/server"
	"github.com/spf13/cobra"
)

var (
	ErrRequired                          = errors.New("only one argument is required")
	ErrOnlyOneArgumentRequired           = errors.New("only one argument is required")
	ErrServerTypeAndTopicRequired        = errors.New("server type and topic are required in this order")
	ErrServerTypeTopicAndMessageRequired = errors.New("server type, topic and message are required in this order")
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  `Start server`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.PersistentFlags().GetString("port")
		if err != nil {
			panic(err)
		}
		s := server.NewServer(port)
		s.Start()
	},
}

func init() {
	serverCmd.PersistentFlags().StringP("port", "p", "3000", "give the port where the server runs")
	rootCmd.AddCommand(serverCmd)
}
