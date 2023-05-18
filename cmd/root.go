/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "p2ppubsub",
	Short: "p2ppubsub contains two different implementation of the same pubsub logic using libp2p and ipfs.",
	Long:  `p2ppubsub contains two different implementation of the same pubsub logic using libp2p and ipfs.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
