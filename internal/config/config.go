package config

import (
	"os"
	"path/filepath"
	"runtime"
)

var maFilesPath string

// GetMaFilesPath returns the path to the maFiles directory
func GetMaFilesPath() string {
	if maFilesPath != "" {
		return maFilesPath
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	// Search order depends on the OS
	var searchPaths []string

	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		searchPaths = []string{
			filepath.Join(appData, "steamguard-cli", "maFiles"),
			filepath.Join(homeDir, "maFiles"),
		}
	} else {
		// Linux and other Unix-like systems
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			configHome = filepath.Join(homeDir, ".config")
		}
		searchPaths = []string{
			filepath.Join(configHome, "steamguard-cli", "maFiles"),
			filepath.Join(homeDir, "maFiles"),
		}
	}

	// Check existing paths
	for _, path := range searchPaths {
		if _, err := os.Stat(filepath.Join(path, "manifest.json")); err == nil {
			return path
		}
	}

	// Return the first path by default
	return searchPaths[0]
}

// SetMaFilesPath sets a custom path to maFiles
func SetMaFilesPath(path string) {
	maFilesPath = path
}
