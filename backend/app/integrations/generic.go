package integrations

import (
	"codejam.io/integrations/discord"
	"codejam.io/integrations/github"
	"codejam.io/logging"
	"encoding/json"
	"strings"
)

var logger = logging.NewLogger(logging.Options{Name: "Integrations", Level: logging.DEBUG})

type IntegrationUser struct {
	IntegrationName string
	UserId          string
	DisplayName     string
	AvatarId      string
}

func getGitHubUser(accessToken string) *IntegrationUser {
	user := github.GetUser(accessToken)
	if user == nil {
		return nil
	} else {
		return &IntegrationUser{
			IntegrationName: "github",
			UserId:          string(user["id"].(json.Number)),
		}
	}
}

func getDiscordUser(accessToken string) *IntegrationUser {
	user := discord.GetUser(accessToken)
	avatar, ok := user["avatar"].(string)
    if !ok {
        avatar = "" // Or set a default avatar URL if preferred
    }
	if user == nil {
		logger.Error("User not found for token: %s", accessToken)
		return nil
	} else {
		return &IntegrationUser{
			IntegrationName: "discord",
			UserId:          user["id"].(string),
			DisplayName:     user["global_name"].(string),
			AvatarId:        avatar,
		}
	}
}

func GetUser(integrationName string, accessToken string) *IntegrationUser {
	switch strings.ToLower(integrationName) {
	case "github":
		return getGitHubUser(accessToken)
	case "discord":
		return getDiscordUser(accessToken)
	default:
		return nil
	}

}
