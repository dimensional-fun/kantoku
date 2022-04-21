package discord

type Channel struct {
	ID                         string                `json:"id"`
	Type                       int                   `json:"type"`
	GuildID                    string                `json:"guild_id,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                       string                `json:"name"`
	Topic                      *string               `json:"topic,omitempty"`
	Nsfw                       bool                  `json:"nsfw,omitempty"`
	LastMessageID              *string               `json:"last_message_id,omitempty"`
	Bitrate                    int                   `json:"bitrate,omitempty"`
	UserLimit                  int                   `json:"user_limit,omitempty"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user"`
	Recipients                 []User                `json:"recipients,omitempty"`
	String                     *string               `json:"string,omitempty"`
	OwnerID                    string                `json:"owner_id,omitempty"`
	ApplicationID              string                `json:"application_id,omitempty"`
	ParentID                   *string               `json:"parent_id,omitempty"`
	LastPinTimestamp           *string               `json:"last_pin_timestamp,omitempty"`
	RtcRegion                  *string               `json:"rtc_region,omitempty"`
	VideoQualityMode           int                   `json:"video_quality_mode,omitempty"`
	MessageCount               int                   `json:"message_count,omitempty"`
	MemberCount                int                   `json:"member_count,omitempty"`
	ThreadMetadata             ThreadChannelMetadata `json:"thread_metadata,omitempty"`
	Member                     ThreadMember          `json:"member,omitempty"`
	DefaultAutoArchiveDuration int                   `json:"default_auto_archive_duration,omitempty"`
	Permissions                string                `json:"permissions,omitempty"`
}

type ThreadMember struct {
	ID            string `json:"id,omitempty"`
	UserID        string `json:"user_id,omitempty"`
	JoinTimestamp string `json:"join_timestamp"`
	Flags         int    `json:"flags"`
}

type ThreadChannelMetadata struct {
	Archived            bool   `json:"archived"`
	AutoArchiveDuration int    `json:"auto_archive_duration"`
	ArchiveTimestamp    string `json:"archive_timestamp"`
	Locked              bool   `json:"locked"`
	Invitable           bool   `json:"invitable,omitempty"`
}

type PermissionOverwrite struct {
	ID    string `json:"id"`
	Type  int    `json:"type"`
	Allow string `json:"allow"`
	Deny  string `json:"deny"`
}
