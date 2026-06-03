package cli

import (
	"strings"
	"testing"

	"github.com/Trickster-ID/envii/internal/model"
)

func TestResolveEnv(t *testing.T) {
	env := &model.Env{Name: "dev"}
	vault := &model.Vault{Projects: []*model.Project{{Name: "api", Envs: []*model.Env{env}}}}

	got, err := resolveEnv(vault, "api", "dev")
	if err != nil {
		t.Fatal(err)
	}
	if got != env {
		t.Fatalf("got %+v, want dev env", got)
	}
}

func TestResolveEnvErrors(t *testing.T) {
	vault := &model.Vault{Projects: []*model.Project{{Name: "api", Envs: []*model.Env{{Name: "dev"}}}}}

	tests := []struct {
		name string
		proj string
		env  string
		want string
	}{
		{name: "missing project", proj: "web", env: "dev", want: `project "web" not found`},
		{name: "missing env", proj: "api", env: "prod", want: `env "prod" not found in project "api"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := resolveEnv(vault, tt.proj, tt.env)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error got %q, want containing %q", err.Error(), tt.want)
			}
		})
	}
}
