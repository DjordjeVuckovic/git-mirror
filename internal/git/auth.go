package git

import "github.com/go-git/go-git/v5/plumbing/transport/http"

type RepoAuth http.BasicAuth

func (ra *RepoAuth) String() string {
	//TODO implement me
	panic("implement me")
}

func (ra *RepoAuth) Name() string {
	panic("todo")
}
