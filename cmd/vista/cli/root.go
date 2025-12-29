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

			command, _ := cmd.Flags().GetString("command")
			interactive, _ := cmd.Flags().GetBool("interactive")

			if command != "" {
				if _, err := manager.ExecuteCommandWithStreams(cmd.Context(), command); err != nil {
					return fmt.Errorf("failed to execute command: %w", err)
				}

				// If interactive is 'false' (default), close immediately to avoid running the TUI
				if !interactive {
					return nil
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

			if err := manager.Shutdown(); err != nil {
				return fmt.Errorf("failed to properly shutdown vfs: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringP("command", "C", "", "")
	cmd.Flags().BoolP("interactive", "i", false, "Used in combination with -C to define, if vista closes immediately after execution or not (default is 'false').")

	cmd.Version = fmt.Sprintf("%s.%s", info.Version, info.Commit)

	return cmd
}
