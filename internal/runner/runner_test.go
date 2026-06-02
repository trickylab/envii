package runner

import (
	"testing"

	"github.com/Trickster-ID/envii/internal/model"
)

func TestDotenv(t *testing.T) {
	env := &model.Env{Vars: []*model.Var{
		{Key: "ZED", Value: "last"},
		{Key: "ABC", Value: "first"},
		{Key: "QUOTED", Value: "has space"},
	}}

	got := Dotenv(env)
	want := "ABC=first\nQUOTED=\"has space\"\nZED=last\n"
	if got != want {
		t.Fatalf("Dotenv mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestRunNoCommand(t *testing.T) {
	if _, err := Run(&model.Env{}, nil); err == nil {
		t.Fatal("expected error for empty argv")
	}
}
