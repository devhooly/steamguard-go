package manifest

import (
	"fmt"

	"github.com/devhooly/steamguard-fork/internal/steamguard"
)

// SteamGuardAccount represents a Steam Guard account (SDA format)
type SteamGuardAccount struct {
	AccountName        string        `json:"account_name"`
	SharedSecret       string        `json:"shared_secret"`
	SerialNumber       string        `json:"serial_number"`
	RevocationCode     string        `json:"revocation_code"`
	URI                string        `json:"uri"`
	ServerTime         int64         `json:"server_time"`
	TokenGID           string        `json:"token_gid"`
	IdentitySecret     string        `json:"identity_secret"`
	Secret1            string        `json:"secret_1"`
	Status             int           `json:"status"`
	DeviceID           string        `json:"device_id"`
	FullyEnrolled      bool          `json:"fully_enrolled"`
	Session            SessionData   `json:"Session"`
}

// SessionData represents Steam session data
type SessionData struct {
	SessionID         string `json:"SessionID"`
	SteamLogin        string `json:"SteamLogin"`
	SteamLoginSecure  string `json:"SteamLoginSecure"`
	WebCookie         string `json:"WebCookie"`
	OAuthToken        string `json:"OAuthToken"`
	SteamID           string `json:"SteamID"`
}

// GenerateCode generates the current Steam Guard code
func (a *SteamGuardAccount) GenerateCode() (string, error) {
	if a.SharedSecret == "" {
		return "", fmt.Errorf("shared secret is missing")
	}

	code, err := steamguard.GenerateSteamGuardCode(a.SharedSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}

	return code, nil
}

// GenerateConfirmationHash generates a hash for confirmation
func (a *SteamGuardAccount) GenerateConfirmationHash(tag string, timestamp int64) (string, error) {
	if a.IdentitySecret == "" {
		return "", fmt.Errorf("identity secret is missing")
	}

	hash, err := steamguard.GenerateConfirmationHash(a.IdentitySecret, tag, timestamp)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash: %w", err)
	}

	return hash, nil
}

// GetDeviceID returns the Device ID (generates if needed)
func (a *SteamGuardAccount) GetDeviceID() string {
	if a.DeviceID == "" {
		a.DeviceID = steamguard.GenerateDeviceID(a.Session.SteamID)
	}
	return a.DeviceID
}
