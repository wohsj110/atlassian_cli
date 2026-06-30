// Package config manages atk-jira configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/wohsj110/atlassian_cli/shared/auth"
	"github.com/wohsj110/atlassian_cli/shared/credstore"
	"github.com/wohsj110/atlassian_cli/shared/keyring"
	"github.com/wohsj110/atlassian_cli/shared/url"
)

// loadShared returns the shared credential store. Accessors can't
// propagate errors, so on corrupt shared store we warn once on stderr
// (so the user sees something is wrong) and fall through to legacy
// reads. Init has a separate code path that surfaces corruption as a
// hard error and refuses to clobber the file.
func loadShared() *credstore.Store {
	// §3.2 relocation-aware, mutation-free runtime resolver: canonical
	// store, transparent read-fallback to the prior hand-rolled location
	// when only it exists, and on an old↔new divergence the canonical
	// store is returned alongside the error so commands keep working
	// while the conflict is surfaced once. `atk-jira init` is the fail-loud
	// mutating gate.
	s, err := credstore.LoadSharedRuntime()
	if err != nil {
		warnCorruptSharedOnce(err)
		if s == nil {
			return &credstore.Store{}
		}
	}
	return s
}

var corruptSharedWarnOnce sync.Once

func warnCorruptSharedOnce(err error) {
	corruptSharedWarnOnce.Do(func() {
		if errors.Is(err, credstore.ErrRelocationConflict) {
			// Readable, not a fallback: the canonical config is in use.
			fmt.Fprintf(os.Stderr, "warning: prior and current shared config diverge (%v); using the current config. Run `atk-jira init` to reconcile.\n", err)
			return
		}
		fmt.Fprintf(os.Stderr, "warning: shared credential store is unreadable (%v); falling back to per-tool config. Run `atk-jira init` to fix.\n", err)
	})
}

// jiraSection returns the resolved connection Section from the shared
// `default` (§2.2: single-sourced — no per-tool override).
func jiraSection() credstore.Section {
	return loadShared().Resolve(credstore.ToolAtkJira)
}

// jiraSectionWithSource returns the resolved value and source for one
// field of the atk_jira section.
func jiraSectionWithSource(field string) (string, credstore.Source) {
	return loadShared().ResolveWithSource(credstore.ToolAtkJira, field)
}

const (
	configFileMode = 0600
	configDirMode  = 0700
)

// Config holds the CLI configuration
type Config struct {
	URL            string `json:"url,omitempty"`
	Domain         string `json:"domain,omitempty"` // Deprecated: use URL instead
	Email          string `json:"email"`
	APIToken       string `json:"api_token"`
	DefaultProject string `json:"default_project,omitempty"`
	AuthMethod     string `json:"auth_method,omitempty"` // "basic" (default) or "bearer"
	CloudID        string `json:"cloud_id,omitempty"`    // Required for bearer auth (gateway URL)
	// Keyring's `omitempty` is a no-op on this struct type — encoding/json
	// emits an empty {} for zero-valued struct fields regardless. Kept for
	// stylistic consistency with the other optional fields; the empty
	// section is harmless on read.
	Keyring KeyringConfig `json:"keyring,omitempty"`
}

// KeyringConfig holds keyring-related user preferences.
type KeyringConfig struct {
	// Backend, when set, requests a specific credstore backend at runtime.
	// Lower precedence than --backend and ATLASSIAN_AGENT_CLI_KEYRING_BACKEND.
	// Valid values: see credstore.ValidBackendNames(). Validation happens
	// inside credstore.Open at startup; an unrecognized value fails closed
	// with an error wrapping ErrBackendNotImplemented.
	Backend string `json:"backend,omitempty"`
}

// configPath returns the path to the canonical shared config file.
func configPath() (string, error) {
	return credstore.DefaultPath()
}

// Load loads Jira-relevant configuration from the canonical shared store.
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	store, err := credstore.Load(path)
	if err != nil {
		return nil, err
	}
	section := store.Resolve(credstore.ToolAtkJira)
	return &Config{
		URL:            section.URL,
		Email:          section.Email,
		APIToken:       section.APIToken,
		AuthMethod:     section.AuthMethod,
		CloudID:        section.CloudID,
		DefaultProject: store.AtkJira.DefaultProject,
	}, nil
}

// Save saves Jira-relevant non-secret configuration to the shared store.
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	store, err := credstore.Load(path)
	if err != nil {
		return err
	}
	store.Default = credstore.Section{
		URL:        credstore.NormalizeBaseURL(cfg.URL),
		Email:      cfg.Email,
		AuthMethod: cfg.AuthMethod,
		CloudID:    cfg.CloudID,
	}
	store.AtkJira.DefaultProject = cfg.DefaultProject
	return store.Save(path)
}

// Clear removes the configuration file
func Clear() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}

// GetURL returns the Jira URL from config or environment.
// Precedence: JIRA_URL → ATLASSIAN_URL → shared default → JIRA_DOMAIN.
func GetURL() string {
	if v := os.Getenv("JIRA_URL"); v != "" {
		return url.NormalizeURL(v)
	}
	if v := os.Getenv("ATLASSIAN_URL"); v != "" {
		return url.NormalizeURL(v)
	}
	if v := jiraSection().URL; v != "" {
		return url.NormalizeURL(v)
	}
	if v := os.Getenv("JIRA_DOMAIN"); v != "" {
		return "https://" + v + ".atlassian.net"
	}
	return ""
}

