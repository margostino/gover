package cmd

import (
	"fmt"
	"github.com/margostino/gover/common"
	"github.com/margostino/gover/http"
	"github.com/margostino/gover/infra"
	"github.com/margostino/gover/repo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

type Params struct {
	AppName           string
	GithubUsername    string
	GithubAccessToken string
	VercelAccessToken string
	GoVersion         string
}

var httpClient = http.New()

type GitRepositoryRequest struct {
	Repo  string `json:"repo"`
	Rtype string `json:"type"`
}

type ProjectRequest struct {
	Name          string                `json:"name"`
	GitRepository *GitRepositoryRequest `json:"gitRepository"`
}

var create = &cobra.Command{
	Use:   "create",
	Short: "Create a new application",
	Long:  `Create a new Go Serverless application on Vercel`,
	Run: func(cmd *cobra.Command, args []string) {
		p := getParams(cmd)
		gotoLink := fmt.Sprintf("https://%s-%s.vercel.app/api/hello\n", p.AppName, p.GithubUsername)
		repository := repo.NewProject(p.AppName, p.GithubUsername, p.GithubAccessToken, p.GoVersion, httpClient)
		vercelProject := infra.NewProject(p.AppName, p.VercelAccessToken, httpClient)
		repository.Create()
		repository.CommitInitial(gotoLink)
		vercelProject.Create()
		repository.Bootstrap()
		log.Printf("⚡️  Go to %s\n", gotoLink)
	},
}

func getParams(cmd *cobra.Command) *Params {
	appName, err := cmd.Flags().GetString("name")
	common.Check(err)
	goVersion, err := cmd.Flags().GetString("go-version")
	githubUsername := getValueFor("GITHUB_USERNAME")
	githubAccessToken := getValueFor("GITHUB_ACCESS_TOKEN")
	vercelAccessToken := getValueFor("VERCEL_ACCESS_TOKEN")
	return &Params{
		AppName:           appName,
		GithubUsername:    githubUsername,
		GithubAccessToken: githubAccessToken,
		VercelAccessToken: vercelAccessToken,
		GoVersion:         goVersion,
	}
}

func getValueFor(key string) string {
	var value string
	value = viper.GetString(key)
	if &value == nil {
		common.Fatal(fmt.Sprintf("Cannot get value for key %s", key))
	}
	return value
}
