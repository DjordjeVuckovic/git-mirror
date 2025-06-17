package mirror

import (
	"fmt"
	"github.com/DjordjeVuckovic/git-mirror/internal/config"
	"github.com/go-git/go-git/v5"
	gitcfg "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"os"
	"path/filepath"
	"strings"
)

type Mirror struct {
	config  config.MirrorConfig
	verbose bool
}

func New(cfg config.MirrorConfig, verbose bool) *Mirror {
	return &Mirror{
		config:  cfg,
		verbose: verbose,
	}
}

func (m *Mirror) Execute() error {
	if m.verbose {
		fmt.Printf("Starting mirror operation for %s\n", m.config.Name)
	}

	// Get authentication for target repository
	targetAuth, err := m.getAuth(m.config.Target.Auth)
	if err != nil {
		return fmt.Errorf("failed to setup target auth: %w", err)
	}

	// Get authentication for source repository
	sourceAuth, err := m.getAuth(m.config.Source.Auth)
	if err != nil {
		return fmt.Errorf("failed to setup source auth: %w", err)
	}

	// Clone target repository to memory
	if m.verbose {
		fmt.Printf("Cloning target repository: %s\n", m.config.Target.URL)
	}

	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:      m.config.Target.URL,
		Auth:     targetAuth,
		Mirror:   true,
		Progress: m.getProgressWriter(),
	})
	if err != nil {
		return fmt.Errorf("failed to clone target repository: %w", err)
	}

	// Create remote for source repository
	if m.verbose {
		fmt.Printf("Adding source remote: %s\n", m.config.Source.URL)
	}

	_, err = repo.CreateRemote(&gitcfg.RemoteConfig{
		Name: "source",
		URLs: []string{m.config.Source.URL},
	})
	if err != nil {
		return fmt.Errorf("failed to create source remote: %w", err)
	}

	// Push to source repository
	if m.verbose {
		fmt.Printf("Pushing to source repository: %s\n", m.config.Source.URL)
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "source",
		Auth:       sourceAuth,
		Progress:   m.getProgressWriter(),
	})
	if err != nil {
		return fmt.Errorf("failed to push to source repository: %w", err)
	}

	if m.verbose {
		fmt.Printf("Mirror operation completed successfully\n")
	}

	return nil
}

func (m *Mirror) getAuth(authConfig config.Auth) (transport.AuthMethod, error) {
	switch authConfig.Method {
	case "token":
		return &http.BasicAuth{
			Username: "token", // GitHub convention
			Password: authConfig.Token,
		}, nil
	case "basic":
		return &http.BasicAuth{
			Username: authConfig.Username,
			Password: authConfig.Password,
		}, nil
	case "ssh":
		if !strings.HasPrefix(authConfig.SSHKey, "/") {
			// Relative path, make it absolute
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get home directory: %w", err)
			}
			authConfig.SSHKey = filepath.Join(home, ".ssh", authConfig.SSHKey)
		}

		sshKey, err := ssh.NewPublicKeysFromFile("git", authConfig.SSHKey, "")
		if err != nil {
			return nil, fmt.Errorf("failed to load SSH key: %w", err)
		}
		return sshKey, nil
	case "none", "":
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported auth method: %s", authConfig.Method)
	}
}

func (m *Mirror) getProgressWriter() *os.File {
	if m.verbose {
		return os.Stdout
	}
	return nil
}
