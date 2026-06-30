package credstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LegacyCreds is a minimal projection of either tool's legacy config
// file. Used by init reconciliation to compare against the shared store
// without depending on each tool's full Config struct.
type LegacyCreds struct {
	Path           string // "" if file was absent
	URL            string
	Email          string
	APIToken       string
	AuthMethod     string
	CloudID        string
	DefaultSpace   string // cfl-only
	DefaultProject string // jtk-only
	OutputFormat   string // cfl-only
}

// Section returns the credential fields, with URL normalized to base.
func (l *LegacyCreds) Section() Section {
	return Section{
		URL:        NormalizeBaseURL(l.URL),
		Email:      l.Email,
		APIToken:   l.APIToken,
		AuthMethod: l.AuthMethod,
		CloudID:    l.CloudID,
	}
}

// SharedLegacyConn is one shared-store section's pre-MON-5328
// connection + token fields. Decoded ONLY by the one-time migration so
// it can still see legacy per-tool credentials before the stripped
// schema is written. The canonical Store no longer exposes these — do
// NOT use for runtime resolution.
type SharedLegacyConn struct {
	URL        string `yaml:"url"`
	Email      string `yaml:"email"`
	APIToken   string `yaml:"api_token"`
	AuthMethod string `yaml:"auth_method"`
	CloudID    string `yaml:"cloud_id"`
}

// SharedLegacyProjection is the migration-only decode of the shared
// config.yml retaining the per-tool connection/token fields the
// canonical Store dropped (§2.2 / MON-5328). Migration-only.
type SharedLegacyProjection struct {
	Path    string
	Default SharedLegacyConn
	AtkCFL  SharedLegacyConn
	AtkJira SharedLegacyConn
}

// LoadSharedLegacyProjection decodes path retaining legacy per-tool
// connection/token fields. Absent file → (nil, nil). Parse failure →
// ErrCorruptStore (same contract as Load) so callers refuse to clobber
// an unreadable file.
func LoadSharedLegacyProjection(path string) (*SharedLegacyProjection, error) {
	data, err := os.ReadFile(path) //nolint:gosec // CLI tool reading its own config
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: reading %s: %s", ErrCorruptStore, path, err.Error())
	}
	var raw struct {
		Default   SharedLegacyConn `yaml:"default"`
		AtkCFL    SharedLegacyConn `yaml:"atk_cfl"`
		AtkJira   SharedLegacyConn `yaml:"atk_jira"`
		LegacyCFL SharedLegacyConn `yaml:"cfl"`
		LegacyJTK SharedLegacyConn `yaml:"jtk"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%w: parsing %s: %s", ErrCorruptStore, path, err.Error())
	}
	cfl := raw.AtkCFL
	if !sharedLegacyConnHasField(cfl) {
		cfl = raw.LegacyCFL
	}
	jira := raw.AtkJira
	if !sharedLegacyConnHasField(jira) {
		jira = raw.LegacyJTK
	}
	return &SharedLegacyProjection{Path: path, Default: raw.Default, AtkCFL: cfl, AtkJira: jira}, nil
}

func sharedLegacyConnHasField(c SharedLegacyConn) bool {
	return c.URL != "" || c.Email != "" || c.APIToken != "" ||
		c.AuthMethod != "" || c.CloudID != ""
}

// LegacyAtkCFLPath returns the canonical cfl legacy config path.
func LegacyAtkCFLPath() string {
	return tooledPath("cfl", "config.yml")
}

// LegacyAtkJiraPath returns the canonical jtk legacy config path. jtk's
// loader uses os.UserConfigDir(), which on macOS is
// ~/Library/Application Support — matching it here is critical so
// macOS users with an existing jtk config are detected by sibling
// init reconciliation.
func LegacyAtkJiraPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join(".", "atk-jira", "config.json")
	}
	return filepath.Join(dir, "atk-jira", "config.json")
}

// LegacyAgentToolPaths returns this fork's pre-shared-config JSON paths.
// They are read only by migration/init reconciliation. Runtime config
// resolution does not consult them.
func LegacyAgentToolPaths(tool string) []string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	base := filepath.Join(dir, "atlassian-agent-cli")
	switch tool {
	case ToolAtkCFL:
		return []string{
			filepath.Join(base, "atk-cfl.json"),
			filepath.Join(base, "atk-confluence.json"),
		}
	case ToolAtkJira:
		return []string{filepath.Join(base, "atk-jira.json")}
	default:
		return nil
	}
}

func tooledPath(toolDir, file string) string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, toolDir, file)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "."+toolDir, file)
	}
	return filepath.Join(home, ".config", toolDir, file)
}

// LoadLegacyAtkCFL reads a cfl YAML legacy config file. An absent file
// returns (nil, nil). Parse failures return ErrCorruptStore so the
// caller can refuse to clobber it.
func LoadLegacyAtkCFL(path string) (*LegacyCreds, error) {
	data, err := os.ReadFile(path) //nolint:gosec // CLI tool reading its own config
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: reading %s: %s", ErrCorruptStore, path, err.Error())
	}
	var raw struct {
		URL          string `yaml:"url"`
		Email        string `yaml:"email"`
		APIToken     string `yaml:"api_token"`
		DefaultSpace string `yaml:"default_space"`
		OutputFormat string `yaml:"output_format"`
		AuthMethod   string `yaml:"auth_method"`
		CloudID      string `yaml:"cloud_id"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%w: parsing %s: %s", ErrCorruptStore, path, err.Error())
	}
	return &LegacyCreds{
		Path:         path,
		URL:          raw.URL,
		Email:        raw.Email,
		APIToken:     raw.APIToken,
		AuthMethod:   raw.AuthMethod,
		CloudID:      raw.CloudID,
		DefaultSpace: raw.DefaultSpace,
		OutputFormat: raw.OutputFormat,
	}, nil
}

