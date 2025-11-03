package qrcode

import (
	"fmt"
	"net/url"

	"github.com/devhooly/steamguard-fork/internal/manifest"
	qr "github.com/skip2/go-qrcode"
)

// GenerateQR generates a QR code for importing into other 2FA applications
func GenerateQR(account *manifest.SteamGuardAccount) (string, error) {
	// Create otpauth URL
	uri := generateOTPAuthURL(account)

	// Generate QR code as ASCII art
	qrCode, err := qr.New(uri, qr.Medium)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Return as string
	return qrCode.ToSmallString(false), nil
}

// GenerateQRPNG generates a QR code as a PNG image
func GenerateQRPNG(account *manifest.SteamGuardAccount, size int) ([]byte, error) {
	uri := generateOTPAuthURL(account)
	
	png, err := qr.Encode(uri, qr.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return png, nil
}

// generateOTPAuthURL creates an otpauth:// URL for 2FA
func generateOTPAuthURL(account *manifest.SteamGuardAccount) string {
	// Format: otpauth://totp/Steam:username?secret=SECRET&issuer=Steam&algorithm=SHA1&digits=5&period=30
	
	params := url.Values{}
	params.Set("secret", account.SharedSecret)
	params.Set("issuer", "Steam")
	params.Set("algorithm", "SHA1")
	params.Set("digits", "5")
	params.Set("period", "30")

	label := url.PathEscape(fmt.Sprintf("Steam:%s", account.AccountName))
	
	return fmt.Sprintf("otpauth://totp/%s?%s", label, params.Encode())
}

// GetOTPAuthURL returns the otpauth:// URL for manual entry
func GetOTPAuthURL(account *manifest.SteamGuardAccount) string {
	return generateOTPAuthURL(account)
}
