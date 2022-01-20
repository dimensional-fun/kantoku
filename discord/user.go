package discord

type UserFlags int

type User struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	Avatar        string    `json:"avatar"`
	Bot           bool      `json:"bot,omitempty"`
	System        bool      `json:"system,omitempty"`
	MFAEnabled    bool      `json:"mfa_enabled,omitempty"`
	Banner        *string   `json:"banner,omitempty"`
	AccentColor   *int      `json:"accent_color,omitempty"`
	Locale        string    `json:"locale,omitempty"`
	Verified      bool      `json:"verified,omitempty"`
	Email         *string   `json:"email,omitempty"`
	Flags         int       `json:"flags,omitempty"`
	PremiumType   int       `json:"premium_type,omitempty"`
	PublicFlags   UserFlags `json:"public_flags,omitempty"`
	Member        *Member   `json:"member,omitempty"`
}
