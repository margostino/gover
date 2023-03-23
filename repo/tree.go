package repo

var treeURLPattern = "%s/%s/%s/git/trees"

type RepoTree struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Sha  string `json:"sha"`
}

type TreeRequest struct {
	Tree     []*RepoTree `json:"tree"`
	BaseTree string      `json:"base_tree"`
}

func (p *Project) createTree() {
	shaResponse := &ShaPayload{}
	treeRequest := buildTreeRequest(p.initialSha, p.gitSha, p.goSha, p.vercelSha, p.helloSha)
	p.httpClient.Post(p.treeURL, treeRequest, p.accessToken, shaResponse)
	p.treeSha = shaResponse.Sha
}

func buildTreeRequest(initialCommitSHA string, gitIgnoreBlobSHA string, goModBlobSHA string, vercelBlobSHA string, helloApiBlobSHA string) *TreeRequest {
	return &TreeRequest{
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
