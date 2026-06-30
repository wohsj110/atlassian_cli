package present

import (
	"fmt"
	"strings"

	sharedpresent "github.com/wohsj110/atlassian_cli/shared/present"

	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/api"
)

// UserPresenter creates presentation models for atk-cfl user output.
type UserPresenter struct{}

// PresentUserOneLiner builds the canonical `atk-cfl me` output as one normalized
// pipe-delimited line on stdout.
func (UserPresenter) PresentUserOneLiner(user *api.User) *sharedpresent.OutputModel {
	return &sharedpresent.OutputModel{
		Sections: []sharedpresent.Section{
			&sharedpresent.MessageSection{
				Kind:    sharedpresent.MessageInfo,
				Message: userOneLiner(user),
				Stream:  sharedpresent.StreamStdout,
			},
		},
	}
}

// PresentUserIDOnly builds the canonical `atk-cfl me --id` output as one
// normalized account ID on stdout.
func (UserPresenter) PresentUserIDOnly(user *api.User) *sharedpresent.OutputModel {
	return &sharedpresent.OutputModel{
		Sections: []sharedpresent.Section{
			&sharedpresent.MessageSection{
				Kind:    sharedpresent.MessageInfo,
				Message: normalizeMeField(user.AccountID),
				Stream:  sharedpresent.StreamStdout,
			},
		},
	}
}

func userOneLiner(user *api.User) string {
	return fmt.Sprintf(
		"%s | %s | %s",
		normalizeMeField(user.AccountID),
		normalizeMeField(user.DisplayName),
		normalizeMeField(user.Email),
	)
}

func normalizeMeField(value string) string {
	if value == "" {
		return "-"
	}
	value = strings.ReplaceAll(value, "\r\n", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "|", `\|`)
	return value
}
