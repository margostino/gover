package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/gover/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

const (
	DefaultEncoding  = "utf-8"
	BaseGithubURL    = "https://api.github.com/repos"
	VercelProjectURL = "https://api.vercel.com/v9/projects"
)

var githubBlobURLPattern = "%s/%s/%s/git/blobs"
var githubTreeURLPattern = "%s/%s/%s/git/trees"
var githubCommitURLPattern = "%s/%s/%s/git/commits"
var githubPatchMasterURLPattern = "%s/%s/%s/git/refs/heads/master"

type HttpRepoBlobRequest struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

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

type RepoTree struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Sha  string `json:"sha"`
}

type HttpRepoTreeRequest struct {
	Tree     []*RepoTree `json:"tree"`
	BaseTree string      `json:"base_tree"`
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

func initialCommit(params *Params, githubClient *github.Client) string {
	message := "initial commit"
	options := &github.RepositoryContentFileOptions{
		Content: []byte(getReadmeContent(params.AppName, params.GithubUsername)),
		Message: &message,
	}

	file := "README.md"
	initialCommit, response, err := githubClient.Repositories.CreateFile(context.Background(), params.GithubUsername, params.AppName, file, options)
	common.Check(err)

	if response.StatusCode == 201 {
		log.Printf("✅  Successful Initial commit: %s\n", initialCommit.GetSHA())
	} else {
		common.Fatal("Initial commit failed")
	}

	return initialCommit.GetSHA()
}

func createRepo(appName string, githubClient *github.Client) {
	repo := &github.Repository{
		Name: &appName,
	}

	repository, response, err := githubClient.Repositories.Create(context.Background(), "", repo)
	common.Check(err)

	if response.StatusCode == 201 {
		log.Printf("✅  Repository %s created successfully", repository.GetHTMLURL())
	} else {
		common.Fatal(fmt.Sprintf("Repository cannot be created. Got status code: %s", response.Status))
	}
}

func getGithubRequest(method string, url string, data interface{}, accessToken string) *http.Request {
	body, err := json.Marshal(data)
	common.Check(err)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	common.Check(err)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func postGithub(httpClient *http.Client, url string, data interface{}, accessToken string) string {
	request := getGithubRequest(http.MethodPost, url, data, accessToken)
	return callGithub(httpClient, request)
}

func patchGithub(httpClient *http.Client, url string, data interface{}, accessToken string) string {
	request := getGithubRequest(http.MethodPatch, url, data, accessToken)
	return callGithub(httpClient, request)
}

func callGithub(httpClient *http.Client, request *http.Request) string {
	var shaPayload ShaPayload
	response, err := httpClient.Do(request)
	common.Check(err)
	err = json.NewDecoder(response.Body).Decode(&shaPayload)
	common.Check(err)
	return shaPayload.Sha
}

func createVercelProject(httpClient *http.Client, appName string, accessToken string) {
	data := &VercelProjectRequest{
		Name: appName,
		GitRepository: &VercelProjectGitRequest{
			Repo: appName,
			Type: "github",
		},
	}
	body, err := json.Marshal(data)
	common.Check(err)
	request, err := http.NewRequest(http.MethodPost, VercelProjectURL, bytes.NewBuffer(body))
	common.Check(err)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Set("Content-Type", "application/json")

	response, err := httpClient.Do(request)
	common.Check(err)

	if response.StatusCode == 200 {
		log.Println("✅  Vercel Project created successfully")
	} else {
		common.Fatal(fmt.Sprintf("Vercel Project cannot be created. Got status code: %s", response.Status))
	}

}

var create = &cobra.Command{
	Use:   "create",
	Short: "Create a new application",
	Long:  `Create a new Go Serverless application on Vercel`,
	Run: func(cmd *cobra.Command, args []string) {
		httpClient := &http.Client{}
		params := getParams(cmd)
		githubClient := getGithubClient(params.GithubAccessToken)
		githubBlobURL := fmt.Sprintf(githubBlobURLPattern, BaseGithubURL, params.GithubUsername, params.AppName)
		githubTreeURL := fmt.Sprintf(githubTreeURLPattern, BaseGithubURL, params.GithubUsername, params.AppName)
		githubCommitURL := fmt.Sprintf(githubCommitURLPattern, BaseGithubURL, params.GithubUsername, params.AppName)
		githubPatchMasterURL := fmt.Sprintf(githubPatchMasterURLPattern, BaseGithubURL, params.GithubUsername, params.AppName)

		createRepo(params.AppName, githubClient)
		initialCommitSHA := initialCommit(params, githubClient)

		createVercelProject(httpClient, params.AppName, params.VercelAccessToken)
		gitIgnoreRequest := getGitIgnoreRequest()
		goModRequest := getGoModRequest(params.GithubUsername, params.AppName, params.GoVersion)
		vercelBlobRequest := getVercelBlobRequest()
		helloApiBlobRequest := getHelloApiRequest()

		gitIgnoreBlobSHA := postGithub(httpClient, githubBlobURL, gitIgnoreRequest, params.GithubAccessToken)
		goModBlobSHA := postGithub(httpClient, githubBlobURL, goModRequest, params.GithubAccessToken)
		vercelBlobSHA := postGithub(httpClient, githubBlobURL, vercelBlobRequest, params.GithubAccessToken)
		helloApiBlobSHA := postGithub(httpClient, githubBlobURL, helloApiBlobRequest, params.GithubAccessToken)

		treeRequest := getTreeRequest(initialCommitSHA, gitIgnoreBlobSHA, goModBlobSHA, vercelBlobSHA, helloApiBlobSHA)
		treeSHA := postGithub(httpClient, githubTreeURL, treeRequest, params.GithubAccessToken)

		commitRequest := getCommitRequest(initialCommitSHA, treeSHA)
		secondCommitSHA := postGithub(httpClient, githubCommitURL, commitRequest, params.GithubAccessToken)

		patchRequest := getUpdateRefRequest(secondCommitSHA)
		patchGithub(httpClient, githubPatchMasterURL, patchRequest, params.GithubAccessToken)
		log.Println("✅  Vercel bootstrap done successfully")
	},
}

func getGithubClient(githubAccessToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
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

func getReadmeContent(appName string, username string) string {
	return fmt.Sprintf("# %s\n\n"+
		"Golang Serverless App hosted by Vercel\n\n"+
		"Go to https://%s-%s.vercel.app/api/hello",
		appName,
		appName,
		username)
}

func getHelloApiRequest() *HttpRepoBlobRequest {
	content := "package api\n\n" +
		"import (\n" +
		"	\"fmt\"\n" +
		"	\"net/http\"\n" +
		")\n\n" +
		"func Hello(w http.ResponseWriter, r *http.Request) {\n" +
		"	fmt.Fprintf(w, \"Hello World!\\n\")\n" +
		"}"
	return &HttpRepoBlobRequest{
		Content:  content,
		Encoding: DefaultEncoding,
	}
}

func getGitIgnoreRequest() *HttpRepoBlobRequest {
	return &HttpRepoBlobRequest{
		Content:  ".vercel\n.env",
		Encoding: DefaultEncoding,
	}
}

func getVercelBlobRequest() *HttpRepoBlobRequest {
	vercelConfig := &VercelConfig{
		Github: &VercelGithubConfig{
			Silent:  true,
			Enabled: true,
		},
		Functions: map[string]*VercelFunctionConfig{
			"api/hello.go": {
				Memory:      1024,
				MaxDuration: 10,
			},
		},
	}
	jsonVercelConfig, err := json.MarshalIndent(vercelConfig, "", "    ")
	common.Check(err)
	return &HttpRepoBlobRequest{
		Content:  string(jsonVercelConfig),
		Encoding: DefaultEncoding,
	}
}

func getGoModRequest(username string, appName string, version string) *HttpRepoBlobRequest {
	return &HttpRepoBlobRequest{
		Content:  fmt.Sprintf("module github.com/%s/%s\n\ngo %s", username, appName, version),
		Encoding: "utf-8",
	}
}

func getCommitRequest(initialCommitSHA string, treeSHA string) *HttpRepoCommitRequest {
	return &HttpRepoCommitRequest{
		Tree:    treeSHA,
		Message: "second commit",
		Parents: []string{initialCommitSHA},
	}
}

func getUpdateRefRequest(sha string) *ShaPayload {
	return &ShaPayload{
		Sha: sha,
	}
}

func getTreeRequest(initialCommitSHA string, gitIgnoreBlobSHA string, goModBlobSHA string, vercelBlobSHA string, helloApiBlobSHA string) *HttpRepoTreeRequest {
	return &HttpRepoTreeRequest{
		Tree: []*RepoTree{
			{
				Path: ".gitignore",
				Mode: "100644",
				Type: "blob",
				Sha:  gitIgnoreBlobSHA,
			},
			{
				Path: "go.mod",
				Mode: "100644",
				Type: "blob",
				Sha:  goModBlobSHA,
			},
			{
				Path: "vercel.json",
				Mode: "100644",
				Type: "blob",
				Sha:  vercelBlobSHA,
			},
			{
				Path: "api/hello.go",
				Mode: "100644",
				Type: "blob",
				Sha:  helloApiBlobSHA,
			},
		},
		BaseTree: initialCommitSHA,
	}
}
