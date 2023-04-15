/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"xbsrebuild/hub"

	"github.com/spf13/cobra"
)

var (
	server string
	port   string
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Run: func(cmd *cobra.Command, args []string) {
		os.RemoveAll("cache")
		os.MkdirAll("cache", os.ModePerm)
		hub.Run(server, port)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().StringVarP(&server, "server", "s", "0.0.0.0", "listen addr")
	serverCmd.Flags().StringVarP(&port, "port", "p", "8282", "listen port")
}
