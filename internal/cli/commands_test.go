package cli

import (
	"bufio"
	"strings"
	"testing"

	"github.com/trickylab/envii/internal/model"
)

func TestPromptProjectSelectExistingByNumber(t *testing.T) {
	api := &model.Project{Name: "api"}
	web := &model.Project{Name: "web"}
	vault := &model.Vault{Projects: []*model.Project{api, web}}

	got, err := promptProject(reader("2\n"), vault)
	if err != nil {
		t.Fatal(err)
	}
	if got != web {
		t.Fatalf("got %+v, want web", got)
	}
	if len(vault.Projects) != 2 {
		t.Fatalf("project count got %d, want 2", len(vault.Projects))
	}
}

func TestPromptProjectSelectExistingByName(t *testing.T) {
	api := &model.Project{Name: "api"}
	vault := &model.Vault{Projects: []*model.Project{api}}

	got, err := promptProject(reader("api\n"), vault)
	if err != nil {
		t.Fatal(err)
	}
	if got != api {
		t.Fatalf("got %+v, want api", got)
	}
	if len(vault.Projects) != 1 {
		t.Fatalf("project count got %d, want 1", len(vault.Projects))
	}
}

func TestPromptProjectCreateNew(t *testing.T) {
	vault := &model.Vault{Projects: []*model.Project{{Name: "api"}}}

	got, err := promptProject(reader("web\n"), vault)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "web" {
		t.Fatalf("got %q, want web", got.Name)
	}
	if len(vault.Projects) != 2 || vault.Projects[1] != got {
		t.Fatalf("new project not appended")
	}
}

func TestPromptProjectRejectsEmpty(t *testing.T) {
	_, err := promptProject(reader("\n"), &model.Vault{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPromptEnvSelectExistingByNumber(t *testing.T) {
	dev := &model.Env{Name: "dev"}
	prod := &model.Env{Name: "prod"}
	project := &model.Project{Name: "api", Envs: []*model.Env{dev, prod}}

	got, err := promptEnv(reader("2\n"), project)
	if err != nil {
		t.Fatal(err)
	}
	if got != prod {
		t.Fatalf("got %+v, want prod", got)
	}
	if len(project.Envs) != 2 {
		t.Fatalf("env count got %d, want 2", len(project.Envs))
	}
}

func TestPromptEnvSelectExistingByName(t *testing.T) {
	dev := &model.Env{Name: "dev"}
	project := &model.Project{Name: "api", Envs: []*model.Env{dev}}

	got, err := promptEnv(reader("dev\n"), project)
	if err != nil {
		t.Fatal(err)
	}
	if got != dev {
		t.Fatalf("got %+v, want dev", got)
	}
	if len(project.Envs) != 1 {
		t.Fatalf("env count got %d, want 1", len(project.Envs))
	}
}

func TestPromptEnvCreateNew(t *testing.T) {
	project := &model.Project{Name: "api", Envs: []*model.Env{{Name: "dev"}}}

	got, err := promptEnv(reader("staging\n"), project)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "staging" {
		t.Fatalf("got %q, want staging", got.Name)
	}
	if len(project.Envs) != 2 || project.Envs[1] != got {
		t.Fatalf("new env not appended")
	}
}

func TestPromptEnvRejectsEmpty(t *testing.T) {
	_, err := promptEnv(reader("\n"), &model.Project{Name: "api"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPromptLineAcceptsEOFWithInput(t *testing.T) {
	got, err := promptLine(reader("project"), "Project: ")
	if err != nil {
		t.Fatal(err)
	}
	if got != "project" {
		t.Fatalf("got %q, want project", got)
	}
}

func TestPromptLineRejectsEmptyEOF(t *testing.T) {
	_, err := promptLine(reader(""), "Project: ")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestChoiceHelpers(t *testing.T) {
	vault := &model.Vault{Projects: []*model.Project{{Name: "api"}}}
	if got, ok := projectChoice(vault, " 1 "); !ok || got.Name != "api" {
		t.Fatalf("expected project choice to select api")
	}
	for _, choice := range []string{"0", "2", "abc", ""} {
		if got, ok := projectChoice(vault, choice); ok || got != nil {
			t.Fatalf("choice %q got %+v/%v, want nil/false", choice, got, ok)
		}
	}

	project := &model.Project{Envs: []*model.Env{{Name: "dev"}}}
	if got, ok := envChoice(project, " 1 "); !ok || got.Name != "dev" {
		t.Fatalf("expected env choice to select dev")
	}
	for _, choice := range []string{"0", "2", "abc", ""} {
		if got, ok := envChoice(project, choice); ok || got != nil {
			t.Fatalf("choice %q got %+v/%v, want nil/false", choice, got, ok)
		}
	}
}

func reader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}
