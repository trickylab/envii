// Package runner executes commands with injected environment variables
// and renders .env exports.
package runner

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/Trickster-ID/envii/internal/model"
)

// Run executes argv with the env vars injected on top of the current
// process environment. It streams stdio and returns the command's exit code.
func Run(env *model.Env, argv []string) (int, error) {
	if len(argv) == 0 {
		return 1, fmt.Errorf("no command provided")
	}

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	for k, v := range env.Map() {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if ok := asExitError(err, &exitErr); ok {
			return exitErr.ExitCode(), nil
		}
		return 1, err
	}
	return 0, nil
}

func asExitError(err error, target **exec.ExitError) bool {
	if e, ok := err.(*exec.ExitError); ok {
		*target = e
		return true
	}
	return false
}

// Dotenv renders an env as .env file content, sorted by key.
func Dotenv(env *model.Env) string {
	vars := make([]*model.Var, len(env.Vars))
	copy(vars, env.Vars)
	sort.Slice(vars, func(i, j int) bool { return vars[i].Key < vars[j].Key })

	var b strings.Builder
	for _, v := range vars {
		b.WriteString(v.Key)
		b.WriteString("=")
		b.WriteString(quote(v.Value))
		b.WriteString("\n")
	}
	return b.String()
}

// quote wraps values containing whitespace or special chars in double quotes.
func quote(s string) string {
	if s == "" {
		return s
	}
	if strings.ContainsAny(s, " \t\n\"'#$") {
		escaped := strings.ReplaceAll(s, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		return `"` + escaped + `"`
	}
	return s
}
