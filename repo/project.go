package repo

import (
	"context"
	"fmt"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/gover/common"
	"github.com/margostino/gover/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"log"
)

const (
	DefaultEncoding = "utf-8"
	BaseGithubURL   = "https://api.github.com/repos"
)

type VercelGithubConfig struct {
	Silent  bool `json:"silent"`
	Enabled bool `json:"enabled"`
}

type VercelFunctionConfig struct {
	Memory      int `json:"memory"`
	MaxDuration int `json:"maxDuration"`
}

type VercelProjectGitRequest struct {
	Repo string `json:"repo"`
	Type string `json:"type"`
}

type VercelProjectRequest struct {
	Name          string                   `json:"name"`
	GitRepository *VercelProjectGitRequest `json:"gitRepository"`
}

type VercelConfig struct {
	Github    *VercelGithubConfig              `json:"github"`
	Functions map[string]*VercelFunctionConfig `json:"functions"`
}

type HttpRepoCommitRequest struct {
	Tree    string   `json:"tree"`
	Message string   `json:"message"`
	Parents []string `json:"parents"`
}

type HttpRepoBlobResponse struct {
	Sha string `json:"sha"`
}

type ShaPayload struct {
	Sha string `json:"sha"`
}

type HttpRepoTreeResponse struct {
	Sha string `json:"sha"`
}

type HttpRepoCommitResponse struct {
	Sha string `json:"sha"`
}

type Params struct {
	AppName           string
	GithubUsername    string
	GithubAccessToken string
	VercelAccessToken string
	GoVersion         string
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

type Project struct {
	name           string
	username       string
	accessToken    string
	githubClient   *github.Client
	httpClient     *http.HttpClient
	goVersion      string
	blobURL        string
	treeURL        string
	commitURL      string
	patchMasterURL string
	initialSha     string
	gitSha         string
	vercelSha      string
	goSha          string
	helloSha       string
	treeSha        string
}

func NewProject(name string, username string, accessToken string, goVersion string, httpClient *http.HttpClient) *Project {
	return &Project{
		name:           name,
		username:       username,
		accessToken:    accessToken,
		githubClient:   getGithubClient(accessToken),
		httpClient:     httpClient,
		goVersion:      goVersion,
		blobURL:        fmt.Sprintf(blobUrlPattern, BaseGithubURL, username, name),
		treeURL:        fmt.Sprintf(treeURLPattern, BaseGithubURL, username, name),
		commitURL:      fmt.Sprintf(commitUrlPattern, BaseGithubURL, username, name),
		patchMasterURL: fmt.Sprintf(patchMasterUrlPattern, BaseGithubURL, username, name),
	}
}

func (p *Project) Create() {
	repo := &github.Repository{
		Name: &p.name,
	}

	repository, response, err := p.githubClient.Repositories.Create(context.Background(), "", repo)
	common.Check(err)

	if response.StatusCode == 201 {
		log.Printf("✅  Repository %s created successfully", repository.GetHTMLURL())
	} else {
		common.Fatal(fmt.Sprintf("Repository cannot be created. Got status code: %s", response.Status))
	}
}

func (p *Project) CommitInitial(gotoLink string) {
	message := "initial commit"
	options := &github.RepositoryContentFileOptions{
		Content: []byte(getReadmeContent(p.name, gotoLink)),
		Message: &message,
	}

	file := "README.md"
	initialCommit, response, err := p.githubClient.Repositories.CreateFile(context.Background(), p.username, p.name, file, options)
	common.Check(err)

	if response.StatusCode == 201 {
		log.Printf("✅  Successful Initial commit: %s\n", initialCommit.GetSHA())
	} else {
		common.Fatal("Initial commit failed")
	}

	p.initialSha = initialCommit.GetSHA()
}

func (p *Project) Bootstrap() {
	p.createBlobs()
	p.createTree()
	p.commit()
	log.Println("✅  Vercel bootstrap done successfully")
}

func getGithubClient(accessToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func getValueFor(key string) string {
	var value string
	value = viper.GetString(key)
	if &value == nil {
		common.Fatal(fmt.Sprintf("Cannot get value for key %s", key))
	}
	return value
}

func getReadmeContent(appName string, gotoLink string) string {
	return fmt.Sprintf("# %s\n\n"+
		"Golang Serverless App hosted by Vercel\n\n"+
		"Go to %s",
		appName,
		gotoLink)
}
