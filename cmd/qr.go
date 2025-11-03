package cmd

import (
	"fmt"
	"os"

	"github.com/devhooly/steamguard-go/internal/qrcode"
	"github.com/spf13/cobra"
)

var qrCmd = &cobra.Command{
	Use:   "qr",
	Short: "Generate QR code to import 2FA secret into other applications",
	Long: `Generates a QR code with your 2FA secret that can be scanned 
in other password management applications (e.g., KeeWeb).

WARNING: Do not use Google Authenticator or Authy - they generate incorrect codes!`,
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

		qr, err := qrcode.GenerateQR(account)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate QR code: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nQR code for account: %s\n\n", account.AccountName)
		fmt.Println(qr)
		fmt.Println("\n⚠️  Do not use Google Authenticator or Authy!")
		fmt.Println("Recommended: KeeWeb, 1Password, Bitwarden")
	},
}

func init() {
	rootCmd.AddCommand(qrCmd)
}

