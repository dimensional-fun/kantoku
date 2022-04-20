package discord

type Sticker struct {
	ID          string  `json:"id"`
	PackID      string  `json:"pack_id,omitempty"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Tags        string  `json:"tags"`
	Asset       string  `json:"asset"`
	Type        int     `json:"type"`
	FormatType  int     `json:"format_type"`
	Available   bool    `json:"available,omitempty"`
	GuildID     string  `json:"guild_id"`
	User        User    `json:"user,omitempty"`
	SortValue   int     `json:"sort_value,omitempty"`
}

type StickerItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FormatType int    `json:"format_type"`
}
