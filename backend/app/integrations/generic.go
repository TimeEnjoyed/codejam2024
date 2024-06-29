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
	ServiceUserName string
	AvatarUrl       string
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
	if user == nil {
		logger.Error("User not found for token: %s", accessToken)
		return nil
	} else {
		var avatar = ""
		if user["avatar"] != nil {
			avatar = user["avatar"].(string)
		}
		return &IntegrationUser{
			IntegrationName: "discord",
			UserId:          user["id"].(string),
			ServiceUserName: user["global_name"].(string),
			AvatarUrl:       avatar,
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
