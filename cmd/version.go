package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = &cobra.Command{
	Use:   "version",
	Short: "Show the version number of Gover",
	Long:  `Current Gover version with additional information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gover API Generator v0.1")
	},
}
