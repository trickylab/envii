package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Trickster-ID/envii/internal/runner"
)

// runCmd: envii run <project> <env> -- <command...>
func runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run <project> <env> -- <command> [args...]",
		Short: "Run a command with the env's variables injected",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(_ *cobra.Command, args []string) error {
			projectName, envName := args[0], args[1]
			argv := args[2:]

			v, _, _, err := loadVault()
			if err != nil {
				return err
			}
			env, err := resolveEnv(v, projectName, envName)
			if err != nil {
				return err
			}

			code, err := runner.Run(env, argv)
			if err != nil {
				return err
			}
			os.Exit(code)
			return nil
		},
	}
}

// exportCmd: envii export <project> <env> [-o file]
func exportCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "export <project> <env>",
		Short: "Export an env as a .env file (stdout by default)",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			v, _, _, err := loadVault()
			if err != nil {
				return err
			}
			env, err := resolveEnv(v, args[0], args[1])
			if err != nil {
				return err
			}

			content := runner.Dotenv(env)
			if out == "" {
				fmt.Print(content)
				return nil
			}
			if err := os.WriteFile(out, []byte(content), 0o600); err != nil {
				return fmt.Errorf("write %s: %w", out, err)
			}
			fmt.Fprintf(os.Stderr, "wrote %s\n", out)
			return nil
		},
	}
	cmd.Flags().StringVarP(&out, "out", "o", "", "output file (default: stdout)")
	return cmd
}
