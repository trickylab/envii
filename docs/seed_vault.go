//go:build ignore

// seed_vault.go creates a demo vault for vhs recording.
// Usage: go run docs/seed_vault.go <vault-path>
package main

import (
	"fmt"
	"os"

	"github.com/trickylab/envii/internal/crypto"
	"github.com/trickylab/envii/internal/model"
	"github.com/trickylab/envii/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: seed_vault.go <path>")
		os.Exit(1)
	}

	// Low work factor so seeding is instant in the demo.
	crypto.SetWorkFactor(10)

	s := &store.Store{Path: os.Args[1]}

	v := model.NewVault()
	v.Projects = []*model.Project{
		{
			Name: "my-api",
			Envs: []*model.Env{
				{
					Name: "dev",
					Vars: []*model.Var{
						{Key: "PORT", Value: "8080"},
						{Key: "DATABASE_URL", Value: "postgres://localhost/myapi_dev", Secret: false},
						{Key: "API_SECRET", Value: "s3cr3t-dev-k3y", Secret: true},
						{Key: "LOG_LEVEL", Value: "debug"},
					},
				},
				{
					Name: "prod",
					Vars: []*model.Var{
						{Key: "PORT", Value: "443"},
						{Key: "DATABASE_URL", Value: "postgres://prod-host/myapi", Secret: false},
						{Key: "API_SECRET", Value: "sup3r-s3cur3-pr0d-k3y", Secret: true},
						{Key: "LOG_LEVEL", Value: "warn"},
					},
				},
			},
		},
		{
			Name: "frontend",
			Envs: []*model.Env{
				{
					Name: "dev",
					Vars: []*model.Var{
						{Key: "NEXT_PUBLIC_API_URL", Value: "http://localhost:8080"},
						{Key: "NEXTAUTH_SECRET", Value: "dev-auth-secret", Secret: true},
					},
				},
			},
		},
	}

	if err := s.Save(v, "demo"); err != nil {
		fmt.Fprintln(os.Stderr, "save:", err)
		os.Exit(1)
	}
	fmt.Println("vault seeded at", os.Args[1])
}
