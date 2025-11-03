package cmd

import (
	"fmt"
	"os"

	"github.com/devhooly/steamguard-fork/internal/steamapi"
	"github.com/spf13/cobra"
)

var (
	acceptAll bool
	rejectAll bool
)

var tradeCmd = &cobra.Command{
	Use:   "trade",
	Short: "Manage trade confirmations",
	Long:  `View and manage pending Steam trade confirmations.`,
	Run: func(cmd *cobra.Command, args []string) {
		if manifestMgr.IsEmpty() {
			fmt.Println("No accounts found. Use 'steamguard setup' to configure.")
			return
		}

		account, err := manifestMgr.GetAccount(username)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := steamapi.NewClient()
		
		// Get list of confirmations
		confirmations, err := client.GetConfirmations(account)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get confirmations: %v\n", err)
			os.Exit(1)
		}

		if len(confirmations) == 0 {
			fmt.Println("No pending confirmations.")
			return
		}

		fmt.Printf("Found confirmations: %d\n\n", len(confirmations))

		if acceptAll {
			// Accept all
			for _, conf := range confirmations {
				err := client.AcceptConfirmation(account, conf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to accept confirmation %s: %v\n", conf.ID, err)
					continue
				}
				fmt.Printf("✓ Accepted: %s\n", conf.Description)
			}
		} else if rejectAll {
			// Reject all
			for _, conf := range confirmations {
				err := client.RejectConfirmation(account, conf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to reject confirmation %s: %v\n", conf.ID, err)
					continue
				}
				fmt.Printf("✗ Rejected: %s\n", conf.Description)
			}
		} else {
			// Show list
			for i, conf := range confirmations {
				fmt.Printf("[%d] %s\n", i+1, conf.Description)
				fmt.Printf("    ID: %s\n", conf.ID)
				fmt.Printf("    Type: %s\n", conf.Type)
				fmt.Println()
			}
			fmt.Println("Use --accept or --reject to manage confirmations.")
		}
	},
}

func init() {
	rootCmd.AddCommand(tradeCmd)
	tradeCmd.Flags().BoolVar(&acceptAll, "accept", false, "Accept all confirmations")
	tradeCmd.Flags().BoolVar(&rejectAll, "reject", false, "Reject all confirmations")
}

