/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"
	"xbsrebuild/xbstools"

	"github.com/spf13/cobra"
)

var (
	jsonFilePath string
	xbsOutPath   string
)

// json2xbsCmd represents the json2xbs command
var json2xbsCmd = &cobra.Command{
	Use:   "json2xbs",
	Short: "json to xbs",
	Run: func(cmd *cobra.Command, args []string) {
		buffer, err := xbstools.LoadFile(jsonFilePath)
		cobra.CheckErr(err)
		xbsBuffer, err := xbstools.Json2XBS(buffer)
		cobra.CheckErr(err)
		os.MkdirAll(filepath.Dir(jsonOutPath), os.ModePerm)
		err = os.WriteFile(xbsOutPath, xbsBuffer, 0644)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(json2xbsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// json2xbsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// json2xbsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	json2xbsCmd.Flags().StringVarP(&jsonFilePath, "json", "i", "", "json file path")
	json2xbsCmd.Flags().StringVarP(&xbsOutPath, "xbs", "o", "", "xbs output path")
	json2xbsCmd.MarkFlagRequired("json")
	json2xbsCmd.MarkFlagRequired("xbs")
}