// LoadLegacyAtkJira reads a jtk JSON legacy config file. An absent file
// returns (nil, nil). Parse failures return ErrCorruptStore.
func LoadLegacyAtkJira(path string) (*LegacyCreds, error) {
	data, err := os.ReadFile(path) //nolint:gosec // CLI tool reading its own config
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: reading %s: %s", ErrCorruptStore, path, err.Error())
	}
	var raw struct {
		URL            string `json:"url"`
		Domain         string `json:"domain"`
		Email          string `json:"email"`
		APIToken       string `json:"api_token"`
		DefaultProject string `json:"default_project"`
		AuthMethod     string `json:"auth_method"`
		CloudID        string `json:"cloud_id"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%w: parsing %s: %s", ErrCorruptStore, path, err.Error())
	}
	url := raw.URL
	if url == "" && raw.Domain != "" {
		url = "https://" + raw.Domain + ".atlassian.net"
	}
	return &LegacyCreds{
		Path:           path,
		URL:            url,
		Email:          raw.Email,
		APIToken:       raw.APIToken,
		AuthMethod:     raw.AuthMethod,
		CloudID:        raw.CloudID,
		DefaultProject: raw.DefaultProject,
	}, nil
}

// LoadLegacyAgentTool reads this fork's early atlassian-agent-cli JSON
// config shape. It accepts "site"/"token"/"auth_type" and maps them to
// the shared migration projection.
func LoadLegacyAgentTool(path, tool string) (*LegacyCreds, error) {
	data, err := os.ReadFile(path) //nolint:gosec // CLI migration reads its own legacy config
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: reading %s: %s", ErrCorruptStore, path, err.Error())
	}
	var raw struct {
		Site           string `json:"site"`
		URL            string `json:"url"`
		Email          string `json:"email"`
		Token          string `json:"token"`
		APIToken       string `json:"api_token"`
		AuthType       string `json:"auth_type"`
		AuthMethod     string `json:"auth_method"`
		CloudID        string `json:"cloud_id"`
		DefaultSpace   string `json:"default_space"`
		DefaultProject string `json:"default_project"`
		OutputFormat   string `json:"output_format"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%w: parsing %s: %s", ErrCorruptStore, path, err.Error())
	}
	url := raw.URL
	if url == "" {
		url = raw.Site
	}
	token := raw.APIToken
	if token == "" {
		token = raw.Token
	}
	authMethod := raw.AuthMethod
	if authMethod == "" {
		authMethod = raw.AuthType
	}
	creds := &LegacyCreds{
		Path:           path,
		URL:            url,
		Email:          raw.Email,
		APIToken:       token,
		AuthMethod:     authMethod,
		CloudID:        raw.CloudID,
		DefaultSpace:   raw.DefaultSpace,
		DefaultProject: raw.DefaultProject,
		OutputFormat:   raw.OutputFormat,
	}
	if tool == ToolAtkCFL && creds.DefaultSpace == "" {
		creds.DefaultSpace = raw.DefaultProject
	}
	return creds, nil
}
