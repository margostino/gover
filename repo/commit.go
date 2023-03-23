package repo

var commitUrlPattern = "%s/%s/%s/git/commits"
var patchMasterUrlPattern = "%s/%s/%s/git/refs/heads/master"

func (p *Project) commit() {
	shaResponse := &ShaPayload{}
	commitRequest := buildCommitRequest(p.initialSha, p.treeSha)
	p.httpClient.Post(p.commitURL, commitRequest, p.accessToken, shaResponse)
	commitSha := shaResponse.Sha
	patchRequest := buildUpdateRefRequest(commitSha)
	p.httpClient.Patch(p.patchMasterURL, patchRequest, p.accessToken, shaResponse)
}

func buildCommitRequest(initialCommitSHA string, treeSHA string) *HttpRepoCommitRequest {
	return &HttpRepoCommitRequest{
		Tree:    treeSHA,
		Message: "second commit",
		Parents: []string{initialCommitSHA},
	}
}

func buildUpdateRefRequest(sha string) *ShaPayload {
	return &ShaPayload{
		Sha: sha,
	}
}
