package discord

type Message struct {
	ID               string              `json:"id"`
	ChannelID        string              `json:"channel_id"`
	GuildID          string              `json:"guild_id,omitempty"`
	Author           *User               `json:"author"`
	Member           *Member             `json:"member,omitempty"`
	Content          string              `json:"content,omitempty"`
	Timestamp        string              `json:"timestamp"`
	EditedTimestamp  *string             `json:"edited_timestamp"`
	TTS              bool                `json:"tts"`
	MentionEveryone  bool                `json:"mention_everyone"`
	Mentions         []User              `json:"mentions"`
	MentionRoles     []Role              `json:"mention_roles"`
	MentionChannels  []ChannelMention    `json:"mention_channels,omitempty"`
	Attachments      []MessageAttachment `json:"attachments"`
	Embeds           []MessageEmbed      `json:"embeds"`
	Reactions        []MessageReaction   `json:"reactions,omitempty"`
	Nonce            string              `json:"nonce,omitempty"`
	Pinned           bool                `json:"pinned"`
	WebhookID        string              `json:"webhook_id,omitempty"`
	Type             int                 `json:"type"`
	Activity         *MessageActivity    `json:"activity,omitempty"`
	Application      *MessageApplication `json:"application,omitempty"`
	ApplicationID    string              `json:"application_id,omitempty"`
	MessageReference *MessageReference   `json:"message_reference,omitempty"`
	Flags            int                 `json:"flags,omitempty"`
	ReferenceMessage *Message            `json:"reference_message,omitempty"`
	Interaction      *MessageInteraction `json:"interaction,omitempty"`
	Thread           *Channel            `json:"thread,omitempty"`
	Components       []MessageComponent  `json:"components,omitempty"`
	StickerItems     []StickerItem       `json:"sticker_items,omitempty"`
	Stickers         []Sticker           `json:"stickers,omitempty"`
}

type ChannelMention struct {
	ID      string `json:"id"`
	GuildID string `json:"guild_id"`
	Type    int    `json:"type"`
	Name    string `json:"name"`
}

type MessageInteraction struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
	Name string `json:"name"`
	User User   `json:"user"`
}

type MessageActivity struct {
	Type    int    `json:"type"`
	PartyID string `json:"party_id,omitempty"`
}

type MessageComponent struct {
	Type        int                `json:"type"`
	CustomID    string             `json:"custom_id,omitempty"`
	Disabled    bool               `json:"disabled,omitempty"`
	Style       int                `json:"style,omitempty"`
	Label       string             `json:"label,omitempty"`
	Emoji       *Emoji             `json:"emoji,omitempty"`
	Url         string             `json:"url,omitempty"`
	Options     []SelectOption     `json:"options,omitempty"`
	Placeholder string             `json:"placeholder,omitempty"`
	MinValues   *float64           `json:"min_values,omitempty"`
	MaxValues   *float64           `json:"max_values,omitempty"`
	Components  []MessageComponent `json:"components,omitempty"`
}

type MessageApplication struct {
	ID          string `json:"id"`
	CoverImage  string `json:"cover_image"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
}

type MessageReference struct {
	MessageID       string `json:"message_id,omitempty"`
	ChannelID       string `json:"channel_id,omitempty"`
	GuildID         string `json:"guild_id,omitempty"`
	FailIfNotExists bool   `json:"fail_if_not_exists,omitempty"`
}

type MessageReaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}

type MessageAttachment struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	Description string `json:"description,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Size        int    `json:"size"`
	Url         string `json:"url"`
	ProxyURL    string `json:"proxy_url"`
	Height      int    `json:"height,omitempty"`
	Width       int    `json:"width,omitempty"`
	Ephemeral   bool   `json:"ephemeral,omitempty"`
}

type MessageEmbed struct {
	URL         string                 `json:"url,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *MessageEmbedVideo     `json:"video,omitempty"`
	Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []MessageEmbedField    `json:"fields,omitempty"`
}

type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type MessageEmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type MessageEmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type MessageEmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type MessageEmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

type MessageEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type MessageEmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}
