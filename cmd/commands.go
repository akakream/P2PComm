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
		serverType, err := cmd.PersistentFlags().GetString("servertype")
		if err != nil {
			panic(err)
		}
		dataPath, err := cmd.PersistentFlags().GetString("data")
		if err != nil {
			panic(err)
		}
		useDatastore, err := cmd.PersistentFlags().GetBool("datastore")
		if err != nil {
			panic(err)
		}

		s := server.NewServer(port, serverType, dataPath, useDatastore)
		s.Start()
	},
}

func init() {
	serverCmd.PersistentFlags().StringP("port", "p", "3001", "give the port where the server runs")
	serverCmd.PersistentFlags().StringP("servertype", "s", "libp2p", "give the type of the server: libp2p or ipfs")
	serverCmd.PersistentFlags().StringP("data", "d", "./data", "give the path to the data folder")
	serverCmd.PersistentFlags().BoolP("datastore", "t", false, "true if you want to use datastore")
	rootCmd.AddCommand(serverCmd)
}
