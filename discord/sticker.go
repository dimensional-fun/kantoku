package discord

type Sticker struct {
	Id          Snowflake `json:"id"`
	PackId      Snowflake `json:"pack_id,omitempty"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Tags        string    `json:"tags"`
	Asset       string    `json:"asset"`
	Type        int       `json:"type"`
	FormatType  int       `json:"format_type"`
	Available   bool      `json:"available,omitempty"`
	GuildId     Snowflake `json:"guild_id"`
	User        User      `json:"user,omitempty"`
	SortValue   int       `json:"sort_value,omitempty"`
}

type StickerItem struct {
	Id         Snowflake `json:"id"`
	Name       string    `json:"name"`
	FormatType int       `json:"format_type"`
}
