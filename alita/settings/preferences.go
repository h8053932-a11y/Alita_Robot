package settings

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/divkix/Alita_Robot/alita/db/lang"
	"github.com/divkix/Alita_Robot/alita/db/rules"
	"github.com/divkix/Alita_Robot/alita/utils/chattz"
	"github.com/divkix/Alita_Robot/alita/utils/digestsettings"
)

// Timezone is the group's display timezone, stored as an IANA name.
var Timezone = &Declaration{
	Key: "timezone", Type: Text, Default: "UTC", Right: Admin,
	read: func(ctx context.Context, chatID int64) any {
		name, err := chattz.NameContext(ctx, chatID)
		if err != nil {
			return "UTC"
		}
		return name
	},
	allowed: func(value any) error {
		_, _, err := chattz.Parse(value.(string))
		return err
	},
	write: func(input WriteContext, value any) error {
		_, err := chattz.Set(input.ChatID, value.(string))
		return err
	},
}

// ReadTimezone preserves repository errors while normalizing invalid stored names.
func ReadTimezoneContext(ctx context.Context, chatID int64) (string, error) {
	name, err := chattz.NameContext(ctx, chatID)
	if err != nil {
		return "", err
	}
	return Timezone.effective(name).(string), nil
}

var Language = &Declaration{
	Key: "language.code", Type: Choice, Default: "en", Right: Admin,
	Options: []string{"en", "es", "fr", "hi", "ru", "pt", "id", "tr"},
	read:    func(ctx context.Context, chatID int64) any { return lang.GetGroupLanguage(ctx, chatID) },
	write: func(input WriteContext, value any) error {
		return lang.ChangeGroupLanguage(input.ChatID, value.(string))
	},
}

var PrivateRules = &Declaration{
	Key: "rules.private_rules", Type: Switch, Default: false, Right: AdminBot,
	read: func(ctx context.Context, chatID int64) any { return rules.GetChatRulesInfo(chatID).Private },
	write: func(input WriteContext, value any) error {
		return rules.SetPrivateRules(input.ChatID, value.(bool))
	},
}
