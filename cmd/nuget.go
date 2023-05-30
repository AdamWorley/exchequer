/*
Copyright © 2023 Adam Worley adam.worley@netwealth.com
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nugetCmd represents the nuget command
var nugetCmd = &cobra.Command{
	Use:   "nuget",
	Short: "Manage NuGet packages",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("please specify a valid command")
	},
}

func init() {
	rootCmd.AddCommand(nugetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nugetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nugetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
