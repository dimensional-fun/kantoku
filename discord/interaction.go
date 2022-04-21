package discord

type Interaction struct {
	ID            string          `json:"id"`
	ApplicationID string          `json:"application_id"`
	Type          int             `json:"type"`
	Data          InteractionData `json:"data,omitempty"`
	GuildID       string          `json:"guild_id,omitempty"`
	ChannelID     string          `json:"channel_id,omitempty"`
	Member        Member          `json:"member,omitempty"`
	User          *User           `json:"user,omitempty"`
	Token         string          `json:"token"`
	Version       int             `json:"version"`
	Message       *Message        `json:"message,omitempty"`
	Locale        string          `json:"locale,omitempty"`
	GuildLocale   string          `json:"guild_locale,omitempty"`
}

type InteractionData struct {
	ID            string                                    `json:"id,omitempty"`
	Name          string                                    `json:"name,omitempty"`
	Type          int                                       `json:"type,omitempty"`
	Resolved      *InteractionResolvedData                  `json:"resolved,omitempty"`
	Options       []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	CustomID      string                                    `json:"custom_id,omitempty"`
	ComponentType int                                       `json:"component_type,omitempty"`
	Values        []SelectOption                            `json:"values,omitempty"`
	TargetID      string                                    `json:"target_id,omitempty"`
}

type InteractionResolvedData struct {
	Users    map[string]User    `json:"users,omitempty"`
	Members  map[string]Member  `json:"members,omitempty"`
	Roles    map[string]Role    `json:"roles,omitempty"`
	Messages map[string]Message `json:"messages,omitempty"`
	Channels map[string]Channel `json:"channels,omitempty"`
}

type ApplicationCommandInteractionDataOption struct {
	Name    string                                    `json:"name"`
	Type    int                                       `json:"type"`
	Value   interface{}                               `json:"value,omitempty"`
	Options []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Focused bool                                      `json:"focused,omitempty"`
}

type SelectOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Emoji       *Emoji `json:"emoji,omitempty"`
	Default     bool   `json:"default,omitempty"`
}

type InteractionResponse struct {
	Type int `json:"type"`
}
