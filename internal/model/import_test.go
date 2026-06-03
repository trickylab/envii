package model

import "testing"

func TestImportVarsEmpty(t *testing.T) {
	env := &Env{Vars: []*Var{{Key: "A", Value: "old"}}}
	summary := ImportVars(env, nil, false)

	assertSummary(t, summary, ImportSummary{})
	assertMapValue(t, env, "A", "old")
}

func TestImportVarsIntoEmptyEnv(t *testing.T) {
	env := &Env{}
	summary := ImportVars(env, []*Var{
		{Key: "A", Value: "one"},
		{Key: "B", Value: "two", Secret: true},
	}, false)

	assertSummary(t, summary, ImportSummary{Added: 2})
	assertMapValue(t, env, "A", "one")
	assertMapValue(t, env, "B", "two")
	if !env.Vars[1].Secret {
		t.Fatal("secret flag was not preserved")
	}
}

func TestImportVarsSkipsDuplicates(t *testing.T) {
	env := &Env{Vars: []*Var{{Key: "A", Value: "old"}}}
	summary := ImportVars(env, []*Var{{Key: "A", Value: "new"}, {Key: "B", Value: "added"}}, false)

	assertSummary(t, summary, ImportSummary{Added: 1, Skipped: 1})
	assertMapValue(t, env, "A", "old")
	assertMapValue(t, env, "B", "added")
}

func TestImportVarsAllSkipped(t *testing.T) {
	env := &Env{Vars: []*Var{{Key: "A", Value: "old-a"}, {Key: "B", Value: "old-b"}}}
	summary := ImportVars(env, []*Var{{Key: "A", Value: "new-a"}, {Key: "B", Value: "new-b"}}, false)

	assertSummary(t, summary, ImportSummary{Skipped: 2})
	assertMapValue(t, env, "A", "old-a")
	assertMapValue(t, env, "B", "old-b")
}

func TestImportVarsOverwrite(t *testing.T) {
	env := &Env{Vars: []*Var{{Key: "A", Value: "old"}}}
	summary := ImportVars(env, []*Var{{Key: "A", Value: "new"}}, true)

	assertSummary(t, summary, ImportSummary{Overwritten: 1})
	assertMapValue(t, env, "A", "new")
}

func TestImportVarsAllOverwritten(t *testing.T) {
	env := &Env{Vars: []*Var{{Key: "A", Value: "old-a"}, {Key: "B", Value: "old-b"}}}
	summary := ImportVars(env, []*Var{{Key: "A", Value: "new-a"}, {Key: "B", Value: "new-b"}}, true)

	assertSummary(t, summary, ImportSummary{Overwritten: 2})
	assertMapValue(t, env, "A", "new-a")
	assertMapValue(t, env, "B", "new-b")
}

func TestImportVarsMixedOverwrite(t *testing.T) {
	env := &Env{Vars: []*Var{{Key: "A", Value: "old"}}}
	summary := ImportVars(env, []*Var{{Key: "A", Value: "new"}, {Key: "B", Value: "added"}}, true)

	assertSummary(t, summary, ImportSummary{Added: 1, Overwritten: 1})
	assertMapValue(t, env, "A", "new")
	assertMapValue(t, env, "B", "added")
}

func TestImportVarsPreservesOrder(t *testing.T) {
	env := &Env{Vars: []*Var{{Key: "A", Value: "old-a"}, {Key: "C", Value: "old-c"}}}
	ImportVars(env, []*Var{{Key: "A", Value: "new-a"}, {Key: "B", Value: "new-b"}}, true)

	want := []string{"A", "C", "B"}
	for i, key := range want {
		if env.Vars[i].Key != key {
			t.Fatalf("var %d key got %q, want %q", i, env.Vars[i].Key, key)
		}
	}
	assertMapValue(t, env, "A", "new-a")
}

func TestImportVarsDuplicateIncoming(t *testing.T) {
	env := &Env{}
	summary := ImportVars(env, []*Var{{Key: "A", Value: "first"}, {Key: "A", Value: "second"}}, false)

	assertSummary(t, summary, ImportSummary{Added: 1, Skipped: 1})
	if len(env.Vars) != 1 {
		t.Fatalf("got %d vars, want 1", len(env.Vars))
	}
	assertMapValue(t, env, "A", "first")
}

func TestImportVarsDuplicateIncomingOverwrite(t *testing.T) {
	env := &Env{}
	summary := ImportVars(env, []*Var{{Key: "A", Value: "first"}, {Key: "A", Value: "second"}}, true)

	assertSummary(t, summary, ImportSummary{Added: 1, Overwritten: 1})
	if len(env.Vars) != 1 {
		t.Fatalf("got %d vars, want 1", len(env.Vars))
	}
	assertMapValue(t, env, "A", "second")
}

func assertSummary(t *testing.T, got, want ImportSummary) {
	t.Helper()
	if got != want {
		t.Fatalf("summary got %+v, want %+v", got, want)
	}
}

func assertMapValue(t *testing.T, env *Env, key, want string) {
	t.Helper()
	if got := env.Map()[key]; got != want {
		t.Fatalf("%s got %q, want %q", key, got, want)
	}
}
