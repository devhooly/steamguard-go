package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long:  `Shows a list of all configured Steam Guard accounts.`,
	Run: func(cmd *cobra.Command, args []string) {
		if manifestMgr.IsEmpty() {
			fmt.Println("No accounts found. Use 'steamguard setup' to configure.")
			return
		}

		accounts := manifestMgr.GetAllAccounts()
		fmt.Printf("Found accounts: %d\n\n", len(accounts))
		
		for i, acc := range accounts {
			fmt.Printf("[%d] %s\n", i+1, acc.AccountName)
			if acc.DeviceID != "" {
				fmt.Printf("    Device ID: %s\n", acc.DeviceID)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

