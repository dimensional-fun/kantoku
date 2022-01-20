package discord

type Role struct {
	Id           Snowflake `json:"id"`
	Name         string    `json:"name"`
	Color        int       `json:"color"`
	Hoist        bool      `json:"hoist"`
	Icon         *string   `json:"icon,omitempty"`
	UnicodeEmoji *string   `json:"unicode_emoji,omitempty"`
	Position     int       `json:"position"`
	Permissions  string    `json:"permissions"`
	Managed      bool      `json:"managed"`
	Mentionable  bool      `json:"mentionable"`
	Tags         []RoleTag `json:"tags,omitempty"`
}

type RoleTag struct {
	BotId             Snowflake `json:"bot_id,omitempty"`
	IntegrationId     Snowflake `json:"integration_id,omitempty"`
	PremiumSubscriber *bool     `json:"premium_subscriber,omitempty"`
}
