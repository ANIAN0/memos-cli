package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	fbconfig "github.com/ANIAN0/memos-cli/internal/config"
	sharedconfig "github.com/ANIAN0/memos-cli/pkg/config"
	"github.com/ANIAN0/memos-cli/pkg/output"
	"github.com/ANIAN0/memos-cli/pkg/version"
)

var (
	cfgFile    string
	jsonMode   bool
	verbose    bool
	noColor    bool
	timeoutSec int
)

var rootCmd = &cobra.Command{
	Use:     "memos-cli",
	Short:   "CLI for Memos note management",
	Long:    "memos-cli provides a shell-callable interface to the Memos HTTP API.",
	Version: version.String(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate timeout
		if timeoutSec <= 0 {
			return fmt.Errorf("timeout must be positive, got %d", timeoutSec)
		}

		// Config loading will be done by subcommands that need it
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command and returns the exit code.
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return output.ExitClientError
	}
	return output.ExitSuccess
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	rootCmd.PersistentFlags().BoolVar(&jsonMode, "json", false, "output JSON format")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colors")
	rootCmd.PersistentFlags().IntVar(&timeoutSec, "timeout", 60, "HTTP timeout in seconds")

	// Set version template
	rootCmd.SetVersionTemplate(`memos-cli {{.Version}}
`)
}

// getTimeout returns the timeout as a time.Duration.
func getTimeout() time.Duration {
	return time.Duration(timeoutSec) * time.Second
}

// loadConfig loads the configuration from file.
func loadConfig() (*fbconfig.Config, error) {
	binaryPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("get executable path: %w", err)
	}

	// Build args list with --config flag if specified
	args := os.Args[1:]
	if cfgFile != "" {
		args = append([]string{"--config", cfgFile}, args...)
	}

	// Build env map
	env := make(map[string]string)
	for _, e := range os.Environ() {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				env[e[:i]] = e[i+1:]
				break
			}
		}
	}

	// Load shared config
	result, err := sharedconfig.LoadConfig("memos-cli", args, env, binaryPath, nil)
	if err != nil {
		return nil, err
	}

	// Parse Memos-specific fields
	data, err := os.ReadFile(result.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	cfg, err := fbconfig.LoadFromBytes(data)
	if err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}