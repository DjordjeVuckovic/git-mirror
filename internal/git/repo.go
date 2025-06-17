package git

import (
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"log"
	"os"
	"path/filepath"
)

type RepoConfig struct {
	Path   string
	Branch string
	URL    string
	Auth   *RepoAuth
}

func CloneOrOpenPrivate(cfg RepoConfig) (*git.Repository, error) {
	if _, err := os.Stat(filepath.Join(cfg.Path, ".git")); os.IsNotExist(err) {
		log.Printf("Cloning private repo from %s...", cfg.URL)
		return git.PlainClone(cfg.Path, false, &git.CloneOptions{
			URL:           cfg.URL,
			ReferenceName: plumbing.NewBranchReferenceName(cfg.Branch),
			SingleBranch:  true,
			Auth:          cfg.Auth,
		})
	}
	log.Println("Opening existing private repo...")
	return git.PlainOpen(cfg.Path)
}