// GetDomain returns the domain from environment only.
// Deprecated: Use GetURL instead.
func GetDomain() string {
	if v := os.Getenv("JIRA_DOMAIN"); v != "" {
		return v
	}
	return ""
}

// GetEmail returns the email from config or environment.
// Precedence: JIRA_EMAIL → ATLASSIAN_EMAIL → shared default.
func GetEmail() string {
	if v := os.Getenv("JIRA_EMAIL"); v != "" {
		return v
	}
	if v := os.Getenv("ATLASSIAN_EMAIL"); v != "" {
		return v
	}
	if v := jiraSection().Email; v != "" {
		return v
	}
	return ""
}

// ResolveAPIToken is the AUTHORITATIVE runtime token resolver: env
// (JIRA_API_TOKEN → ATLASSIAN_API_TOKEN) then the OS keyring, running the
// one-time §1.8 migration. A keyring error PROPAGATES — it must never be
// folded into an empty token (that would silently de-authenticate every
// command). This is the single migrating entry point; APIClient uses it.
func ResolveAPIToken() (string, error) {
	tok, _, err := keyring.ResolveToken(credstore.ToolAtkJira)
	return tok, err
}

// GetAPIToken returns the API token via the NON-migrating keyring path
// (env → keyring), swallowing keyring errors to an empty string. It is
// used only by diagnostics (`config show` source column) and the
// IsConfigured gate; the authoritative, error-propagating path is
// ResolveAPIToken. The token is no longer read from the plaintext config
// file or shared store.
func GetAPIToken() string {
	tok, _, err := keyring.ResolveTokenNoMigrate(credstore.ToolAtkJira)
	if err != nil {
		return ""
	}
	return tok
}

// IsConfigured returns true if the NON-SECRET config is complete and a
// token is resolvable (env or keyring, non-migrating). The token left
// the plaintext config store, so completeness is composed from both
// halves. For bearer auth: URL + Cloud ID + token; for basic: URL +
// email + token.
func IsConfigured() bool {
	if GetAuthMethod() == auth.AuthMethodBearer {
		return GetURL() != "" && GetCloudID() != "" && GetAPIToken() != ""
	}
	return GetURL() != "" && GetEmail() != "" && GetAPIToken() != ""
}

// GetAuthMethod returns the auth method from config or environment.
// Precedence: JIRA_AUTH_METHOD → ATLASSIAN_AUTH_METHOD → shared default → "basic"
// Invalid values are ignored and fall through to the next source.
func GetAuthMethod() string {
	v, _ := GetAuthMethodWithSource()
	return v
}

// GetAuthMethodWithSource returns the auth method and its source.
// Precedence: JIRA_AUTH_METHOD → ATLASSIAN_AUTH_METHOD → shared default → "basic"
// Invalid values are skipped and fall through to the next source.
// Validation happens at entry points (api.New, init --auth-method) not here.
func GetAuthMethodWithSource() (value, source string) {
	if v := os.Getenv("JIRA_AUTH_METHOD"); v != "" {
		if auth.ValidateAuthMethod(v) == nil {
			return v, "env (JIRA_AUTH_METHOD)"
		}
	}
	if v := os.Getenv("ATLASSIAN_AUTH_METHOD"); v != "" {
		if auth.ValidateAuthMethod(v) == nil {
			return v, "env (ATLASSIAN_AUTH_METHOD)"
		}
	}
	if v, src := jiraSectionWithSource("auth_method"); v != "" && auth.ValidateAuthMethod(v) == nil {
		return v, string(src)
	}
	return auth.AuthMethodBasic, "default"
}

// GetCloudID returns the Atlassian Cloud ID from config or environment.
// Precedence: JIRA_CLOUD_ID → ATLASSIAN_CLOUD_ID → shared default.
func GetCloudID() string {
	v, _ := GetCloudIDWithSource()
	return v
}

// GetCloudIDWithSource returns the Cloud ID and its source.
// Precedence: JIRA_CLOUD_ID → ATLASSIAN_CLOUD_ID → shared default.
func GetCloudIDWithSource() (value, source string) {
	if v := os.Getenv("JIRA_CLOUD_ID"); v != "" {
		return v, "env (JIRA_CLOUD_ID)"
	}
	if v := os.Getenv("ATLASSIAN_CLOUD_ID"); v != "" {
		return v, "env (ATLASSIAN_CLOUD_ID)"
	}
	if v, src := jiraSectionWithSource("cloud_id"); v != "" {
		return v, string(src)
	}
	return "", "-"
}

// GetDefaultProject returns the default project from config or environment.
// Precedence: JIRA_DEFAULT_PROJECT → shared atk_jira.default_project.
func GetDefaultProject() string {
	if v := os.Getenv("JIRA_DEFAULT_PROJECT"); v != "" {
		return v
	}
	if v := loadShared().AtkJira.DefaultProject; v != "" {
		return v
	}
	return ""
}

// Path returns the path to the config file
func Path() string {
	path, _ := configPath()
	return path
}
