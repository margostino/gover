package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var create = &cobra.Command{
	Use:   "create",
	Short: "Create a new application",
	Long:  `Create a new Go Serverless application on Vercel`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gover API Generator v0.1")
	},
}
