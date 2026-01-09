package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mwantia/vista/internal/app"
	"github.com/mwantia/vista/internal/vfs"
	"github.com/spf13/cobra"
)

func NewRootCommand(info VersionInfo) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "vista <uri>",
		Short:         "",
		Long:          "",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var uri string
			if len(args) == 1 {
				uri = strings.TrimSpace(args[0])
			}

			manager, err := vfs.NewManager(cmd.Context(), uri)
			if err != nil {
				return fmt.Errorf("failed to initialize vfs: %w", err)
			}

			interactive, _ := cmd.Flags().GetBool("interactive")
			script, _ := cmd.Flags().GetString("script")
			command, _ := cmd.Flags().GetString("command")

			if script != "" {
				// Fail if both parts are defined
				if command != "" {
					return fmt.Errorf("failed to parse: Using --script together with --command is not supported")
				}

				verbose, _ := cmd.Flags().GetBool("script-verbose")
				continu, _ := cmd.Flags().GetBool("script-continue")
				// Build source command
				command = fmt.Sprintf("source %s", script)
				if verbose {
					command += " -v"
				}
				if continu {
					command += " -c"
				}
			}

			if command != "" {
				if _, err := manager.ExecuteCommandWithStreams(cmd.Context(), command); err != nil {
					return fmt.Errorf("failed to execute command: %w", err)
				}
				// If interactive is 'false' (default), close immediately to avoid running the TUI
				if !interactive {
					return manager.Shutdown()
				}
			}

			model := app.New(manager)
			opts := []tea.ProgramOption{
				tea.WithContext(cmd.Context()),
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			}

			p := tea.NewProgram(model, opts...)
			if _, err := p.Run(); err != nil {
				return fmt.Errorf("error running vista: %w", err)
			}

			return nil
		},
	}

	// Existing flags
	cmd.Flags().BoolP("interactive", "i", false, "Keep TUI open after -C command (default is 'false')")
	cmd.Flags().StringP("command", "C", "", "Execute a single command and exit")
	cmd.Flags().StringP("script", "S", "", "Execute commands from a script file and exit")

	cmd.Flags().Bool("script-verbose", false, "Print each script command before executing (requires -S)")
	cmd.Flags().Bool("script-continue", false, "Continue script execution on errors (requires -S)")

	cmd.Version = fmt.Sprintf("%s.%s", info.Version, info.Commit)

	return cmd
}
