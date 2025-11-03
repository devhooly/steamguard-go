package steamapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/devhooly/steamguard-fork/internal/manifest"
)

const steamLoginBase = "https://login.steampowered.com"


// Client represents a client for working with Steam API
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new Steam API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Confirmation represents a trade/market confirmation
type Confirmation struct {
	ID          string
	Key         string
	Description string
	Type        string
	Creator     string
	Time        int64
}

// Login performs Steam login (simplified version)
func (c *Client) Login(username, password string) (*manifest.SessionData, error) {
	// NOTE: Full Steam login implementation is very complex and requires
	// handling RSA encryption, captcha, email/mobile codes, etc.
	// This is a simplified stub for demonstration purposes
	
	return nil, fmt.Errorf("Login function requires full Steam Auth protocol implementation")
}



// GetConfirmations gets a list of pending confirmations
func (c *Client) GetConfirmations(account *manifest.SteamGuardAccount) ([]*Confirmation, error) {
	timestamp := time.Now().Unix()
	
	// Generate hashes for confirmations
	hashConf, err := account.GenerateConfirmationHash("conf", timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate conf hash: %w", err)
	}

	deviceID := account.GetDeviceID()
	steamID := account.Session.SteamID

	// Build URL
	params := url.Values{}
	params.Set("p", deviceID)
	params.Set("a", steamID)
	params.Set("k", hashConf)
	params.Set("t", fmt.Sprintf("%d", timestamp))
	params.Set("m", "android")
	params.Set("tag", "conf")

	confURL := fmt.Sprintf("%s/mobileconf/conf?%s", steamLoginBase, params.Encode())

	// Execute request
	req, err := http.NewRequest("GET", confURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add cookies
	c.addSessionCookies(req, account)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get confirmations: %d - %s", resp.StatusCode, string(body))
	}

	// Parse response (simplified version, actual response is HTML)
	// In reality, HTML page parsing is needed
	return []*Confirmation{}, nil
}

// AcceptConfirmation accepts a confirmation
func (c *Client) AcceptConfirmation(account *manifest.SteamGuardAccount, conf *Confirmation) error {
	return c.respondToConfirmation(account, conf, "allow")
}

// RejectConfirmation rejects a confirmation
func (c *Client) RejectConfirmation(account *manifest.SteamGuardAccount, conf *Confirmation) error {
	return c.respondToConfirmation(account, conf, "cancel")
}

// respondToConfirmation sends a response to a confirmation
func (c *Client) respondToConfirmation(account *manifest.SteamGuardAccount, conf *Confirmation, op string) error {
	timestamp := time.Now().Unix()
	
	// Generate hash for specific operation
	hashDetails, err := account.GenerateConfirmationHash("details"+conf.ID, timestamp)
	if err != nil {
		return fmt.Errorf("failed to generate hash: %w", err)
	}

	deviceID := account.GetDeviceID()
	steamID := account.Session.SteamID

	// Build parameters
	params := url.Values{}
	params.Set("op", op)
	params.Set("p", deviceID)
	params.Set("a", steamID)
	params.Set("k", hashDetails)
	params.Set("t", fmt.Sprintf("%d", timestamp))
	params.Set("m", "android")
	params.Set("tag", "conf")
	params.Set("cid", conf.ID)
	params.Set("ck", conf.Key)

	confURL := fmt.Sprintf("%s/mobileconf/ajaxop", steamLoginBase)

	req, err := http.NewRequest("POST", confURL, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.addSessionCookies(req, account)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to respond to confirmation: %d - %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var result struct {
		Success bool `json:"success"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("operation was not successful")
	}

	return nil
}

// addSessionCookies adds session cookies to the request
func (c *Client) addSessionCookies(req *http.Request, account *manifest.SteamGuardAccount) {
	if account.Session.SessionID != "" {
		req.AddCookie(&http.Cookie{
			Name:  "sessionid",
			Value: account.Session.SessionID,
		})
	}
	if account.Session.SteamLoginSecure != "" {
		req.AddCookie(&http.Cookie{
			Name:  "steamLoginSecure",
			Value: account.Session.SteamLoginSecure,
		})
	}
}
