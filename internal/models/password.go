package password

// import "time"

// Password represents a password entry in the password manager.
type Password struct {
	ID       int    `json:"id"`
	Site     string `json:"site"`
	Username string `json:"username"`
	Pass     string `json:"pass"`
	// Note      string
	// CreatedAt time.Time
	// UpdatedAt time.Time
}
