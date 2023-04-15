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
	xbsFilePath string
	jsonOutPath string
)

// xbs2jsonCmd represents the xbs2json command
var xbs2jsonCmd = &cobra.Command{
	Use:   "xbs2json",
	Short: "xbs to json",
	Run: func(cmd *cobra.Command, args []string) {
		buffer, err := xbstools.LoadFile(xbsFilePath)
		cobra.CheckErr(err)
		jsonBuffer, err := xbstools.XBS2Json(buffer)
		cobra.CheckErr(err)
		os.MkdirAll(filepath.Dir(jsonOutPath), os.ModePerm)
		err = os.WriteFile(jsonOutPath, jsonBuffer, 0644)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(xbs2jsonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// xbs2jsonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// xbs2jsonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	xbs2jsonCmd.Flags().StringVarP(&xbsFilePath, "xbs", "i", "", "xbs file path")
	xbs2jsonCmd.Flags().StringVarP(&jsonOutPath, "json", "o", "", "xbs output path")
	xbs2jsonCmd.MarkFlagRequired("xbs")
	xbs2jsonCmd.MarkFlagRequired("json")
}
