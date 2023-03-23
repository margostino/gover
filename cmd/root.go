package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	root = &cobra.Command{
		Use:   "gover",
		Short: "A generator for GO Serverless functions hosted by Vercel",
		Long:  `Gover is a CLI for bootstrapping Go Serverless functions hosted by Vercel. This application is a tool to generate the needed files to quickly create a Go application.`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func Execute() error {
	return root.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	root.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gover.yml)")
	root.PersistentFlags().StringP("go-version", "", "1.18", "set Go version")
	root.PersistentFlags().StringP("github-username", "", "", "set github username")
	root.PersistentFlags().StringP("github-access-token", "", "", "set github access token")
	root.PersistentFlags().StringP("vercel-access-token", "", "", "set vercel access token")

	//root.PersistentFlags().StringP("name", "", "my-go-serverless-app", "set application name")
	//root.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	//root.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	//viper.BindPFlag("author", root.PersistentFlags().Lookup("author"))
	//viper.BindPFlag("useViper", root.PersistentFlags().Lookup("viper"))
	//viper.SetDefault("author", "maj.dagostino@gmail.com")
	//viper.SetDefault("license", "apache")

	create.PersistentFlags().StringP("name", "n", "my-go-serverless-app", "set application name")
	//create.PersistentFlags().StringP("github_username", "u", "", "set github username")
	//create.PersistentFlags().StringP("github_access_token", "p", "", "set github access token")

	//version.PersistentFlags().StringP("github_username", "u", "", "set github username")
	//version.PersistentFlags().StringP("github_access_token", "p", "", "set github access token")

	viper.BindPFlag("github_username", root.PersistentFlags().Lookup("github-username"))
	viper.BindPFlag("github_access_token", root.PersistentFlags().Lookup("github-access-token"))
	viper.BindPFlag("vercel_access_token", root.PersistentFlags().Lookup("vercel-access-token"))

	//viper.SetDefault("name", "maj.dagostino@gmail.com")
	root.AddCommand(version)
	root.AddCommand(create)
	//create.Flags().StringP("name", "n", "first-go-serverless-app", "set application name")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yml")
		viper.SetConfigName(".gover")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
