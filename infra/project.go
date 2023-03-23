package infra

import (
	"github.com/margostino/gover/http"
	"log"
)

const ProjectURL = "https://api.vercel.com/v9/projects"

type Project struct {
	name        string
	accessToken string
	httpClient  *http.HttpClient
}

type GitRepositoryRequest struct {
	Repo  string `json:"repo"`
	Rtype string `json:"type"`
}

type ProjectRequest struct {
	Name          string                `json:"name"`
	GitRepository *GitRepositoryRequest `json:"gitRepository"`
}

func NewProject(name string, accessToken string, httpClient *http.HttpClient) *Project {
	return &Project{
		name:        name,
		accessToken: accessToken,
		httpClient:  httpClient,
	}
}

func (p *Project) Create() {
	var response interface{}
	data := &ProjectRequest{
		Name: p.name,
		GitRepository: &GitRepositoryRequest{
			Repo:  p.name,
			Rtype: "github",
		},
	}
	p.httpClient.Post(ProjectURL, data, p.accessToken, response)
	log.Println("✅  Vercel Project created successfully")
}
