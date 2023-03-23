package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = &cobra.Command{
	Use:   "version",
	Short: "Show the version number of Gover",
	Long:  `Current Gover version with additional information`,
	Run: func(cmd *cobra.Command, args []string) {
		s := viper.GetString("github_username")
		x := viper.GetString("github_access_token")
		a := viper.GetString("GITHUB_ACCESS_TOKEN")
		println(s)
		println(a)
		println(x)
		fmt.Println("Gover API Generator v0.1 by @margostino")
	},
}
