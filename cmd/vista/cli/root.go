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
		Use:           "gosync",
		Short:         "GoSync Sync S3 Client",
		Long:          "A production-ready, cross-platform sync client for S3 (MinIO) that provides true bidirectional synchronization similar to MegaSync, Dropbox, or OneDrive.",
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

	cmd.Version = fmt.Sprintf("%s.%s", info.Version, info.Commit)
	return cmd
}
