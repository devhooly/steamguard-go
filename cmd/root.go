package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devhooly/steamguard-fork/internal/config"
	"github.com/devhooly/steamguard-fork/internal/manifest"
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	username    string
	manifestMgr *manifest.Manager
)

var rootCmd = &cobra.Command{
	Use:   "steamguard",
	Short: "Steam Guard CLI - utility for generating 2FA codes and managing Steam confirmations",
	Long: `steamguard-cli - command line utility for working with Steam Mobile Authenticator.
	
Features:
  - Generate 2FA codes for Steam login
  - Manage trade and market confirmations
  - Encrypted storage of 2FA secrets
  - Generate QR codes for importing into other applications
  - Compatible with Steam Desktop Authenticator's maFiles format`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize manifest manager
		maFilesPath := config.GetMaFilesPath()
		var err error
		manifestMgr, err = manifest.NewManager(maFilesPath)
		if err != nil {
			return fmt.Errorf("failed to load manifest: %w", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// By default, generate code for the first account or specified username
		if manifestMgr.IsEmpty() {
			fmt.Println("No accounts found. Use 'steamguard setup' to configure.")
			return
		}

		account, err := manifestMgr.GetAccount(username)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		code, err := account.GenerateCode()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate code: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(code)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Steam username")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to configuration file")
}

func initConfig() {
	if cfgFile != "" {
		// Use the specified path
		config.SetMaFilesPath(filepath.Dir(cfgFile))
	}
}

