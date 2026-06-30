package config

import (
	"os"

	"github.com/wohsj110/atlassian_cli/shared/credstore"
	"github.com/wohsj110/atlassian_cli/shared/keyring"
)

// ValuesWithSources holds all config values with their source information.
// This is a projection helper that inspects env vars and config file to determine
// where each value came from. Used by commands to pass resolved values to presenters.
type ValuesWithSources struct {
	URL         string
	URLSource   string
	Email       string
	EmailSource string
	// The API token VALUE is never projected — the keyring is the source
	// of truth and §1.12 forbids displaying it (or any prefix/suffix).
	// Only presence + source + non-secret keyring metadata are shown.
	TokenConfigured   bool
	TokenSource       string
	KeyringRef        string
	KeyringBackend    string
	KeyringPassphrase string // file backend only; "" otherwise
	DefaultProject    string
	ProjectSource     string
	AuthMethod        string
	AuthMethodSrc     string
	CloudID           string
	CloudIDSrc        string
	Path              string
}

// GetValuesWithSources returns all config values with their source information.
func GetValuesWithSources() ValuesWithSources {
	url, urlSrc := GetURLWithSource()
	email, emailSrc := GetEmailWithSource()
	project, projectSrc := GetDefaultProjectWithSource()
	authMethod, authMethodSrc := GetAuthMethodWithSource()
	cloudID, cloudIDSrc := GetCloudIDWithSource()

	// Non-migrating: `config show` is diagnostic and must stay usable
	// even during an unresolved §1.8 conflict. A keyring error is folded
	// into a clear source label rather than crashing show.
	kr, err := keyring.InspectForTool(credstore.ToolAtkJira)
	if err != nil {
		kr.TokenSource = "keyring error: " + err.Error()
	}

	return ValuesWithSources{
		URL:               url,
		URLSource:         urlSrc,
		Email:             email,
		EmailSource:       emailSrc,
		TokenConfigured:   kr.TokenConfigured,
		TokenSource:       kr.TokenSource,
		KeyringRef:        kr.Ref,
		KeyringBackend:    keyringBackendLabel(kr),
		KeyringPassphrase: kr.PassphraseSource,
		DefaultProject:    project,
		ProjectSource:     projectSrc,
		AuthMethod:        authMethod,
		AuthMethodSrc:     authMethodSrc,
		CloudID:           cloudID,
		CloudIDSrc:        cloudIDSrc,
		Path:              Path(),
	}
}

// keyringBackendLabel renders the backend and how it was selected, e.g.
// "keychain (auto)". Empty when the keyring could not be opened.
func keyringBackendLabel(kr keyring.Info) string {
	if kr.Backend == "" {
		return ""
	}
	if kr.BackendSource == "" {
		return kr.Backend
	}
	return kr.Backend + " (" + kr.BackendSource + ")"
}

// GetURLWithSource returns the URL and its source.
// Precedence: JIRA_URL → ATLASSIAN_URL → shared default → JIRA_DOMAIN.
func GetURLWithSource() (value, source string) {
	if os.Getenv("JIRA_URL") != "" {
		return GetURL(), "env (JIRA_URL)"
	}
	if os.Getenv("ATLASSIAN_URL") != "" {
		return GetURL(), "env (ATLASSIAN_URL)"
	}
	if v, src := jiraSectionWithSource("url"); v != "" {
		return GetURL(), string(src)
	}
	if os.Getenv("JIRA_DOMAIN") != "" {
		return GetURL(), "env (JIRA_DOMAIN, deprecated)"
	}
	return "", "-"
}

// GetEmailWithSource returns the email and its source.
// Precedence: JIRA_EMAIL → ATLASSIAN_EMAIL → shared default.
func GetEmailWithSource() (value, source string) {
	if os.Getenv("JIRA_EMAIL") != "" {
		return GetEmail(), "env (JIRA_EMAIL)"
	}
	if os.Getenv("ATLASSIAN_EMAIL") != "" {
		return GetEmail(), "env (ATLASSIAN_EMAIL)"
	}
	if v, src := jiraSectionWithSource("email"); v != "" {
		return v, string(src)
	}
	return "", "-"
}

// GetDefaultProjectWithSource returns the default project and its source.
// Precedence: JIRA_DEFAULT_PROJECT → shared atk_jira.default_project.
func GetDefaultProjectWithSource() (value, source string) {
	if os.Getenv("JIRA_DEFAULT_PROJECT") != "" {
		return GetDefaultProject(), "env (JIRA_DEFAULT_PROJECT)"
	}
	if v := loadShared().AtkJira.DefaultProject; v != "" {
		return v, "shared atk_jira"
	}
	return "", "-"
}
