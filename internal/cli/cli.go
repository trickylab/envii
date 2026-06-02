// Package cli wires the cobra command tree.
package cli

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/Trickster-ID/envii/internal/model"
	"github.com/Trickster-ID/envii/internal/store"
	"github.com/Trickster-ID/envii/internal/tui"
)

var vaultPath string

// Execute runs the root command.
func Execute(version string) error {
	root := &cobra.Command{
		Use:     "envii",
		Short:   "Encrypted .env & secrets manager with a TUI",
		Version: version,
		RunE:    runTUI,
	}
	root.PersistentFlags().StringVar(&vaultPath, "vault", "", "path to vault file (default: ~/.config/envii/vault.age)")

	root.AddCommand(runCmd(), exportCmd())
	return root.Execute()
}

// runTUI is the default action: open the interactive UI.
func runTUI(cmd *cobra.Command, _ []string) error {
	s, err := store.New(vaultPath)
	if err != nil {
		return err
	}

	var vault *model.Vault
	var pass string

	if s.Exists() {
		pass, err = promptPassphrase("Passphrase: ")
		if err != nil {
			return err
		}
		vault, err = s.Load(pass)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("No vault found. Creating a new one.")
		pass, err = promptNewPassphrase()
		if err != nil {
			return err
		}
		vault = model.NewVault()
		if err := s.Save(vault, pass); err != nil {
			return err
		}
	}

	p := tea.NewProgram(tui.New(vault, s, pass), tea.WithAltScreen())
	_, err = p.Run()
	return err
}

// loadVault is a helper shared by non-TUI subcommands.
func loadVault() (*model.Vault, *store.Store, string, error) {
	s, err := store.New(vaultPath)
	if err != nil {
		return nil, nil, "", err
	}
	if !s.Exists() {
		return nil, nil, "", errors.New("no vault found; run `envii` to create one")
	}
	pass, err := promptPassphrase("Passphrase: ")
	if err != nil {
		return nil, nil, "", err
	}
	v, err := s.Load(pass)
	if err != nil {
		return nil, nil, "", err
	}
	return v, s, pass, nil
}

// resolveEnv looks up a project/env pair in the vault.
func resolveEnv(v *model.Vault, projectName, envName string) (*model.Env, error) {
	p := v.FindProject(projectName)
	if p == nil {
		return nil, fmt.Errorf("project %q not found", projectName)
	}
	e := p.FindEnv(envName)
	if e == nil {
		return nil, fmt.Errorf("env %q not found in project %q", envName, projectName)
	}
	return e, nil
}

func promptPassphrase(label string) (string, error) {
	// Allow bypassing the interactive prompt via env var (useful for demos/CI).
	if p := os.Getenv("ENVII_PASSPHRASE"); p != "" {
		return p, nil
	}
	fmt.Fprint(os.Stderr, label)
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", fmt.Errorf("read passphrase: %w", err)
	}
	return string(b), nil
}

func promptNewPassphrase() (string, error) {
	p1, err := promptPassphrase("New passphrase: ")
	if err != nil {
		return "", err
	}
	p2, err := promptPassphrase("Confirm passphrase: ")
	if err != nil {
		return "", err
	}
	if p1 != p2 {
		return "", errors.New("passphrases do not match")
	}
	if len(p1) < 1 {
		return "", errors.New("passphrase cannot be empty")
	}
	return p1, nil
}
