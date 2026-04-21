// Package main is the entry point for the envoy-cli tool.
// It wires together the internal packages and exposes a command-line interface
// for managing and syncing .env files across environments.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envoy-cli/internal/config"
	"github.com/envoy-cli/internal/env"
	"github.com/envoy-cli/internal/storage"
	"github.com/envoy-cli/internal/sync"
)

var (
	// Version is set at build time via -ldflags.
	Version = "dev"
	// Commit is the git commit SHA set at build time.
	Commit = "none"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envoy",
		Short: "Manage and sync .env files across environments",
		Long: `envoy-cli is a lightweight tool for managing, encrypting,
and syncing .env files across local and remote environments.`,
		SilenceUsage: true,
	}

	root.AddCommand(
		newVersionCmd(),
		newPushCmd(),
		newPullCmd(),
		newDiffCmd(),
		newLintCmd(),
	)

	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the envoy version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "envoy %s (%s)\n", Version, Commit)
		},
	}
}

func newPushCmd() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push local .env file to the configured remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadFromCWD()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			local := storage.NewLocalBackend(file)
			remote := storage.NewLocalBackend(cfg.RemotePath)
			syncer := sync.NewSyncer(local, remote, nil)

			if err := syncer.Push(cmd.Context()); err != nil {
				return fmt.Errorf("push: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "pushed successfully")
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", ".env", "path to the local .env file")
	return cmd
}

func newPullCmd() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull remote .env file to the local path",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadFromCWD()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			local := storage.NewLocalBackend(file)
			remote := storage.NewLocalBackend(cfg.RemotePath)
			syncer := sync.NewSyncer(local, remote, nil)

			if err := syncer.Pull(cmd.Context()); err != nil {
				return fmt.Errorf("pull: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "pulled successfully")
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", ".env", "path to the local .env file")
	return cmd
}

func newDiffCmd() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Show differences between local and remote .env files",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadFromCWD()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			local := storage.NewLocalBackend(file)
			remote := storage.NewLocalBackend(cfg.RemotePath)
			syncer := sync.NewSyncer(local, remote, nil)

			diffs, err := syncer.Diff(cmd.Context())
			if err != nil {
				return fmt.Errorf("diff: %w", err)
			}

			if !env.HasDiff(diffs) {
				fmt.Fprintln(cmd.OutOrStdout(), "no differences found")
				return nil
			}

			for _, d := range diffs {
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", d.Type, d.Key)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", ".env", "path to the local .env file")
	return cmd
}

func newLintCmd() *cobra.Command {
	var file string
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint a .env file for common issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("open %s: %w", file, err)
			}
			defer f.Close()

			parser := env.NewParser()
			entries, err := parser.Parse(f)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			linter := env.NewLinter()
			violations := linter.Lint(entries)
			report := env.NewLintReport(violations)

			if jsonOut {
				return report.WriteJSON(cmd.OutOrStdout())
			}
			return report.WriteText(cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", ".env", "path to the .env file to lint")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output results as JSON")
	return cmd
}
