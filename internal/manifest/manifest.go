package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/devhooly/steamguard-go/internal/crypto"
)

// Manifest represents the manifest.json file
type Manifest struct {
	Encrypted     bool              `json:"encrypted"`
	FirstRun      bool              `json:"first_run"`
	Entries       []ManifestEntry   `json:"entries"`
	PeriodicCheck bool              `json:"periodic_checking"`
	PeriodicTime  int               `json:"periodic_checking_interval"`
	AutoConfirm   []AutoConfirmRule `json:"auto_confirm_trades"`
}

// ManifestEntry represents an entry in the manifest for one account
type ManifestEntry struct {
	Encryption *EncryptionParams `json:"encryption,omitempty"`
	Filename   string            `json:"filename"`
	SteamID    string            `json:"steamid"`
}

// EncryptionParams encryption parameters for an account
type EncryptionParams struct {
	IV   string `json:"iv"`
	Salt string `json:"salt"`
}

// AutoConfirmRule rule for automatic confirmation
type AutoConfirmRule struct {
	Action string `json:"action"`
}

// Manager manages the manifest and accounts
type Manager struct {
	mu          sync.RWMutex
	manifest    *Manifest
	accounts    map[string]*SteamGuardAccount
	maFilesPath string
	passkey     string
}

// NewManager creates a new manifest manager
func NewManager(maFilesPath string) (*Manager, error) {
	mgr := &Manager{
		maFilesPath: maFilesPath,
		accounts:    make(map[string]*SteamGuardAccount),
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(maFilesPath, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	manifestPath := filepath.Join(maFilesPath, "manifest.json")
	
	// Check if manifest exists
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// Create new manifest
		mgr.manifest = &Manifest{
			Encrypted:  false,
			FirstRun:   true,
			Entries:    []ManifestEntry{},
			AutoConfirm: []AutoConfirmRule{},
		}
		return mgr, mgr.Save()
	}

	// Load existing manifest
	if err := mgr.Load(); err != nil {
		return nil, err
	}

	return mgr, nil
}

// Load loads the manifest and all accounts
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	manifestPath := filepath.Join(m.maFilesPath, "manifest.json")
	
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	m.manifest = &Manifest{}
	if err := json.Unmarshal(data, m.manifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Load all accounts
	for _, entry := range m.manifest.Entries {
		account, err := m.loadAccount(entry)
		if err != nil {
			return fmt.Errorf("failed to load account %s: %w", entry.SteamID, err)
		}
		m.accounts[account.AccountName] = account
	}

	return nil
}

// loadAccount loads one account from file
func (m *Manager) loadAccount(entry ManifestEntry) (*SteamGuardAccount, error) {
	accountPath := filepath.Join(m.maFilesPath, entry.Filename)
	
	data, err := os.ReadFile(accountPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read account file: %w", err)
	}

	// If data is encrypted, decrypt it
	if m.manifest.Encrypted && entry.Encryption != nil {
		if m.passkey == "" {
			// Ask for password
			fmt.Print("Enter password to decrypt: ")
			fmt.Scanln(&m.passkey)
		}

		decrypted, err := crypto.Decrypt(data, m.passkey, entry.Encryption.IV, entry.Encryption.Salt)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt: %w", err)
		}
		data = decrypted
	}

	account := &SteamGuardAccount{}
	if err := json.Unmarshal(data, account); err != nil {
		return nil, fmt.Errorf("failed to parse account: %w", err)
	}

	return account, nil
}

// Save saves the manifest
func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	manifestPath := filepath.Join(m.maFilesPath, "manifest.json")
	
	data, err := json.MarshalIndent(m.manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// AddAccount adds a new account
func (m *Manager) AddAccount(account *SteamGuardAccount) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	filename := fmt.Sprintf("%s.maFile", account.AccountName)
	accountPath := filepath.Join(m.maFilesPath, filename)

	// Serialize account
	data, err := json.MarshalIndent(account, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize account: %w", err)
	}

	// If encryption is needed
	var encParams *EncryptionParams
	if m.manifest.Encrypted {
		if m.passkey == "" {
			fmt.Print("Enter password for encryption: ")
			fmt.Scanln(&m.passkey)
		}

		encrypted, iv, salt, err := crypto.Encrypt(data, m.passkey)
		if err != nil {
			return fmt.Errorf("failed to encrypt: %w", err)
		}
		data = encrypted
		encParams = &EncryptionParams{
			IV:   iv,
			Salt: salt,
		}
	}

	// Write file
	if err := os.WriteFile(accountPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write account file: %w", err)
	}

	// Add to manifest
	entry := ManifestEntry{
		Encryption: encParams,
		Filename:   filename,
		SteamID:    account.Session.SteamID,
	}
	m.manifest.Entries = append(m.manifest.Entries, entry)
	m.accounts[account.AccountName] = account

	// Save manifest
	return m.saveUnlocked()
}

// saveUnlocked saves the manifest without locking (for internal use)
func (m *Manager) saveUnlocked() error {
	manifestPath := filepath.Join(m.maFilesPath, "manifest.json")
	
	data, err := json.MarshalIndent(m.manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// GetAccount returns an account by name (or the first one if name is not specified)
func (m *Manager) GetAccount(username string) (*SteamGuardAccount, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if username == "" {
		// Return the first account
		if len(m.accounts) == 0 {
			return nil, fmt.Errorf("no accounts")
		}
		for _, account := range m.accounts {
			return account, nil
		}
	}

	account, ok := m.accounts[username]
	if !ok {
		return nil, fmt.Errorf("account %s not found", username)
	}

	return account, nil
}

// GetAllAccounts returns all accounts
func (m *Manager) GetAllAccounts() []*SteamGuardAccount {
	m.mu.RLock()
	defer m.mu.RUnlock()

	accounts := make([]*SteamGuardAccount, 0, len(m.accounts))
	for _, account := range m.accounts {
		accounts = append(accounts, account)
	}

	return accounts
}

// IsEmpty checks if there are any accounts
func (m *Manager) IsEmpty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.accounts) == 0
}
