package repo

import (
	"encoding/json"
	"fmt"
	"github.com/margostino/gover/common"
)

var blobUrlPattern = "%s/%s/%s/git/blobs"

type BlobRequest struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func (p *Project) createBlobs() {
	shaResponse := &ShaPayload{}
	gitRequest := buildGitBlobRequest()
	goRequest := buildGoBlobRequest(p.username, p.name, p.goVersion)
	vercelRequest := buildVercelBlobRequest()
	helloRequest := buildHelloBlobRequest()

	p.httpClient.Post(p.blobURL, gitRequest, p.accessToken, shaResponse)
	p.gitSha = shaResponse.Sha

	p.httpClient.Post(p.blobURL, goRequest, p.accessToken, shaResponse)
	p.goSha = shaResponse.Sha

	p.httpClient.Post(p.blobURL, vercelRequest, p.accessToken, shaResponse)
	p.vercelSha = shaResponse.Sha

	p.httpClient.Post(p.blobURL, helloRequest, p.accessToken, shaResponse)
	p.helloSha = shaResponse.Sha
}

func buildBlobRequest(content string) *BlobRequest {
	return &BlobRequest{
		Content:  content,
		Encoding: DefaultEncoding,
	}
}

func buildHelloBlobRequest() *BlobRequest {
	content := "package api\n\n" +
		"import (\n" +
		"	\"fmt\"\n" +
		"	\"net/http\"\n" +
		")\n\n" +
		"func Hello(w http.ResponseWriter, r *http.Request) {\n" +
		"	fmt.Fprintf(w, \"Hello World!\\n\")\n" +
		"}"
	return buildBlobRequest(content)
}

func buildGitBlobRequest() *BlobRequest {
	return buildBlobRequest(".vercel\n.env")
}

func buildVercelBlobRequest() *BlobRequest {
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
	return buildBlobRequest(string(jsonVercelConfig))
}

func buildGoBlobRequest(username string, appName string, version string) *BlobRequest {
	content := fmt.Sprintf("module github.com/%s/%s\n\ngo %s", username, appName, version)
	return buildBlobRequest(content)
}
