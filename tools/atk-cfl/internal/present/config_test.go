package present

import (
	"errors"
	"testing"

	"github.com/wohsj110/atlassian_cli/shared/keyring"
	sharedpresent "github.com/wohsj110/atlassian_cli/shared/present"
	"github.com/wohsj110/atlassian_cli/shared/testutil"

	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/api"
)

func TestConfigPresenter_PresentTestSuccess(t *testing.T) {
	t.Parallel()

	progressModel := ConfigPresenter{}.PresentTestProgress()
	progressMsg := requireMessageSection(t, progressModel, 0)
	testutil.Equal(t, sharedpresent.StreamStderr, progressMsg.Stream)
	testutil.Equal(t, true, progressMsg.NoNewline)
	testutil.Equal(t, "Testing connection... ", progressMsg.Message)

	connectionModel := ConfigPresenter{}.PresentTestConnectionSuccess()
	connectionMsg := requireMessageSection(t, connectionModel, 0)
	testutil.Equal(t, sharedpresent.StreamStderr, connectionMsg.Stream)
	testutil.Equal(t, "success!\n", connectionMsg.Message)

	model := ConfigPresenter{}.PresentTestSuccess(&api.User{
		AccountID:   "acct-1",
		DisplayName: "Test User",
		Email:       "test@example.com",
	})

	msg := requireMessageSection(t, model, 0)
	testutil.Equal(t, sharedpresent.StreamStderr, msg.Stream)
	testutil.Equal(t, "Authentication successful\nAPI access verified\n\nAuthenticated as: Test User (test@example.com)\nAccount ID: acct-1", msg.Message)
}

func TestConfigPresenter_PresentTestSuccessFallback(t *testing.T) {
	t.Parallel()

	model := ConfigPresenter{}.PresentTestSuccess(nil)

	msg := requireMessageSection(t, model, 0)
	testutil.Equal(t, sharedpresent.StreamStderr, msg.Stream)
	testutil.Equal(t, "Your atk-cfl configuration is working correctly.", msg.Message)
}

func TestConfigPresenter_PresentTestFailure(t *testing.T) {
	t.Parallel()

	model := ConfigPresenter{}.PresentTestFailure()

	msg := requireMessageSection(t, model, 0)
	testutil.Equal(t, sharedpresent.StreamStderr, msg.Stream)
	testutil.Equal(t, "failed!\n\nTroubleshooting:\n  - Verify your URL is correct (should include https://)\n  - Check your email and API token\n  - Ensure your API token hasn't expired\n  - Verify you have permission to access Confluence\n\nTo regenerate an API token:\n  https://id.atlassian.com/manage-profile/security/api-tokens", msg.Message)
}

func TestConfigPresenter_PresentClearDefaultPlan(t *testing.T) {
	t.Parallel()

	model := ConfigPresenter{}.PresentClearDefaultPlan(keyring.ClearPlan{Ref: "atlassian-agent-cli", ToolKey: "api_token"})

	msg := requireMessageSection(t, model, 0)
	testutil.Equal(t, sharedpresent.StreamStderr, msg.Stream)
	testutil.Equal(t, "This will delete key \"api_token\" from keyring atlassian-agent-cli.\nWarning: this is the SHARED token (api_token). atk-jira will also lose access (atk-cfl and atk-jira resolve the same key).", msg.Message)
}

func TestConfigPresenter_PresentClearAllPlan(t *testing.T) {
	t.Parallel()

	model := ConfigPresenter{}.PresentClearAllPlan(keyring.ClearPlan{
		Ref:                 "atlassian-agent-cli",
		ExistingKeys:        []string{"api_token", "other"},
		SharedConfigPath:    "/tmp/shared.yml",
		OldSharedConfigPath: "/tmp/old-shared.yml",
		LegacyPaths:         []string{"/tmp/.atk-cfl.yml"},
	}, errors.New("locked"))

	msg := requireMessageSection(t, model, 0)
	testutil.Equal(t, sharedpresent.StreamStderr, msg.Stream)
	testutil.Equal(t, "This will remove the ENTIRE shared keyring bundle atlassian-agent-cli (keys: api_token, other).\nIt will also delete the shared config file: /tmp/shared.yml\nIt will also delete the prior shared config file: /tmp/old-shared.yml\nIt will scrub the legacy plaintext file: /tmp/.atk-cfl.yml\nNote: the keyring could not be opened (locked); plaintext artifacts will still be cleaned, but the keyring bundle will be left intact.", msg.Message)
}

func TestConfigPresenter_PresentClearNoStoredToken(t *testing.T) {
	t.Parallel()

	model := ConfigPresenter{}.PresentClearNoStoredToken(keyring.ClearPlan{Ref: "atlassian-agent-cli", EnvActive: []string{"CFL_API_TOKEN"}})

	msg := requireMessageSection(t, model, 0)
	testutil.Equal(t, sharedpresent.StreamStderr, msg.Stream)
	testutil.Equal(t, "No stored API token in keyring atlassian-agent-cli for atk-cfl; nothing to clear.\nNote: CFL_API_TOKEN still set in the environment and will continue to override at runtime (not cleared).", msg.Message)
}

func TestConfigPresenter_PresentClearCancelled(t *testing.T) {
	t.Parallel()

	model := ConfigPresenter{}.PresentClearCancelled()

	msg := requireMessageSection(t, model, 0)
	testutil.Equal(t, sharedpresent.StreamStderr, msg.Stream)
	testutil.Equal(t, "Cancelled. Nothing was cleared.", msg.Message)
}

func TestConfigPresenter_PresentClearDefaultSuccess(t *testing.T) {
	t.Parallel()

	model := ConfigPresenter{}.PresentClearDefaultSuccess(keyring.ClearPlan{Ref: "atlassian-agent-cli", ToolKey: "api_token"})

	msg := requireMessageSection(t, model, 0)
	testutil.Equal(t, sharedpresent.StreamStderr, msg.Stream)
	testutil.Equal(t, "Removed key \"api_token\" from keyring atlassian-agent-cli.", msg.Message)
}

func TestConfigPresenter_PresentClearAllSuccess(t *testing.T) {
	t.Parallel()

	model := ConfigPresenter{}.PresentClearAllSuccess(keyring.ClearPlan{EnvActive: []string{"ATLASSIAN_API_TOKEN"}})

	msg := requireMessageSection(t, model, 0)
	testutil.Equal(t, sharedpresent.StreamStderr, msg.Stream)
	testutil.Equal(t, "Removed the shared keyring bundle and config file.\nNote: ATLASSIAN_API_TOKEN still set in the environment and will continue to override at runtime (not cleared).", msg.Message)
}
