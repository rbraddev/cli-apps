/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"pScan/scan"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Aliases: []string{"l"},
	Use:     "list",
	Short:   "List hosts in host list",

	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")

		return listAction(os.Stdout, hostsFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listAction(out io.Writer, hostFile string, args []string) error {
	hl := &scan.HostList{}

	if err := hl.Load(hostFile); err != nil {
		return err
	}

	for _, h := range hl.Hosts {
		if _, err := fmt.Fprintln(out, h); err != nil {
			return err
		}
	}

	return nil
}
